---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_namespace"
description: |-
  Manages a ServiceBus Namespace.
---

# azurerm_servicebus_namespace

Manages a ServiceBus Namespace.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "terraform-servicebus"
  location = "West Europe"
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "tfex-servicebus-namespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"

  tags = {
    source = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the ServiceBus Namespace resource . Changing this forces a
    new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the namespace.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Required) Defines which tier to use. Options are `Basic`, `Standard` or `Premium`. Please note that setting this field to `Premium` will force the creation of a new resource.

* `identity` - (Optional) An `identity` block as defined below.

* `capacity` - (Optional) Specifies the capacity. When `sku` is `Premium`, capacity can be `1`, `2`, `4`, `8` or `16`. When `sku` is `Basic` or `Standard`, capacity can be `0` only.

* `customer_managed_key` - (Optional) An `customer_managed_key` block as defined below.

* `local_auth_enabled` - (Optional) Whether or not SAS authentication is enabled for the Service Bus namespace. Defaults to `true`.

* `public_network_access_enabled` - (Optional) Is public network access enabled for the Service Bus Namespace? Defaults to `true`.

* `minimum_tls_version` - (Optional) The minimum supported TLS version for this Service Bus Namespace. Valid values are: `1.0`, `1.1` and `1.2`. The current default minimum TLS version is `1.2`.

* `zone_redundant` - (Optional) Whether or not this resource is zone redundant. `sku` needs to be `Premium`. Defaults to `false`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this ServiceBus Namespace. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this ServiceBus namespace.

~> **NOTE:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.


---

-> **Note:** Once customer-managed key encryption has been enabled, it cannot be disabled.

A `customer_managed_key` block supports the following:


* `key_vault_key_id` - (Required) The ID of the Key Vault Key which should be used to Encrypt the data in this ServiceBus Namespace.

* `identity_id` - (Required) The ID of the User Assigned Identity that has access to the key.

* `infrastructure_encryption_enabled` - (Optional) Used to specify whether enable Infrastructure Encryption (Double Encryption).

## Attributes Reference

The following attributes are exported:

* `id` - The ServiceBus Namespace ID.

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this ServiceBus Namespace.

---

A `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this ServiceBus Namespace.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this ServiceBus Namespace.


The following attributes are exported only if there is an authorization rule named
`RootManageSharedAccessKey` which is created automatically by Azure.

* `default_primary_connection_string` - The primary connection string for the authorization
    rule `RootManageSharedAccessKey`.

* `default_secondary_connection_string` - The secondary connection string for the
    authorization rule `RootManageSharedAccessKey`.

* `default_primary_key` - The primary access key for the authorization rule `RootManageSharedAccessKey`.

* `default_secondary_key` - The secondary access key for the authorization rule `RootManageSharedAccessKey`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the ServiceBus Namespace.
* `update` - (Defaults to 30 minutes) Used when updating the ServiceBus Namespace.
* `read` - (Defaults to 5 minutes) Used when retrieving the ServiceBus Namespace.
* `delete` - (Defaults to 30 minutes) Used when deleting the ServiceBus Namespace.

## Import

Service Bus Namespace can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_servicebus_namespace.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ServiceBus/namespaces/sbns1
```
