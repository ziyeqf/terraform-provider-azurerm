---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_scheduled_query_rules_alert_v2"
description: |-
  Manages an AlertingAction Scheduled Query Rules Version 2 resource within Azure Monitor
---

# azurerm_monitor_scheduled_query_rules_alert_v2

Manages an AlertingAction Scheduled Query Rules Version 2 resource within Azure Monitor

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_application_insights" "example" {
  name                = "example-ai"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_monitor_action_group" "example" {
  name                = "example-mag"
  resource_group_name = azurerm_resource_group.example.name
  short_name          = "test mag"
}

resource "azurerm_monitor_scheduled_query_rules_alert_v2" "example" {
  name                = "example-msqrv2"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  evaluation_frequency = "PT10M"
  window_duration      = "PT10M"
  scopes               = [azurerm_application_insights.example.id]
  severity             = 4
  criteria {
    query                   = <<-QUERY
      requests
        | summarize CountByCountry=count() by client_CountryOrRegion
      QUERY
    time_aggregation_method = "Maximum"
    threshold               = 17.5
    operator                = "LessThan"

    resource_id_column    = "client_CountryOrRegion"
    metric_measure_column = "CountByCountry"
    dimension {
      name     = "client_CountryOrRegion"
      operator = "Exclude"
      values   = ["123"]
    }
    failing_periods {
      minimum_failing_periods_to_trigger_alert = 1
      number_of_evaluation_periods             = 1
    }
  }

  auto_mitigation_enabled          = true
  workspace_alerts_storage_enabled = false
  description                      = "example sqr"
  display_name                     = "example-sqr"
  enabled                          = true
  query_time_range_override        = "PT1H"
  skip_query_validation            = true
  action {
    action_groups = [azurerm_monitor_action_group.example.id]
    custom_properties = {
      key  = "value"
      key2 = "value2"
    }
  }

  tags = {
    key  = "value"
    key2 = "value2"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Monitor Scheduled Query Rule. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Monitor Scheduled Query Rule should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the Monitor Scheduled Query Rule should exist. Changing this forces a new resource to be created.

* `criteria` - (Required) A `criteria` block as defined below.

* `evaluation_frequency` - (Required) How often the scheduled query rule is evaluated, represented in ISO 8601 duration format. 

* `scopes` - (Required) Specifies the list of resource ids that this scheduled query rule is scoped to. Changing this forces a new resource to be created.

* `severity` - (Required) Severity of the alert. Should be an integer between 0 and 4. Value of 0 is severest. 

* `window_duration` - (Required) Specifies the period of time in ISO 8601 duration format on which the Scheduled Query Rule will be executed (bin size).

* `action` - (Optional) An `action` block as defined below.

* `auto_mitigation_enabled` - (Optional) Specifies the flag that indicates whether the alert should be automatically resolved or not. Value should be `true` or `false`. The default is `false`.

* `workspace_alerts_storage_enabled` - (Optional) Specifies the flag which indicates whether this scheduled query rule check if storage is configured. Value should be `true` or `false`. The default is `false`. 

* `description` - (Optional) Specifies the description of the scheduled query rule.

* `display_name` - (Optional) Specifies the display name of the alert rule.

* `enabled` - (Optional) Specifies the flag which indicates whether this scheduled query rule is enabled. Value should be `true` or `false`. The default is `true`.

* `mute_actions_after_alert_duration` - (Optional) Mute actions for the chosen period of time in ISO 8601 duration format after the alert is fired. 

-> **NOTE** `auto_mitigation_enabled` and `mute_actions_after_alert_duration` are mutually exclusive and cannot both be set.

* `query_time_range_override` - (Optional) If specified then overrides the query time range, default is `window_duration`*`number_of_evaluation_periods`.

* `skip_query_validation` - (Optional) Specifies the flag which indicates whether the provided query should be validated or not. The default is false.

* `tags` - (Optional) A mapping of tags which should be assigned to the Monitor Scheduled Query Rule.

* `target_resource_types` - (Optional) List of resource type of the target resource(s) on which the alert is created/updated. For example if the scope is a resource group and targetResourceTypes is `Microsoft.Compute/virtualMachines`, then a different alert will be fired for each virtual machine in the resource group which meet the alert criteria. 

---

An `action` block supports the following:

* `action_groups` - (Optional) List of Action Group resource ids to invoke when the alert fires.

* `custom_properties` - (Optional) Specifies the properties of an alert payload.

---

A `criteria` block supports the following:

* `operator` - (Required) Specifies the criteria operator. Possible values are `Equals`, `GreaterThan`, `GreaterThanOrEqual`, `LessThan`,and `LessThanOrEqual`. 

* `query` - (Required) The query to run on logs. The results returned by this query are used to populate the alert. 

* `threshold` - (Required) Specifies the criteria threshold value that activates the alert.

* `time_aggregation_method` - (Required) The type of aggregation to apply to the data points in aggregation granularity. Possible values are `Average`, `Count`, `Maximum`, `Minimum`,and `Total`.

* `dimension` - (Optional) A `dimension` block as defined below.

* `failing_periods` - (Optional) A `failing_periods` block as defined below.

* `metric_measure_column` - (Optional) Specifies the column containing the metric measure number.

* `resource_id_column` - (Optional) Specifies the column containing the resource id. The content of the column must be an uri formatted as resource id. 

---

A `dimension` block supports the following:

* `name` - (Required) Name of the dimension.

* `operator` - (Required) Operator for dimension values. Possible values are `Exclude`,and `Include`.

* `values` - (Required) List of dimension values. Use a wildcard `*` to collect all.

---

A `failing_periods` block supports the following:

* `minimum_failing_periods_to_trigger_alert` - (Required) Specifies the number of violations to trigger an alert. Should be smaller or equal to `number_of_evaluation_periods`. Possible value is integer between 1 and 6.

* `number_of_evaluation_periods` - (Required) Specifies the number of aggregated look-back points. The look-back time window is calculated based on the aggregation granularity `window_duration` and the selected number of aggregated points. Possible value is integer between 1 and 6. 

## Attributes Reference

The following Attributes are exported:

* `id` - The id of the Monitor Scheduled Query Rule.

* `created_with_api_version` - The api-version used when creating this alert rule.

* `is_a_legacy_log_analytics_rule` - True if this alert rule is a legacy Log Analytic Rule.

* `is_workspace_alerts_storage_configured` - The flag indicates whether this Scheduled Query Rule has been configured to be stored in the customer's storage.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Monitor Scheduled Query Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Monitor Scheduled Query Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Monitor Scheduled Query Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Monitor Scheduled Query Rule.

## Import

Monitor Scheduled Query Rule Alert can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_scheduled_query_rules_alert_v2.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Insights/scheduledQueryRules/rule1
```
