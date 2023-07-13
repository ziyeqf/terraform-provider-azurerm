---
subcategory: "Service Networking"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_adp"
description: |-
  Manages a application load balancer.
---

# azurerm_adp

Manages a application load balancer.

## Example Usage

```hcl
resource "azurerm_adp" "example" {
  name = "example"
  resource_group_name = "example"
  location = "West Europe"
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the application load balancer should exist. Changing this forces a new application load balancer to be created.

* `name` - (Required) The name which should be used for this application load balancer. Changing this forces a new application load balancer to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the application load balancer should exist. Changing this forces a new application load balancer to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the application load balancer.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the application load balancer.

* `configuration_endpoint` - A `configuration_endpoint` block as defined below.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the application load balancer.
* `read` - (Defaults to 5 minutes) Used when retrieving the application load balancer.
* `update` - (Defaults to 30 minutes) Used when updating the application load balancer.
* `delete` - (Defaults to 30 minutes) Used when deleting the application load balancer.

## Import

application load balancers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_adp.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/alb1
```