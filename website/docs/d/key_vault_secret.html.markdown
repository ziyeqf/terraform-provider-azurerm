---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_secret"
description: |-
  Gets information about an existing Key Vault Secret.
---

# Data Source: azurerm_key_vault_secret

Use this data source to access information about an existing Key Vault Secret.

~> **Note:** All arguments including the secret value will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
data "azurerm_key_vault_secret" "example" {
  name         = "secret-sauce"
  key_vault_id = data.azurerm_key_vault.existing.id
}

output "secret_value" {
  value     = data.azurerm_key_vault_secret.example.value
  sensitive = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - Specifies the name of the Key Vault Secret.

* `key_vault_id` - Specifies the ID of the Key Vault instance where the Secret resides, available on the `azurerm_key_vault` Data Source / Resource. 

**NOTE:** The vault must be in the same subscription as the provider. If the vault is in another subscription, you must create an aliased provider for that subscription.

## Attributes Reference

The following attributes are exported:

* `id` - The Key Vault Secret ID.
* `resource_id` - The (Versioned) ID for this Key Vault Secret. This property points to a specific version of a Key Vault Secret, as such using this won't auto-rotate values if used in other Azure Services.
* `resource_versionless_id` - The Versionless ID of the Key Vault Secret. This property allows other Azure Services (that support it) to auto-rotate their value when the Key Vault Secret is updated.
* `value` - The value of the Key Vault Secret.
* `version` - The current version of the Key Vault Secret.
* `content_type` - The content type for the Key Vault Secret.
* `tags` - Any tags assigned to this resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 30 minutes) Used when retrieving the Key Vault Secret.
