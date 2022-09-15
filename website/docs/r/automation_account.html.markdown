---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_account"
description: |-
  Manages a Automation Account.
---

# azurerm_automation_account

Manages a Automation Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_automation_account" "example" {
  name                = "example-account"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Basic"

  tags = {
    environment = "development"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Automation Account. Changing this forces a new resource to be created.

* `public_network_access_enabled` - (Optional) Whether public network access is allowed for the container registry. Defaults to `true`.

* `resource_group_name` - (Required) The name of the resource group in which the Automation Account is created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku_name` - (Required) The SKU of the account - only `Basic` is supported at this time.

* `local_authentication_enabled` - (Optional) Whether requests using non-AAD authentication are blocked.

---

* `identity` - (Optional) An `identity` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `encryption` - (Optional) An `encryption` block as defined below.

---

An `identity` block supports the following:

* `type` - (Required) The type of identity used for this Automation Account. Possible values are `SystemAssigned`, `UserAssigned` and `SystemAssigned, UserAssigned`.

* `identity_ids` - (Optional) The ID of the User Assigned Identity which should be assigned to this Automation Account.

-> **Note:** `identity_ids` is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

--

An `encryption` block supports the following:

* `user_assigned_identity_id` - (Optional) The User Assigned Managed Identity ID to be used for accessing the Customer Managed Key for encryption.

* `key_source` - (Optional) The source of the encryption key. Possible values are `Microsoft.Keyvault` and `Microsoft.Storage`.

* `key_vault_key_id` - (Required) The ID of the Key Vault Key which should be used to Encrypt the data in this Automation Account.

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Automation Account.

* `identity` - An `identity` block as defined below.

* `dsc_server_endpoint` - The DSC Server Endpoint associated with this Automation Account.

* `dsc_primary_access_key` - The Primary Access Key for the DSC Endpoint associated with this Automation Account.

* `dsc_secondary_access_key` - The Secondary Access Key for the DSC Endpoint associated with this Automation Account.

* `hybrid_service_url` - The URL of automation hybrid service which is used for hybrid worker on-boarding With this Automation Account.
---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation Account.
* `update` - (Defaults to 30 minutes) Used when updating the Automation Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the Automation Account.

## Import

Automation Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_account.account1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1
```
