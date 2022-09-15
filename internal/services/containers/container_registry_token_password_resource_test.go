package containers_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ContainerRegistryTokenPasswordResource struct {
	Expiry time.Time
}

func TestAccContainerRegistryTokenPassword_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_token_password", "test")
	r := ContainerRegistryTokenPasswordResource{Expiry: time.Now().Add(time.Hour)}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password1.0.value"),
	})
}

func TestAccContainerRegistryTokenPassword_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_token_password", "test")
	r := ContainerRegistryTokenPasswordResource{Expiry: time.Now().Add(time.Hour)}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password1.0.value", "password2.0.value"),
	})
}

func TestAccContainerRegistryTokenPassword_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_token_password", "test")
	r := ContainerRegistryTokenPasswordResource{Expiry: time.Now().Add(time.Hour)}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password1.0.value"),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password1.0.value", "password2.0.value"),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password1.0.value"),
	})
}

func TestAccContainerRegistryTokenPassword_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_token_password", "test")
	r := ContainerRegistryTokenPasswordResource{Expiry: time.Now().Add(time.Hour)}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r ContainerRegistryTokenPasswordResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	client := clients.Containers.TokensClient

	id, err := parse.ContainerRegistryTokenPasswordID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.RegistryName, id.TokenName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	props := resp.TokenProperties
	if props == nil {
		return nil, fmt.Errorf("checking for presence of existing %s: unexpected nil tokenProperties", id)
	}
	cred := props.Credentials
	if cred == nil {
		return nil, fmt.Errorf("checking for presence of existing %s: unexpected nil tokenProperties.credentials", id)
	}
	pwds := cred.Passwords
	if pwds == nil {
		return nil, fmt.Errorf("checking for presence of existing %s: unexpected nil tokenProperties.credentials.passwords", id)
	}
	// ACR token with no password returns a empty array for ".password"
	return utils.Bool(len(*pwds) != 0), nil
}

func (r ContainerRegistryTokenPasswordResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_token_password" "test" {
  container_registry_token_id = azurerm_container_registry_token.test.id
  password1 {}
}
`, template)
}

func (r ContainerRegistryTokenPasswordResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_token_password" "test" {
  container_registry_token_id = azurerm_container_registry_token.test.id
  password1 {
    expiry = %q
  }
  password2 {
    expiry = %q
  }
}
`, template, r.Expiry.Format(time.RFC3339), r.Expiry.Format(time.RFC3339))
}

func (r ContainerRegistryTokenPasswordResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_token_password" "import" {
  container_registry_token_id = azurerm_container_registry_token.test.id
  password1 {}
}
`, template)
}

func (r ContainerRegistryTokenPasswordResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
}

# use system wide scope map for tests
data "azurerm_container_registry_scope_map" "pull_repos" {
  name                    = "_repositories_pull"
  container_registry_name = azurerm_container_registry.test.name
  resource_group_name     = azurerm_container_registry.test.resource_group_name
}

resource "azurerm_container_registry_token" "test" {
  name                    = "testtoken-%d"
  resource_group_name     = azurerm_resource_group.test.name
  container_registry_name = azurerm_container_registry.test.name
  scope_map_id            = data.azurerm_container_registry_scope_map.pull_repos.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
