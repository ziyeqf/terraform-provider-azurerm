package mobilenetwork_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/simpolicy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MobileNetworkSimPolicyResource struct{}

func TestAccMobileNetworkSimPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_sim_policy", "test")
	r := MobileNetworkSimPolicyResource{}
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

func TestAccMobileNetworkSimPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_sim_policy", "test")
	r := MobileNetworkSimPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMobileNetworkSimPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_sim_policy", "test")
	r := MobileNetworkSimPolicyResource{}
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

func TestAccMobileNetworkSimPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_sim_policy", "test")
	r := MobileNetworkSimPolicyResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r MobileNetworkSimPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := simpolicy.ParseSimPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.MobileNetwork.SIMPolicyClient
	resp, err := client.SimPoliciesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r MobileNetworkSimPolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-mn-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mobile_network" "test" {
  name                = "acctest-mn-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[2]s"
  mobile_country_code = "001"
  mobile_network_code = "01"
}

resource "azurerm_mobile_network_slice" "test" {
  name              = "acctest-mns-%[1]d"
  mobile_network_id = azurerm_mobile_network.test.id
  location          = "%[2]s"
  single_network_slice_selection_assistance_information {
    slice_service_type = 1
  }
}


resource "azurerm_mobile_network_service" "test" {
  name               = "acctest-mns-%[1]d"
  mobile_network_id  = azurerm_mobile_network.test.id
  location           = "%[2]s"
  service_precedence = 0

  pcc_rules {
    rule_name               = "default-rule"
    rule_precedence         = 1
    traffic_control_enabled = true

    service_data_flow_templates {
      direction      = "Uplink"
      template_name  = "IP-to-server"
      ports          = []
      protocol       = ["ip"]
      remote_ip_list = ["10.3.4.0/24"]
    }

  }
}

resource "azurerm_mobile_network_data_network" "test" {
  name              = "acctest-mndn-%[1]d"
  mobile_network_id = azurerm_mobile_network.test.id
  location          = "%[2]s"
}

`, data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkSimPolicyResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_mobile_network_sim_policy" "test" {
  name                          = "acctest-mnsp-%d"
  mobile_network_id             = azurerm_mobile_network.test.id
  location                      = "%s"
  default_slice_id              = azurerm_mobile_network_slice.test.id
  registration_timer_in_seconds = 3240

  user_equipment_aggregate_maximum_bit_rate {
    downlink = "1 Gbps"
    uplink   = "500 Mbps"
  }

  slice_configurations {
    default_data_network_id = azurerm_mobile_network_data_network.test.id
    slice_id                = azurerm_mobile_network_slice.test.id
    data_network_configurations {
      data_network_id                         = azurerm_mobile_network_data_network.test.id
      allocation_and_retention_priority_level = 9
      default_session_type                    = "IPv4"
      qos_indicator                           = 9
      preemption_capability                   = "NotPreempt"
      preemption_vulnerability                = "Preemptable"
      allowed_services_ids                    = [azurerm_mobile_network_service.test.id]
      session_aggregate_maximum_bit_rate {
        downlink = "1 Gbps"
        uplink   = "500 Mbps"
      }
    }
  }

  tags = {
    key = "value"
  }

}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkSimPolicyResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_sim_policy" "import" {
  name                          = azurerm_mobile_network_sim_policy.test.name
  mobile_network_id             = azurerm_mobile_network.test.id
  default_slice_id              = azurerm_mobile_network_slice.test.id
  location                      = "%s"
  registration_timer_in_seconds = 3240

  slice_configurations {
    default_data_network_id = azurerm_mobile_network_data_network.test.id
    slice_id                = azurerm_mobile_network_slice.test.id

    data_network_configurations {
      allocation_and_retention_priority_level = 9
      default_session_type                    = "IPv4"
      qos_indicator                           = 9
      preemption_capability                   = "NotPreempt"
      preemption_vulnerability                = "Preemptable"
      allowed_services_ids                    = [azurerm_mobile_network_service.test.id]
      data_network_id                         = azurerm_mobile_network_data_network.test.id
      session_aggregate_maximum_bit_rate {
        downlink = "1 Gbps"
        uplink   = "500 Mbps"
      }
    }
  }

  user_equipment_aggregate_maximum_bit_rate {
    downlink = "1 Gbps"
    uplink   = "500 Mbps"
  }
  tags = {
    key = "value"
  }

}
`, config, data.Locations.Primary)
}

func (r MobileNetworkSimPolicyResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_sim_policy" "test" {
  name                          = "acctest-mnsp-%d"
  mobile_network_id             = azurerm_mobile_network.test.id
  location                      = "%s"
  default_slice_id              = azurerm_mobile_network_slice.test.id
  registration_timer_in_seconds = 3240
  rfsp_index                    = 1

  slice_configurations {
    default_data_network_id = azurerm_mobile_network_data_network.test.id
    slice_id                = azurerm_mobile_network_slice.test.id
    data_network_configurations {
      allocation_and_retention_priority_level = 9
      default_session_type                    = "IPv4"
      qos_indicator                           = 9
      preemption_capability                   = "NotPreempt"
      preemption_vulnerability                = "Preemptable"
      allowed_services_ids                    = [azurerm_mobile_network_service.test.id]
      data_network_id                         = azurerm_mobile_network_data_network.test.id
      max_buffered_packets                    = 200
      session_aggregate_maximum_bit_rate {
        downlink = "1 Gbps"
        uplink   = "500 Mbps"
      }
    }
  }

  user_equipment_aggregate_maximum_bit_rate {
    downlink = "1 Gbps"
    uplink   = "500 Mbps"
  }
  tags = {
    key = "value"
  }

}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkSimPolicyResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s


resource "azurerm_mobile_network_sim_policy" "test" {
  name                          = "acctest-mnsp-%d"
  mobile_network_id             = azurerm_mobile_network.test.id
  location                      = "%s"
  default_slice_id              = azurerm_mobile_network_slice.test.id
  registration_timer_in_seconds = 3240
  rfsp_index                    = 1

  slice_configurations {
    default_data_network_id = azurerm_mobile_network_data_network.test.id
    slice_id                = azurerm_mobile_network_slice.test.id
    data_network_configurations {
      allocation_and_retention_priority_level = 9
      default_session_type                    = "IPv4"
      qos_indicator                           = 9
      preemption_capability                   = "NotPreempt"
      preemption_vulnerability                = "Preemptable"
      allowed_services_ids                    = [azurerm_mobile_network_service.test.id]
      data_network_id                         = azurerm_mobile_network_data_network.test.id
      session_aggregate_maximum_bit_rate {
        downlink = "1 Gbps"
        uplink   = "500 Mbps"
      }
    }
  }

  user_equipment_aggregate_maximum_bit_rate {
    downlink = "1 Gbps"
    uplink   = "500 Mbps"
  }
  tags = {
    key = "value2"
  }

}
`, template, data.RandomInteger, data.Locations.Primary)
}
