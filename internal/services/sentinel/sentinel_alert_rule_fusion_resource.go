package sentinel

import (
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-07-01-preview/alertrules"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	loganalyticsParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSentinelAlertRuleFusion() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSentinelAlertRuleFusionCreateUpdate,
		Read:   resourceSentinelAlertRuleFusionRead,
		Update: resourceSentinelAlertRuleFusionCreateUpdate,
		Delete: resourceSentinelAlertRuleFusionDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.AlertRuleID(id)
			return err
		}, importSentinelAlertRule(alertrules.AlertRuleKindFusion)),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"log_analytics_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: alertrules.ValidateWorkspaceID,
			},

			"alert_rule_template_guid": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceSentinelAlertRuleFusionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)

	workspaceID, err := alertrules.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	id := alertrules.NewAlertRuleID(workspaceID.SubscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)

	if d.IsNewResource() {
		resp, err := client.AlertRulesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existing Sentinel Alert Rule Fusion %q: %+v", id, err)
			}
		}

		model := resp.Model
		if model != nil {
			id := alertRuleID(model)
			if id != nil && *id != "" {
				return tf.ImportAsExistsError("azurerm_sentinel_alert_rule_fusion", *id)
			}
		}
	}

	params := alertrules.FusionAlertRule{
		Properties: &alertrules.FusionAlertRuleProperties{
			AlertRuleTemplateName: d.Get("alert_rule_template_guid").(string),
			Enabled:               d.Get("enabled").(bool),
		},
	}

	if !d.IsNewResource() {
		resp, err := client.AlertRulesGet(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving Sentinel Alert Rule Fusion %q: %+v", id, err)
		}

		model := *resp.Model
		if model == nil {
			return fmt.Errorf("retrieving Sentinel Alert Rule Fusion %q: model is nil", id)
		}

		if err := assertAlertRuleKind(model, alertrules.AlertRuleKindFusion); err != nil {
			return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
		}
	}

	if _, err := client.AlertRulesCreateOrUpdate(ctx, id, params); err != nil {
		return fmt.Errorf("creating Sentinel Alert Rule Fusion %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelAlertRuleFusionRead(d, meta)
}

func resourceSentinelAlertRuleFusionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := alertrules.ParseAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.AlertRulesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Sentinel Alert Rule Fusion %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Sentinel Alert Rule Fusion %q: %+v", id, err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("retrieving Sentinel Alert Rule Fusion %q: model is nil", id)
	}

	if err := assertAlertRuleKind(*model, alertrules.AlertRuleKindFusion); err != nil {
		return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
	}
	rule := (*model).(alertrules.FusionAlertRule)

	d.Set("name", rule.Name)

	workspaceId := loganalyticsParse.NewLogAnalyticsWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)
	d.Set("log_analytics_workspace_id", workspaceId.ID())

	if prop := rule.Properties; prop != nil {
		d.Set("enabled", prop.Enabled)
		d.Set("alert_rule_template_guid", prop.AlertRuleTemplateName)
	}

	return nil
}

func resourceSentinelAlertRuleFusionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := alertrules.ParseAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.AlertRulesDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting Sentinel Alert Rule Fusion %q: %+v", id, err)
	}

	return nil
}
