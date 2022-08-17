package sentinel

import (
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-07-01-preview/automationrules"
	"log"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	loganalyticsParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/parse"
	loganalyticsValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSentinelAutomationRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSentinelAutomationRuleCreateUpdate,
		Read:   resourceSentinelAutomationRuleRead,
		Update: resourceSentinelAutomationRuleCreateUpdate,
		Delete: resourceSentinelAutomationRuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AutomationRuleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(5 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"log_analytics_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: loganalyticsValidate.LogAnalyticsWorkspaceID,
			},

			"display_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"order": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 1000),
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"expiration": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppress.RFC3339Time,
				ValidateFunc:     validation.IsRFC3339Time,
			},

			"condition": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"property": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyAccountAadTenantId),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyAccountAadUserId),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyAccountNTDomain),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyAccountName),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyAccountObjectGuid),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyAccountPUID),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyAccountSid),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyAccountUPNSuffix),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyAzureResourceResourceId),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyAzureResourceSubscriptionId),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyCloudApplicationAppId),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyCloudApplicationAppName),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyDNSDomainName),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyFileDirectory),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyFileHashValue),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyFileName),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyHostAzureID),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyHostNTDomain),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyHostName),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyHostNetBiosName),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyHostOSVersion),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyIPAddress),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyIncidentDescription),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyIncidentProviderName),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyIncidentRelatedAnalyticRuleIds),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyIncidentSeverity),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyIncidentStatus),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyIncidentTactics),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyIncidentTitle),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyIoTDeviceId),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyIoTDeviceModel),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyIoTDeviceName),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyIoTDeviceOperatingSystem),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyIoTDeviceType),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyIoTDeviceVendor),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyMailMessageDeliveryAction),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyMailMessageDeliveryLocation),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyMailMessagePOneSender),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyMailMessagePTwoSender),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyMailMessageRecipient),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyMailMessageSenderIP),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyMailMessageSubject),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyMailboxDisplayName),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyMailboxPrimaryAddress),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyMailboxUPN),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyMalwareCategory),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyMalwareName),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyProcessCommandLine),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyProcessId),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyRegistryKey),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyRegistryValueData),
								string(automationrules.AutomationRulePropertyConditionSupportedPropertyUrl),
							}, false),
						},

						"operator": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(automationrules.AutomationRulePropertyConditionSupportedOperatorContains),
								string(automationrules.AutomationRulePropertyConditionSupportedOperatorEndsWith),
								string(automationrules.AutomationRulePropertyConditionSupportedOperatorEquals),
								string(automationrules.AutomationRulePropertyConditionSupportedOperatorNotContains),
								string(automationrules.AutomationRulePropertyConditionSupportedOperatorNotEndsWith),
								string(automationrules.AutomationRulePropertyConditionSupportedOperatorNotEquals),
								string(automationrules.AutomationRulePropertyConditionSupportedOperatorNotStartsWith),
								string(automationrules.AutomationRulePropertyConditionSupportedOperatorStartsWith),
							}, false),
						},

						"values": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"action_incident": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"order": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"status": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(automationrules.IncidentStatusActive),
								string(automationrules.IncidentStatusClosed),
								string(automationrules.IncidentStatusNew),
							}, false),
						},

						"classification": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(automationrules.IncidentClassificationUndetermined),
								string(automationrules.IncidentClassificationBenignPositive) + "_" + string(automationrules.IncidentClassificationReasonSuspiciousButExpected),
								string(automationrules.IncidentClassificationFalsePositive) + "_" + string(automationrules.IncidentClassificationReasonIncorrectAlertLogic),
								string(automationrules.IncidentClassificationFalsePositive) + "_" + string(automationrules.IncidentClassificationReasonInaccurateData),
								string(automationrules.IncidentClassificationTruePositive) + "_" + string(automationrules.IncidentClassificationReasonSuspiciousActivity),
							}, false),
						},

						"classification_comment": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"labels": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"owner_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"severity": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(automationrules.IncidentSeverityHigh),
								string(automationrules.IncidentSeverityInformational),
								string(automationrules.IncidentSeverityLow),
								string(automationrules.IncidentSeverityMedium),
							}, false),
						},
					},
				},
				AtLeastOneOf: []string{"action_incident", "action_playbook"},
			},

			"action_playbook": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"order": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"logic_app_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"tenant_id": {
							Type: pluginsdk.TypeString,
							// We'll use the current tenant id if this property is absent.
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IsUUID,
						},
					},
				},
				AtLeastOneOf: []string{"action_incident", "action_playbook"},
			},
		},
	}
}

func resourceSentinelAutomationRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AutomationRulesClient
	tenantId := meta.(*clients.Client).Account.TenantId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	workspaceId, err := automationrules.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	id := automationrules.NewAutomationRuleID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, name)

	if d.IsNewResource() {
		resp, err := client.AutomationRulesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_sentinel_automation_rule", id.ID())
		}
	}

	actions, err := expandAutomationRuleActions(d, tenantId)
	if err != nil {
		return err
	}
	params := automationrules.AutomationRule{
		Properties: automationrules.AutomationRuleProperties{
			DisplayName: d.Get("display_name").(string),
			Order:       int64(d.Get("order").(int)),
			TriggeringLogic: automationrules.AutomationRuleTriggeringLogic{
				IsEnabled:    d.Get("enabled").(bool),
				TriggersOn:   "Incidents", // This is the only supported enum for now. The reason why there is no enum in SDK, see: https://github.com/Azure/azure-sdk-for-go/issues/14589
				TriggersWhen: "Created",   // This is the only supported enum for now. The reason why there is no enum in SDK, see: https://github.com/Azure/azure-sdk-for-go/issues/14589
				Conditions:   expandAutomationRuleConditions(d.Get("condition").([]interface{})),
			},
			Actions: *actions,
		},
	}

	if expiration := d.Get("expiration").(string); expiration != "" {
		t, _ := time.Parse(time.RFC3339, expiration)
		params.Properties.TriggeringLogic.ExpirationTimeUtc = utils.String(date.Time{Time: t}.String())
	}

	_, err = client.AutomationRulesCreateOrUpdate(ctx, id, params)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelAutomationRuleRead(d, meta)
}

func resourceSentinelAutomationRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AutomationRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := automationrules.ParseAutomationRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.AutomationRulesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("retrieving %s: model is nil", id)
	}

	d.Set("name", model.Name)
	d.Set("log_analytics_workspace_id", loganalyticsParse.NewLogAnalyticsWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID())

	prop := model.Properties
	d.Set("display_name", prop.DisplayName)

	var order int
	order = int(prop.Order)
	d.Set("order", order)

	tl := prop.TriggeringLogic
	enabled := tl.IsEnabled

	d.Set("enabled", enabled)

	var expiration string
	if tl.ExpirationTimeUtc != nil {
		t, err := time.Parse(time.RFC3339, *tl.ExpirationTimeUtc)
		if err != nil {
			expiration = t.Format(time.RFC3339)
		}
	}

	d.Set("expiration", expiration)

	if err := d.Set("condition", flattenAutomationRuleConditions(tl.Conditions)); err != nil {
		return fmt.Errorf("setting `condition`: %v", err)
	}

	actionIncident, actionPlaybook := flattenAutomationRuleActions(&prop.Actions)

	if err := d.Set("action_incident", actionIncident); err != nil {
		return fmt.Errorf("setting `action_incident`: %v", err)
	}
	if err := d.Set("action_playbook", actionPlaybook); err != nil {
		return fmt.Errorf("setting `action_playbook`: %v", err)
	}

	return nil
}

func resourceSentinelAutomationRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AutomationRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := automationrules.ParseAutomationRuleID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.AutomationRulesDelete(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandAutomationRuleConditions(input []interface{}) *[]automationrules.AutomationRuleCondition {
	if len(input) == 0 {
		return nil
	}

	out := make([]automationrules.AutomationRuleCondition, 0, len(input))
	for _, b := range input {
		b := b.(map[string]interface{})
		operator := automationrules.AutomationRulePropertyConditionSupportedOperator(b["operator"].(string))
		name := automationrules.AutomationRulePropertyConditionSupportedProperty(b["property"].(string))
		out = append(out, &automationrules.AutomationRulePropertyValuesCondition{
			Operator:       &operator,
			PropertyName:   &name,
			PropertyValues: utils.ExpandStringSlice(b["values"].([]interface{})),
		})
	}
	return &out
}

func flattenAutomationRuleConditions(conditions *[]automationrules.AutomationRuleCondition) interface{} {
	if conditions == nil {
		return nil
	}

	out := make([]interface{}, 0, len(*conditions))
	for _, condition := range *conditions {
		condition := condition.(automationrules.AutomationRulePropertyValuesCondition)

		var (
			property string
			operator string
			values   []interface{}
		)

		if condition.PropertyName != nil {
			property = string(*condition.PropertyName)
		}

		if condition.Operator != nil {
			operator = string(*condition.Operator)
		}

		if condition.PropertyValues != nil {
			values = utils.FlattenStringSlice(condition.PropertyValues)
		}

		out = append(out, map[string]interface{}{
			"property": property,
			"operator": operator,
			"values":   values,
		})
	}
	return out
}

func expandAutomationRuleActions(d *pluginsdk.ResourceData, defaultTenantId string) (*[]automationrules.AutomationRuleAction, error) {
	actionIncident, err := expandAutomationRuleActionIncident(d.Get("action_incident").([]interface{}))
	if err != nil {
		return nil, err
	}
	actionPlaybook := expandAutomationRuleActionPlaybook(d.Get("action_playbook").([]interface{}), defaultTenantId)

	if len(actionIncident)+len(actionPlaybook) == 0 {
		return nil, nil
	}

	out := make([]automationrules.AutomationRuleAction, 0, len(actionIncident)+len(actionPlaybook))
	out = append(out, actionIncident...)
	out = append(out, actionPlaybook...)
	return &out, nil
}

func flattenAutomationRuleActions(input *[]automationrules.AutomationRuleAction) (actionIncident []interface{}, actionPlaybook []interface{}) {
	if input == nil {
		return nil, nil
	}

	actionIncident = make([]interface{}, 0)
	actionPlaybook = make([]interface{}, 0)

	for _, action := range *input {
		switch action := action.(type) {
		case automationrules.AutomationRuleModifyPropertiesAction:
			actionIncident = append(actionIncident, flattenAutomationRuleActionIncident(action))
		case automationrules.AutomationRuleRunPlaybookAction:
			actionPlaybook = append(actionPlaybook, flattenAutomationRuleActionPlaybook(action))
		}
	}

	return
}

func expandAutomationRuleActionIncident(input []interface{}) ([]automationrules.AutomationRuleAction, error) {
	if len(input) == 0 {
		return nil, nil
	}

	out := make([]automationrules.AutomationRuleAction, 0, len(input))
	for _, b := range input {
		b := b.(map[string]interface{})

		status := automationrules.IncidentStatus(b["status"].(string))
		l := strings.Split(b["classification"].(string), "_")
		classification, clr := l[0], ""
		if len(l) == 2 {
			clr = l[1]
		}
		classificationComment := b["classification_comment"].(string)

		// sanity check on classification
		if status == automationrules.IncidentStatusClosed && classification == "" {
			return nil, fmt.Errorf("`classification` is required when `status` is set to `Closed`")
		}
		if status != automationrules.IncidentStatusClosed {
			if classification != "" {
				return nil, fmt.Errorf("`classification` can't be set when `status` is not set to `Closed`")
			}
			if classificationComment != "" {
				return nil, fmt.Errorf("`classification_comment` can't be set when `status` is not set to `Closed`")
			}
		}

		var labelsPtr *[]automationrules.IncidentLabel
		if labelStrsPtr := utils.ExpandStringSlice(b["labels"].([]interface{})); labelStrsPtr != nil && len(*labelStrsPtr) > 0 {
			labels := make([]automationrules.IncidentLabel, 0, len(*labelStrsPtr))
			for _, label := range *labelStrsPtr {
				labels = append(labels, automationrules.IncidentLabel{
					LabelName: label,
				})
			}
			labelsPtr = &labels
		}

		var ownerPtr *automationrules.IncidentOwnerInfo
		if ownerIdStr := b["owner_id"].(string); ownerIdStr != "" {
			ownerId, err := uuid.FromString(ownerIdStr)
			if err != nil {
				return nil, fmt.Errorf("getting `owner_id`: %v", err)
			}
			ownerPtr = &automationrules.IncidentOwnerInfo{
				ObjectId: utils.String(ownerId.String()),
			}
		}

		severity := b["severity"].(string)

		// sanity check on the whole incident action
		if severity == "" && ownerPtr == nil && labelsPtr == nil && status == "" {
			return nil, fmt.Errorf("at least one of `severity`, `owner_id`, `labels` or `status` should be specified")
		}

		tCls := automationrules.IncidentClassification(classification)
		tClr := automationrules.IncidentClassificationReason(clr)
		tSvy := automationrules.IncidentSeverity(severity)

		out = append(out, automationrules.AutomationRuleModifyPropertiesAction{
			Order: int64(b["order"].(int)),
			ActionConfiguration: &automationrules.IncidentPropertiesAction{
				Status:                &status,
				Classification:        &tCls,
				ClassificationComment: &classificationComment,
				ClassificationReason:  &tClr,
				Labels:                labelsPtr,
				Owner:                 ownerPtr,
				Severity:              &tSvy,
			},
		})
	}

	return out, nil
}

func flattenAutomationRuleActionIncident(input automationrules.AutomationRuleModifyPropertiesAction) map[string]interface{} {
	order := int(input.Order)

	var (
		status      string
		clsf        string
		clsfComment string
		clsfReason  string
		labels      []interface{}
		owner       string
		severity    string
	)

	if cfg := input.ActionConfiguration; cfg != nil {
		status = string(*cfg.Status)
		clsf = string(*cfg.Classification)
		if cfg.ClassificationComment != nil {
			clsfComment = *cfg.ClassificationComment
		}
		clsfReason = string(*cfg.ClassificationReason)

		if cfg.Labels != nil {
			for _, label := range *cfg.Labels {
				labels = append(labels, label.LabelName)
			}
		}

		if cfg.Owner != nil && cfg.Owner.ObjectId != nil {
			owner = *cfg.Owner.ObjectId
		}

		severity = string(*cfg.Severity)
	}

	classification := clsf
	if clsfReason != "" {
		classification = classification + "_" + clsfReason
	}

	return map[string]interface{}{
		"order":                  order,
		"status":                 status,
		"classification":         classification,
		"classification_comment": clsfComment,
		"labels":                 labels,
		"owner_id":               owner,
		"severity":               severity,
	}
}

func expandAutomationRuleActionPlaybook(input []interface{}, defaultTenantId string) []automationrules.AutomationRuleAction {
	if len(input) == 0 {
		return nil
	}

	out := make([]automationrules.AutomationRuleAction, 0, len(input))
	for _, b := range input {
		b := b.(map[string]interface{})

		tenantId := defaultTenantId
		if tid := b["tenant_id"].(string); tid != "" {
			tenantId = tid
		}

		out = append(out, automationrules.AutomationRuleRunPlaybookAction{
			Order: int64(b["order"].(int)),
			ActionConfiguration: &automationrules.PlaybookActionProperties{
				LogicAppResourceId: utils.String(b["logic_app_id"].(string)),
				TenantId:           &tenantId,
			},
		})
	}
	return out
}

func flattenAutomationRuleActionPlaybook(input automationrules.AutomationRuleRunPlaybookAction) map[string]interface{} {
	order := int(input.Order)

	var (
		logicAppId string
		tenantId   string
	)

	if cfg := input.ActionConfiguration; cfg != nil {
		if cfg.LogicAppResourceId != nil {
			logicAppId = *cfg.LogicAppResourceId
		}

		if cfg.TenantId != nil {
			tenantId = *cfg.TenantId
		}
	}

	return map[string]interface{}{
		"order":        order,
		"logic_app_id": logicAppId,
		"tenant_id":    tenantId,
	}
}
