package recoveryservices_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationrecoveryservicesproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const HostName = "acctest-nested-server"

type HyperVHostTestResource struct{}

func (r HyperVHostTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.RecoveryServices.VaultReplicationProvider

	parsedFabricId, err := replicationrecoveryservicesproviders.ParseReplicationFabricID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.ListByReplicationFabricsComplete(ctx, *parsedFabricId)
	for _, item := range resp.Items {
		if item.Properties != nil && item.Properties.FriendlyName != nil && *item.Properties.FriendlyName == HostName {
			return utils.Bool(true), nil
		}
	}

	return utils.Bool(false), nil
}

func (HyperVHostTestResource) virtualMachineExists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) error {
	id, err := parse.VirtualMachineID(state.ID)
	if err != nil {
		return err
	}

	resp, err := client.Compute.VMClient.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s does not exist", *id)
		}

		return fmt.Errorf("Bad: Get on client: %+v", err)
	}

	return nil
}

func (HyperVHostTestResource) rebootVirtualMachine() func(context.Context, *clients.Client, *pluginsdk.InstanceState) error {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		client := clients.Compute.VMClient
		id, err := parse.VirtualMachineID(state.ID)
		if err != nil {
			return err
		}

		future, err := client.Restart(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("restart %s: %+v", id, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for restart of %s: %+v", id, err)
		}

		return nil
	}
}

func TestAccSiteRecoveryHyperTest_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_site_recovery_services_vault_hyperv_site", "hybrid")
	r := HyperVHostTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.setupTemplate(data),
			Check: acceptance.ComposeTestCheckFunc(
				data.CheckWithClientForResource(r.virtualMachineExists, "azurerm_windows_virtual_machine.host"),
				data.CheckWithClientForResource(r.rebootVirtualMachine(), "azurerm_windows_virtual_machine.host"),
			),
		},
		{
			Config: r.template(data),
		},
		{
			Config: r.basic(data),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
	})
}

func (r HyperVHostTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
`, r.template(data))
}

func (HyperVHostTestResource) keyVault(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "hybird" {
  name                = local.keyvault_name
  resource_group_name = azurerm_resource_group.hybrid.name
  location            = azurerm_resource_group.hybrid.location
  sku_name            = "standard"
  tenant_id           = data.azurerm_client_config.current.tenant_id

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Backup",
      "Create",
      "Decrypt",
      "Delete",
      "Encrypt",
      "Get",
      "Import",
      "List",
      "Purge",
      "Recover",
      "Restore",
      "Sign",
      "UnwrapKey",
      "Update",
      "Verify",
      "WrapKey",
    ]

    secret_permissions = [
      "Backup",
      "Delete",
      "Get",
      "List",
      "Purge",
      "Recover",
      "Restore",
      "Set",
    ]

    certificate_permissions = [
      "Create",
      "Delete",
      "DeleteIssuers",
      "Get",
      "GetIssuers",
      "Import",
      "List",
      "ListIssuers",
      "ManageContacts",
      "ManageIssuers",
      "Purge",
      "SetIssuers",
      "Update",
    ]
  }

  enabled_for_deployment          = true
  enabled_for_template_deployment = true
}

resource "azurerm_key_vault_certificate" "winrm" {
  name         = local.cert_name
  key_vault_id = azurerm_key_vault.hybird.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      extended_key_usage = ["1.3.6.1.5.5.7.3.1"]

      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=${local.vm_name}"
      validity_in_months = 12
    }
  }
}

`)
}

func (HyperVHostTestResource) securityGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_network_security_group" "hybrid" {
  name                = local.nsg_name
  location            = azurerm_resource_group.hybrid.location
  resource_group_name = azurerm_resource_group.hybrid.name

  security_rule {
    name                       = "allow-rdp"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "3389"
    source_address_prefix      = "167.220.255.65"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "allow-winrm"
    priority                   = 101
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "5986"
    source_address_prefix      = "167.220.255.65"
    destination_address_prefix = "*"
  }


  lifecycle {
    ignore_changes = [security_rule]
  }
}

resource "azurerm_network_interface_security_group_association" "hybrid" {
  network_interface_id      = azurerm_network_interface.host.id
  network_security_group_id = azurerm_network_security_group.hybrid.id
}
`)
}

func (HyperVHostTestResource) recovery(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_recovery_services_vault" "hybrid" {
  name                = local.recovery_vault_name
  location            = azurerm_resource_group.hybrid.location
  resource_group_name = azurerm_resource_group.hybrid.name
  sku                 = "Standard"

  soft_delete_enabled = false
}

resource "azurerm_site_recovery_services_vault_hyperv_site" "hybrid" {
  name              = local.recovery_site_name
  recovery_vault_id = azurerm_recovery_services_vault.hybrid.id
}

resource "azurerm_recovery_services_vault_hyperv_host_registration_key" "hybrid" {
  site_recovery_services_vault_hyperv_site_id = azurerm_site_recovery_services_vault_hyperv_site.hybrid.id
  validate_in_hours                           = 120
}

`)
}

func (r HyperVHostTestResource) setupTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

provider "azurerm" {
  features {}
}

provider "random" {}

