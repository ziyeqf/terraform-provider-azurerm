---
subcategory: "Service Networking"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_alb_association"
description: |-
  Manages a application load balancer association.
---

# azurerm_alb_association

Manages a application load balancer association.

## Example Usage

```hcl
resource "azurerm_alb_association" "example" {
  name = "example"
  location = "West Europe"
  container_application_gateway_id = "TODO"
  subnet_id = "TODO"
}
```

## Arguments Reference

The following arguments are supported:

* `container_application_gateway_id` - (Required) The ID of the TODO. Changing this forces a new application load balancer association to be created.

* `location` - (Required) The Azure Region where the application load balancer association should exist. Changing this forces a new application load balancer association to be created.

* `name` - (Required) The name which should be used for this application load balancer association. Changing this forces a new application load balancer association to be created.

* `subnet_id` - (Required) The ID of the TODO. Changing this forces a new application load balancer association to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the application load balancer association.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the application load balancer association.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the application load balancer association.
* `read` - (Defaults to 5 minutes) Used when retrieving the application load balancer association.
* `update` - (Defaults to 30 minutes) Used when updating the application load balancer association.
* `delete` - (Defaults to 30 minutes) Used when deleting the application load balancer association.

## Import

application load balancer associations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_alb_association.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/alb1/associations/association1
```