---
subcategory: "Service Networking"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_alb_frontend"
description: |-
  Manages a application load balancer frontend.
---

# azurerm_alb_frontend

Manages a application load balancer frontend.

## Example Usage

```hcl
resource "azurerm_alb_frontend" "example" {
  name = "example"
  location = "West Europe"
  container_application_gateway_id = "TODO"
}
```

## Arguments Reference

The following arguments are supported:

* `container_application_gateway_id` - (Required) The ID of the TODO. Changing this forces a new application load balancer frontend to be created.

* `location` - (Required) The Azure Region where the application load balancer frontend should exist. Changing this forces a new application load balancer frontend to be created.

* `name` - (Required) The name which should be used for this application load balancer frontend. Changing this forces a new application load balancer frontend to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the application load balancer frontend.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the application load balancer frontend.

* `fully_qualified_domain_name` - TODO.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the application load balancer frontend.
* `read` - (Defaults to 5 minutes) Used when retrieving the application load balancer frontend.
* `update` - (Defaults to 30 minutes) Used when updating the application load balancer frontend.
* `delete` - (Defaults to 30 minutes) Used when deleting the application load balancer frontend.

## Import

application load balancer frontends can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_alb_frontend.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ServiceNetworking/trafficControllers/alb1/frontends/frontend1
```