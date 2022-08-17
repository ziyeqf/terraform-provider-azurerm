package sentinel

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-07-01-preview/alertruletemplates"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceSentinelAlertRuleTemplate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceSentinelAlertRuleTemplateRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"name", "display_name"},
			},

			"display_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"name", "display_name"},
			},

			"log_analytics_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: alertruletemplates.ValidateWorkspaceID,
			},

			"scheduled_template": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"description": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"tactics": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"severity": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"query": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"query_frequency": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"query_period": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"trigger_operator": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"trigger_threshold": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"security_incident_template": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"description": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"product_filter": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSentinelAlertRuleTemplateRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRuleTemplatesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	displayName := d.Get("display_name").(string)
	workspaceID, err := alertruletemplates.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}

	// Either "name" or "display_name" must have been specified, constrained by the pluginsdk.
	var resp alertruletemplates.AlertRuleTemplate
	var nameToLog string
	if name != "" {
		nameToLog = name
		resp, err = getAlertRuleTemplateByName(ctx, client, workspaceID, name)
	} else {
		nameToLog = displayName
		resp, err = getAlertRuleTemplateByDisplayName(ctx, client, workspaceID, displayName)
	}
	if err != nil {
		return fmt.Errorf("retrieving Sentinel Alert Rule Template %q (Workspace %q / Resource Group %q): %+v", nameToLog, workspaceID.WorkspaceName, workspaceID.ResourceGroupName, err)
	}

	switch template := resp.(type) {
	case alertruletemplates.MLBehaviorAnalyticsAlertRuleTemplate:
		err = setForMLBehaviorAnalyticsAlertRuleTemplate(d, &template)
	case alertruletemplates.FusionAlertRuleTemplate:
		err = setForFusionAlertRuleTemplate(d, &template)
	case alertruletemplates.MicrosoftSecurityIncidentCreationAlertRuleTemplate:
		err = setForMsSecurityIncidentAlertRuleTemplate(d, &template)
	case alertruletemplates.ScheduledAlertRuleTemplate:
		err = setForScheduledAlertRuleTemplate(d, &template)
	default:
		return fmt.Errorf("unknown template type of Sentinel Alert Rule Template %q (Workspace %q / Resource Group %q) ID", nameToLog, workspaceID.WorkspaceName, workspaceID.ResourceGroupName)
	}

	if err != nil {
		return fmt.Errorf("setting ResourceData for Sentinel Alert Rule Template %q (Workspace %q / Resource Group %q) ID", nameToLog, workspaceID.WorkspaceName, workspaceID.ResourceGroupName)
	}

	return nil
}

func getAlertRuleTemplateByName(ctx context.Context, client *alertruletemplates.AlertRuleTemplatesClient, workspaceID *alertruletemplates.WorkspaceId, name string) (res alertruletemplates.AlertRuleTemplate, err error) {
	id := alertruletemplates.NewAlertRuleTemplateID(workspaceID.SubscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)
	template, err := client.AlertRuleTemplatesGet(ctx, id)
	if err != nil {
		return nil, err
	}

	return template.Model, nil
}

