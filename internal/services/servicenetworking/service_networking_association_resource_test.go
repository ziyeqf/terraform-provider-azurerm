package servicenetworking_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-05-01-preview/associationsinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AssociationResource struct{}

func (r AssociationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := associationsinterface.ParseAssociationID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("while parsing resource ID: %+v", err)
	}

	resp, err := clients.ServiceNetworking.ServiceNetworkingClient.AssociationsInterface.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("while checking existence for %q: %+v", id.String(), err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestAccServiceNetworkingAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_networking_association", "test")

	// for preview only, remove before merge
	data.Locations.Primary = "northeurope"
	r := AssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceNetworkingAssociation_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_networking_association", "test")

	// for preview only, remove before merge
	data.Locations.Primary = "northeurope"
	r := AssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceNetworkingAssociation_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_networking_association", "test")

	// for preview only, remove before merge
	data.Locations.Primary = "northeurope"
	r := AssociationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r AssociationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%[1]s

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%[2]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "delegation"

    service_delegation {
      name    = "Microsoft.ServiceNetworking/trafficControllers"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

`, TrafficControllerResource{}.basic(data), data.RandomInteger)
}

func (r AssociationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

resource "azurerm_service_networking_association" "test" {
  name                             = "acct-%d"
  container_application_gateway_id = azurerm_service_networking_container_application_gateway.test.id
  subnet_id                        = azurerm_subnet.test.id
  location                         = azurerm_service_networking_container_application_gateway.test.location
}
`, r.template(data), data.RandomInteger)
}

func (r AssociationResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

resource "azurerm_service_networking_association" "test" {
  name                             = "acct-%d"
  container_application_gateway_id = azurerm_service_networking_container_application_gateway.test.id
  subnet_id                        = azurerm_subnet.test.id
  location                         = azurerm_service_networking_container_application_gateway.test.location
  tags = {
    key = "value"
  }
}
`, r.template(data), data.RandomInteger)
}
