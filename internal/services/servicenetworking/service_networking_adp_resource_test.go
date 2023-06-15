package servicenetworking_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-05-01-preview/trafficcontrollerinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type TrafficControllerResource struct{}

func (r TrafficControllerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := trafficcontrollerinterface.ParseTrafficControllerID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("while parsing resource ID: %+v", err)
	}

	resp, err := clients.ServiceNetworking.ServiceNetworkingClient.TrafficControllerInterface.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("while checking existence for %q: %+v", id.String(), err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestAccServiceNetworkingADP_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_networking_adp", "test")
	r := TrafficControllerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func (r TrafficControllerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acct-sn-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r TrafficControllerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

resource "azurerm_service_networking_adp" "test" {
  name                = "acct-%d"
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}