func getAlertRuleTemplateByDisplayName(ctx context.Context, client *alertruletemplates.AlertRuleTemplatesClient, workspaceID *alertruletemplates.WorkspaceId, name string) (res alertruletemplates.AlertRuleTemplate, err error) {
	templates, err := client.AlertRuleTemplatesListComplete(ctx, *workspaceID)
	if err != nil {
		return nil, err
	}
	var results []alertruletemplates.AlertRuleTemplate

	for _, template := range templates.Items {
		switch template := template.(type) {
		case alertruletemplates.FusionAlertRuleTemplate:
			if template.Properties != nil && template.Properties.DisplayName != nil && *template.Properties.DisplayName == name {
				results = append(results, template)
			}
		case alertruletemplates.MLBehaviorAnalyticsAlertRuleTemplate:
			if template.Name != nil && *template.Name == name {
				results = append(results, template)
			}
		case alertruletemplates.MicrosoftSecurityIncidentCreationAlertRuleTemplate:
			if template.Properties != nil && template.Properties.DisplayName != nil && *template.Properties.DisplayName == name {
				results = append(results, template)
			}
		case alertruletemplates.ScheduledAlertRuleTemplate:
			if template.Properties != nil && template.Properties.DisplayName != nil && *template.Properties.DisplayName == name {
				results = append(results, template)
			}
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no Alert Rule Template found with display name: %s", name)
	}
	if len(results) > 1 {
		return nil, fmt.Errorf("more than one Alert Rule Template found with display name: %s", name)
	}
	return results[0], nil
}

func setForScheduledAlertRuleTemplate(d *pluginsdk.ResourceData, template *alertruletemplates.ScheduledAlertRuleTemplate) error {
	if template.Id == nil || *template.Id == "" {
		return errors.New("empty or nil ID")
	}
	id, err := alertruletemplates.ParseAlertRuleTemplateIDInsensitively(*template.Id)
	if err != nil {
		return err
	}
	d.SetId(id.ID())
	d.Set("name", template.Name)
	if prop := template.Properties; template.Properties != nil {
		d.Set("display_name", prop.DisplayName)
	}

	return d.Set("scheduled_template", flattenScheduledAlertRuleTemplate(template.Properties))
}

func setForMsSecurityIncidentAlertRuleTemplate(d *pluginsdk.ResourceData, template *alertruletemplates.MicrosoftSecurityIncidentCreationAlertRuleTemplate) error {
	if template.Id == nil || *template.Id == "" {
		return errors.New("empty or nil ID")
	}
	id, err := alertruletemplates.ParseAlertRuleTemplateIDInsensitively(*template.Id)
	if err != nil {
		return err
	}
	d.SetId(id.ID())
	d.Set("name", template.Name)

	if prop := template.Properties; template.Properties != nil {
		d.Set("display_name", prop.DisplayName)
	}

	return d.Set("security_incident_template", flattenMsSecurityIncidentAlertRuleTemplate(template.Properties))
}

func setForFusionAlertRuleTemplate(d *pluginsdk.ResourceData, template *alertruletemplates.FusionAlertRuleTemplate) error {
	if template.Id == nil || *template.Id == "" {
		return errors.New("empty or nil ID")
	}
	id, err := alertruletemplates.ParseAlertRuleTemplateIDInsensitively(*template.Id)
	if err != nil {
		return err
	}
	d.SetId(id.ID())
	d.Set("name", template.Name)

	if prop := template.Properties; template.Properties != nil {
		d.Set("display_name", prop.DisplayName)
	}

	return nil
}

func setForMLBehaviorAnalyticsAlertRuleTemplate(d *pluginsdk.ResourceData, template *alertruletemplates.MLBehaviorAnalyticsAlertRuleTemplate) error {
	if template.Id == nil || *template.Id == "" {
		return errors.New("empty or nil ID")
	}
	id, err := alertruletemplates.ParseAlertRuleTemplateIDInsensitively(*template.Id)
	if err != nil {
		return err
	}
	d.SetId(id.ID())
	d.Set("name", template.Name)

	if prop := template.Properties; template.Properties != nil {
		d.Set("display_name", prop)
	}

	return nil
}

func flattenScheduledAlertRuleTemplate(input *alertruletemplates.ScheduledAlertRuleTemplateProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	description := ""
	if input.Description != nil {
		description = *input.Description
	}

	tactics := []interface{}{}
	if input.Tactics != nil {
		tactics = myFlattenAlertRuleScheduledTactics(input.Tactics)
	}

	query := ""
	if input.Query != nil {
		query = *input.Query
	}

	queryFrequency := ""
	if input.QueryFrequency != nil {
		queryFrequency = *input.QueryFrequency
	}

	queryPeriod := ""
	if input.QueryPeriod != nil {
		queryPeriod = *input.QueryPeriod
	}

	triggerThreshold := 0
	if input.TriggerThreshold != nil {
		triggerThreshold = int(*input.TriggerThreshold)
	}

	return []interface{}{
		map[string]interface{}{
			"description":       description,
			"tactics":           tactics,
			"severity":          input.Severity,
			"query":             query,
			"query_frequency":   queryFrequency,
			"query_period":      queryPeriod,
			"trigger_operator":  input.TriggerOperator,
			"trigger_threshold": triggerThreshold,
		},
	}
}

func flattenMsSecurityIncidentAlertRuleTemplate(input *alertruletemplates.MicrosoftSecurityIncidentCreationAlertRuleTemplateProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	description := ""
	if input.Description != nil {
		description = *input.Description
	}

	return []interface{}{
		map[string]interface{}{
			"description":    description,
			"product_filter": input.ProductFilter,
		},
	}
}

func myFlattenAlertRuleScheduledTactics(input *[]alertruletemplates.AttackTactic) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		output = append(output, string(e))
	}

	return output
}
