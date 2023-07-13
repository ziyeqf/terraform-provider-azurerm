package servicenetworking_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-05-01-preview/frontendsinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type FrontendResource struct{}

func (r FrontendResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := frontendsinterface.ParseFrontendID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("while parsing resource ID: %+v", err)
	}

	resp, err := clients.ServiceNetworking.ServiceNetworkingClient.FrontendsInterface.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("while checking existence for %q: %+v", id.String(), err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestAccServiceNetworkingFrontend_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_alb_frontend", "test")

	// for preview only, remove before merge
	data.Locations.Primary = "northeurope"
	r := FrontendResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fully_qualified_domain_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceNetworkingFrontend_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_alb_frontend", "test")

	// for preview only, remove before merge
	data.Locations.Primary = "northeurope"
	r := FrontendResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fully_qualified_domain_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceNetworkingFrontend_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_alb_frontend", "test")

	// for preview only, remove before merge
	data.Locations.Primary = "northeurope"
	r := FrontendResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fully_qualified_domain_name").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fully_qualified_domain_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (r FrontendResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

resource "azurerm_alb_frontend" "test" {
  name                             = "acct-frnt-%d"
  container_application_gateway_id = azurerm_alb.test.id
  location                         = azurerm_alb.test.location
}
`, TrafficControllerResource{}.basic(data), data.RandomInteger)
}

func (r FrontendResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

resource "azurerm_alb_frontend" "test" {
  name                             = "acct-frnt-%d"
  container_application_gateway_id = azurerm_alb.test.id
  location                         = azurerm_alb.test.location
  tags = {
    "tag1" = "value1"
  }
}
`, TrafficControllerResource{}.basic(data), data.RandomInteger)
}