locals {
  rg_name             = "acctest-nested-rg-%[1]d"
  location            = "%[2]s"
  vn_name             = "acctest-nested-vn-%[1]d"
  ip_name             = "acctest-nested-ip-%[1]d"
  vm_name             = "acctest-nested-vm-%[1]d"
  nic_name            = "acctest-nested-nic-%[1]d"
  disk_name           = "acctest-nested-disk-%[1]d"
  keyvault_name       = "acctkv%[1]d"
  nsg_name            = "acctest-nested-nsg-%[1]d"
  recovery_vault_name = "acctest-nested-recovery-vault-%[1]d"
  recovery_site_name  = "acctest-nested-recovery-site-%[1]d"
  admin_name          = "acctestadmin"
  cert_name           = "acctestcert"
}

resource "azurerm_resource_group" "hybrid" {
  name     = local.rg_name
  location = local.location
}

resource "azurerm_virtual_network" "hybrid" {
  name                = local.vn_name
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.hybrid.location
  resource_group_name = azurerm_resource_group.hybrid.name
}

resource "azurerm_subnet" "hybrid" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.hybrid.name
  virtual_network_name = azurerm_virtual_network.hybrid.name
  address_prefixes     = ["10.0.10.0/24"]
}

resource "azurerm_public_ip" "host" {
  name                = local.ip_name
  resource_group_name = azurerm_resource_group.hybrid.name
  location            = azurerm_resource_group.hybrid.location
  allocation_method   = "Static"
}

resource "azurerm_network_interface" "host" {
  name                = local.nic_name
  location            = azurerm_resource_group.hybrid.location
  resource_group_name = azurerm_resource_group.hybrid.name

  enable_ip_forwarding = true

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.hybrid.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.host.id
  }
}

resource "random_password" "host" {
  length           = 16
  special          = true
  override_special = "!#$%%*()-_=+[]{}:?"
}

resource "azurerm_windows_virtual_machine" "host" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.hybrid.name
  location            = azurerm_resource_group.hybrid.location
  size                = "Standard_D8as_v5"
  admin_username      = local.admin_name
  admin_password      = random_password.host.result
  computer_name       = "nested-Host"

  network_interface_ids = [
    azurerm_network_interface.host.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Premium_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2022-Datacenter"
    version   = "latest"
  }

  identity {
    type = "SystemAssigned"
  }

  lifecycle {
    ignore_changes = [tags, identity]
  }

  additional_unattend_content {
    setting = "AutoLogon"
    content = "<AutoLogon><Password><Value>${random_password.host.result}</Value></Password><Enabled>true</Enabled><LogonCount>1</LogonCount><Username>${local.admin_name}</Username></AutoLogon>"
  }

  winrm_listener {
    protocol        = "Https"
    certificate_url = azurerm_key_vault_certificate.winrm.secret_id
  }

  secret {
    key_vault_id = azurerm_key_vault.hybird.id

    certificate {
      store = "My"
      url   = azurerm_key_vault_certificate.winrm.secret_id
    }
  }

  connection {
    host     = self.public_ip_address
    type     = "winrm"
    user     = self.admin_username
    password = self.admin_password
    port     = 5986
    https    = true
    use_ntlm = true
    insecure = true
    timeout  = "60m"
    # script_path = "c:/windows/temp/terraform_%%RAND%%.ps1"
  }

  provisioner "remote-exec" {
    inline = [
      "powershell -command \"Set-NetConnectionProfile -InterfaceAlias Ethernet -NetworkCategory Private\"",
      "mkdir c:\\Disks",
      "mkdir C:\\Machines",
      "powershell -command \"Install-WindowsFeature -Name Hyper-V,Hyper-V-Powershell,Hyper-V-Tools -IncludeManagementTools\"",
    ]
  }
}

%[3]s

%[4]s

%[5]s
`, data.RandomInteger, data.Locations.Primary, r.recovery(data), r.keyVault(data), r.securityGroup(data))
}

func (r HyperVHostTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "null_resource" "setup_provider" {
  connection {
    host     = azurerm_windows_virtual_machine.host.public_ip_address
    type     = "winrm"
    user     = azurerm_windows_virtual_machine.host.admin_username
    password = azurerm_windows_virtual_machine.host.admin_password
    port     = 5986
    https    = true
    use_ntlm = true
    insecure = true
    # script_path = "c:/windows/temp/terraform_%%RAND%%.ps1"
    timeout = "60m"
  }

  provisioner "file" {
    content     = azurerm_recovery_services_vault_hyperv_host_registration_key.hybrid.xml_content
    destination = "c:/temp/hyperv-credential"
  }

  provisioner "file" {
    source      = "./scripts/01-provider_setup.ps1"
    destination = "c:/temp/01-provider_setup.ps1"
  }

  provisioner "remote-exec" {
    inline = [
      "curl -o C:\\Disks\\VM1.vhd \"https://software-static.download.prss.microsoft.com/pr/download/17763.737.amd64fre.rs5_release_svc_refresh.190906-2324_server_serverdatacentereval_en-us_1.vhd\" -L",
      "curl -o C:\\AzureSiteRecoveryProvider.exe \"https://aka.ms/downloaddra_eus\" -L",
      "cd c:\\temp",
      "powershell -ExecutionPolicy Bypass -File 01-provider_setup.ps1",
      "C:\\AzureSiteRecoveryProvider.exe /x:C:\\AzureSiteRecoveryProvider /q",
      "C:\\AzureSiteRecoveryProvider\\SETUPDR.EXE /i",
      "cd \"C:\\Program Files\\Microsoft Azure Site Recovery Provider\"",
      ".\\DRConfigurator.exe /r /Friendlyname \"%s\" /Credentials \"C:\\temp\\hyperv-credential\""
    ]
  }
}`, r.setupTemplate(data), HostName)
}