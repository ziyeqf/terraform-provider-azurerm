package sentinel

import (
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-07-01-preview/alertrules"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/rickb777/date/period"
)

func resourceSentinelAlertRuleScheduled() *pluginsdk.Resource {
	var entityMappingTypes = []string{
		string(alertrules.EntityMappingTypeAccount),
		string(alertrules.EntityMappingTypeAzureResource),
		string(alertrules.EntityMappingTypeCloudApplication),
		string(alertrules.EntityMappingTypeDNS),
		string(alertrules.EntityMappingTypeFile),
		string(alertrules.EntityMappingTypeFileHash),
		string(alertrules.EntityMappingTypeHost),
		string(alertrules.EntityMappingTypeIP),
		string(alertrules.EntityMappingTypeMailbox),
		string(alertrules.EntityMappingTypeMailCluster),
		string(alertrules.EntityMappingTypeMailMessage),
		string(alertrules.EntityMappingTypeMalware),
		string(alertrules.EntityMappingTypeProcess),
		string(alertrules.EntityMappingTypeRegistryKey),
		string(alertrules.EntityMappingTypeRegistryValue),
		string(alertrules.EntityMappingTypeSecurityGroup),
		string(alertrules.EntityMappingTypeSubmissionMail),
		string(alertrules.EntityMappingTypeURL),
	}
	return &pluginsdk.Resource{
		Create: resourceSentinelAlertRuleScheduledCreateUpdate,
		Read:   resourceSentinelAlertRuleScheduledRead,
		Update: resourceSentinelAlertRuleScheduledCreateUpdate,
		Delete: resourceSentinelAlertRuleScheduledDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.AlertRuleID(id)
			return err
		}, importSentinelAlertRule(alertrules.AlertRuleKindScheduled)),

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

			"display_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"alert_rule_template_guid": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"alert_rule_template_version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"event_grouping": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"aggregation_method": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(alertrules.EventGroupingAggregationKindAlertPerResult),
								string(alertrules.EventGroupingAggregationKindSingleAlert),
							}, false),
						},
					},
				},
			},

			"tactics": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(alertrules.AttackTacticCollection),
						string(alertrules.AttackTacticCommandAndControl),
						string(alertrules.AttackTacticCredentialAccess),
						string(alertrules.AttackTacticDefenseEvasion),
						string(alertrules.AttackTacticDiscovery),
						string(alertrules.AttackTacticExecution),
						string(alertrules.AttackTacticExfiltration),
						string(alertrules.AttackTacticImpact),
						string(alertrules.AttackTacticInitialAccess),
						string(alertrules.AttackTacticLateralMovement),
						string(alertrules.AttackTacticPersistence),
						string(alertrules.AttackTacticPrivilegeEscalation),
						string(alertrules.AttackTacticPreAttack),
					}, false),
				},
			},

			"incident_configuration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"create_incident": {
							Required: true,
							Type:     pluginsdk.TypeBool,
						},
						"grouping": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MaxItems: 1,
							MinItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  true,
									},
									"lookback_duration": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validate.ISO8601Duration,
										Default:      "PT5M",
									},
									"reopen_closed_incidents": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
									},
									"entity_matching_method": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  alertrules.MatchingMethodAnyAlert,
										ValidateFunc: validation.StringInSlice([]string{
											string(alertrules.MatchingMethodAnyAlert),
											string(alertrules.MatchingMethodSelected),
											string(alertrules.MatchingMethodAllEntities),
										}, false),
									},
									"group_by_entities": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice(entityMappingTypes, false),
										},
									},
									"group_by_alert_details": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												string(alertrules.AlertDetailDisplayName),
												string(alertrules.AlertDetailSeverity),
											},
												false),
										},
									},
									"group_by_custom_details": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
								},
							},
						},
					},
				},
			},

			"severity": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(alertrules.AlertSeverityHigh),
					string(alertrules.AlertSeverityMedium),
					string(alertrules.AlertSeverityLow),
					string(alertrules.AlertSeverityInformational),
				}, false),
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"query": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"query_frequency": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      "PT5H",
				ValidateFunc: validate.ISO8601DurationBetween("PT5M", "P14D"),
			},

			"query_period": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      "PT5H",
				ValidateFunc: validate.ISO8601DurationBetween("PT5M", "P14D"),
			},

			"trigger_operator": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(alertrules.TriggerOperatorGreaterThan),
				ValidateFunc: validation.StringInSlice([]string{
					string(alertrules.TriggerOperatorGreaterThan),
					string(alertrules.TriggerOperatorLessThan),
					string(alertrules.TriggerOperatorEqual),
					string(alertrules.TriggerOperatorNotEqual),
				}, false),
			},

			"trigger_threshold": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"suppression_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},
			"suppression_duration": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      "PT5H",
				ValidateFunc: validate.ISO8601DurationBetween("PT5M", "PT24H"),
			},
			"alert_details_override": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"description_format": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"display_name_format": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"severity_column_name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"tactics_column_name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
			"custom_details": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
			"entity_mapping": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 5,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"entity_type": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(entityMappingTypes, false),
						},
						"field_mapping": {
							Type:     pluginsdk.TypeList,
							MaxItems: 3,
							Required: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"identifier": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"column_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceSentinelAlertRuleScheduledCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	workspaceID, err := alertrules.ParseAlertRuleID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	id := alertrules.NewAlertRuleID(workspaceID.SubscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)

	if d.IsNewResource() {
		resp, err := client.AlertRulesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existing Sentinel Alert Rule Scheduled %q: %+v", id, err)
			}
		}

		model := resp.Model
		if model == nil {
			return fmt.Errorf("existing Sentinel Alert Rule Scheduled %q not found", id)
		}

		id := alertRuleID(model)
		if id != nil && *id != "" {
			return tf.ImportAsExistsError("azurerm_sentinel_alert_rule_scheduled", *id)
		}
	}

	// Sanity checks

	// query frequency must <= query period: ensure there is no gaps in the overall query coverage.
	queryFreq := d.Get("query_frequency").(string)
	queryFreqDuration := period.MustParse(queryFreq).DurationApprox()

	queryPeriod := d.Get("query_period").(string)
	queryPeriodDuration := period.MustParse(queryPeriod).DurationApprox()
	if queryFreqDuration > queryPeriodDuration {
		return fmt.Errorf("`query_frequency`(%v) should not be larger than `query period`(%v), which introduce gaps in the overall query coverage", queryFreq, queryPeriod)
	}

	// query frequency must <= suppression duration: otherwise suppression has no effect.
	suppressionDuration := d.Get("suppression_duration").(string)
	suppressionEnabled := d.Get("suppression_enabled").(bool)
	if suppressionEnabled {
		suppressionDurationDuration := period.MustParse(suppressionDuration).DurationApprox()
		if queryFreqDuration > suppressionDurationDuration {
			return fmt.Errorf("`query_frequency`(%v) should not be larger than `suppression_duration`(%v), which makes suppression pointless", queryFreq, suppressionDuration)
		}
	}

	param := alertrules.ScheduledAlertRule{
		Properties: &alertrules.ScheduledAlertRuleProperties{
			Description:           utils.String(d.Get("description").(string)),
			DisplayName:           d.Get("display_name").(string),
			Tactics:               expandAlertRuleScheduledTactics(d.Get("tactics").(*pluginsdk.Set).List()),
			IncidentConfiguration: expandAlertRuleScheduledIncidentConfiguration(d.Get("incident_configuration").([]interface{})),
			Severity:              alertrules.AlertSeverity(d.Get("severity").(string)),
			Enabled:               d.Get("enabled").(bool),
			Query:                 d.Get("query").(string),
			QueryFrequency:        queryFreq,
			QueryPeriod:           queryPeriod,
			SuppressionEnabled:    suppressionEnabled,
			SuppressionDuration:   suppressionDuration,
			TriggerOperator:       alertrules.TriggerOperator(d.Get("trigger_operator").(string)),
			TriggerThreshold:      int64(d.Get("trigger_threshold").(int)),
		},
	}

	if v, ok := d.GetOk("alert_rule_template_guid"); ok {
		param.Properties.AlertRuleTemplateName = utils.String(v.(string))
	}
	if v, ok := d.GetOk("alert_rule_template_version"); ok {
		param.Properties.TemplateVersion = utils.String(v.(string))
	}
	if v, ok := d.GetOk("event_grouping"); ok {
		param.Properties.EventGroupingSettings = expandAlertRuleScheduledEventGroupingSetting(v.([]interface{}))
	}
	if v, ok := d.GetOk("alert_details_override"); ok {
		param.Properties.AlertDetailsOverride = expandAlertRuleScheduledAlertDetailsOverride(v.([]interface{}))
	}
	if v, ok := d.GetOk("custom_details"); ok {
		tmp := v.(map[string]string)
		param.Properties.CustomDetails = &tmp
	}
	if v, ok := d.GetOk("entity_mapping"); ok {
		param.Properties.EntityMappings = expandAlertRuleScheduledEntityMapping(v.([]interface{}))
	}

	if !d.IsNewResource() {
		resp, err := client.AlertRulesGet(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving Sentinel Alert Rule Scheduled %q: %+v", id, err)
		}

		model := *resp.Model
		if model == nil {
			return fmt.Errorf("existing Sentinel Alert Rule Scheduled %q not found", id)
		}

		if err := assertAlertRuleKind(model, alertrules.AlertRuleKindScheduled); err != nil {
			return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
		}
	}

	if _, err := client.AlertRulesCreateOrUpdate(ctx, id, param); err != nil {
		return fmt.Errorf("creating Sentinel Alert Rule Scheduled %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelAlertRuleScheduledRead(d, meta)
}

func resourceSentinelAlertRuleScheduledRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] Sentinel Alert Rule Scheduled %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Sentinel Alert Rule Scheduled %q: %+v", id, err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("existing Sentinel Alert Rule Scheduled %q not found", id)
	}

	if err := assertAlertRuleKind(*model, alertrules.AlertRuleKindScheduled); err != nil {
		return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
	}
	rule := (*model).(alertrules.ScheduledAlertRule)

	d.Set("name", rule.Name)

	workspaceId := alertrules.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)
	d.Set("log_analytics_workspace_id", workspaceId.ID())

	if prop := rule.Properties; prop != nil {
		d.Set("description", prop.Description)
		d.Set("display_name", prop.DisplayName)
		if err := d.Set("tactics", flattenAlertRuleScheduledTactics(prop.Tactics)); err != nil {
			return fmt.Errorf("setting `tactics`: %+v", err)
		}
		if err := d.Set("incident_configuration", flattenAlertRuleScheduledIncidentConfiguration(prop.IncidentConfiguration)); err != nil {
			return fmt.Errorf("setting `incident_configuration`: %+v", err)
		}
		d.Set("severity", string(prop.Severity))
		d.Set("enabled", prop.Enabled)
		d.Set("query", prop.Query)
		d.Set("query_frequency", prop.QueryFrequency)
		d.Set("query_period", prop.QueryPeriod)
		d.Set("trigger_operator", string(prop.TriggerOperator))
		d.Set("trigger_threshold", int(prop.TriggerThreshold))
		d.Set("suppression_enabled", prop.SuppressionEnabled)
		d.Set("suppression_duration", prop.SuppressionDuration)
		d.Set("alert_rule_template_guid", prop.AlertRuleTemplateName)
		d.Set("alert_rule_template_version", prop.TemplateVersion)

		if err := d.Set("event_grouping", flattenAlertRuleScheduledEventGroupingSetting(prop.EventGroupingSettings)); err != nil {
			return fmt.Errorf("setting `event_grouping`: %+v", err)
		}
		if err := d.Set("alert_details_override", flattenAlertRuleScheduledAlertDetailsOverride(prop.AlertDetailsOverride)); err != nil {
			return fmt.Errorf("setting `alert_details_override`: %+v", err)
		}
		if err := d.Set("custom_details", prop.CustomDetails); err != nil {
			return fmt.Errorf("setting `custom_details`: %+v", err)
		}
		if err := d.Set("entity_mapping", flattenAlertRuleScheduledEntityMapping(prop.EntityMappings)); err != nil {
			return fmt.Errorf("setting `entity_mapping`: %+v", err)
		}
	}

	return nil
}

func resourceSentinelAlertRuleScheduledDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := alertrules.ParseAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.AlertRulesDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting Sentinel Alert Rule Scheduled %q: %+v", id, err)
	}

	return nil
}

func expandAlertRuleScheduledTactics(input []interface{}) *[]alertrules.AttackTactic {
	result := make([]alertrules.AttackTactic, 0)

	for _, e := range input {
		result = append(result, alertrules.AttackTactic(e.(string)))
	}

	return &result
}

func flattenAlertRuleScheduledTactics(input *[]alertrules.AttackTactic) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		output = append(output, string(e))
	}

	return output
}

func expandAlertRuleScheduledIncidentConfiguration(input []interface{}) *alertrules.IncidentConfiguration {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})

	output := &alertrules.IncidentConfiguration{
		CreateIncident:        raw["create_incident"].(bool),
		GroupingConfiguration: expandAlertRuleScheduledGrouping(raw["grouping"].([]interface{})),
	}

	return output
}

func flattenAlertRuleScheduledIncidentConfiguration(input *alertrules.IncidentConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	createIncident := false
	createIncident = input.CreateIncident

	return []interface{}{
		map[string]interface{}{
			"create_incident": createIncident,
			"grouping":        flattenAlertRuleScheduledGrouping(input.GroupingConfiguration),
		},
	}
}

func expandAlertRuleScheduledGrouping(input []interface{}) *alertrules.GroupingConfiguration {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})

	output := &alertrules.GroupingConfiguration{
		Enabled:              raw["enabled"].(bool),
		ReopenClosedIncident: raw["reopen_closed_incidents"].(bool),
		LookbackDuration:     raw["lookback_duration"].(string),
		MatchingMethod:       alertrules.MatchingMethod(raw["entity_matching_method"].(string)),
	}

	groupByEntitiesList := raw["group_by_entities"].([]interface{})
	groupByEntities := make([]alertrules.EntityMappingType, len(groupByEntitiesList))
	for idx, t := range groupByEntitiesList {
		groupByEntities[idx] = alertrules.EntityMappingType(t.(string))
	}
	output.GroupByEntities = &groupByEntities

	groupByAlertDetailsList := raw["group_by_alert_details"].([]interface{})
	groupByAlertDetails := make([]alertrules.AlertDetail, len(groupByAlertDetailsList))
	for idx, t := range groupByAlertDetailsList {
		groupByAlertDetails[idx] = alertrules.AlertDetail(t.(string))
	}
	output.GroupByAlertDetails = &groupByAlertDetails

	output.GroupByCustomDetails = utils.ExpandStringSlice(raw["group_by_custom_details"].([]interface{}))

	return output
}

func expandAlertRuleScheduledEventGroupingSetting(input []interface{}) *alertrules.EventGroupingSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	result := alertrules.EventGroupingSettings{}

	if aggregationKind := v["aggregation_method"].(string); aggregationKind != "" {
		t := alertrules.EventGroupingAggregationKind(aggregationKind)
		result.AggregationKind = &t
	}

	return &result
}

func flattenAlertRuleScheduledGrouping(input *alertrules.GroupingConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	enabled := false
	enabled = input.Enabled

	lookbackDuration := ""
	lookbackDuration = input.LookbackDuration

	reopenClosedIncidents := false
	reopenClosedIncidents = input.ReopenClosedIncident

	var groupByEntities []interface{}
	if input.GroupByEntities != nil {
		for _, entity := range *input.GroupByEntities {
			groupByEntities = append(groupByEntities, string(entity))
		}
	}

	var groupByAlertDetails []interface{}
	if input.GroupByAlertDetails != nil {
		for _, detail := range *input.GroupByAlertDetails {
			groupByAlertDetails = append(groupByAlertDetails, string(detail))
		}
	}

	var groupByCustomDetails []interface{}
	if input.GroupByCustomDetails != nil {
		for _, detail := range *input.GroupByCustomDetails {
			groupByCustomDetails = append(groupByCustomDetails, detail)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"enabled":                 enabled,
			"lookback_duration":       lookbackDuration,
			"reopen_closed_incidents": reopenClosedIncidents,
			"entity_matching_method":  string(input.MatchingMethod),
			"group_by_entities":       groupByEntities,
			"group_by_alert_details":  groupByAlertDetails,
			"group_by_custom_details": groupByCustomDetails,
		},
	}
}

func flattenAlertRuleScheduledEventGroupingSetting(input *alertrules.EventGroupingSettings) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var aggregationKind string
	if *input.AggregationKind != "" {
		aggregationKind = string(*input.AggregationKind)
	}

	return []interface{}{
		map[string]interface{}{
			"aggregation_method": aggregationKind,
		},
	}
}

func expandAlertRuleScheduledAlertDetailsOverride(input []interface{}) *alertrules.AlertDetailsOverride {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	b := input[0].(map[string]interface{})
	output := &alertrules.AlertDetailsOverride{}

	if v := b["description_format"]; v != "" {
		output.AlertDescriptionFormat = utils.String(v.(string))
	}
	if v := b["display_name_format"]; v != "" {
		output.AlertDisplayNameFormat = utils.String(v.(string))
	}
	if v := b["severity_column_name"]; v != "" {
		output.AlertSeverityColumnName = utils.String(v.(string))
	}
	if v := b["tactics_column_name"]; v != "" {
		output.AlertTacticsColumnName = utils.String(v.(string))
	}

	return output
}

func flattenAlertRuleScheduledAlertDetailsOverride(input *alertrules.AlertDetailsOverride) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var descriptionFormat string
	if input.AlertDescriptionFormat != nil {
		descriptionFormat = *input.AlertDescriptionFormat
	}

	var displayNameFormat string
	if input.AlertDisplayNameFormat != nil {
		displayNameFormat = *input.AlertDisplayNameFormat
	}

	var severityColumnName string
	if input.AlertSeverityColumnName != nil {
		severityColumnName = *input.AlertSeverityColumnName
	}

	var tacticsColumnName string
	if input.AlertTacticsColumnName != nil {
		tacticsColumnName = *input.AlertTacticsColumnName
	}

	return []interface{}{
		map[string]interface{}{
			"description_format":   descriptionFormat,
			"display_name_format":  displayNameFormat,
			"severity_column_name": severityColumnName,
			"tactics_column_name":  tacticsColumnName,
		},
	}
}

func expandAlertRuleScheduledEntityMapping(input []interface{}) *[]alertrules.EntityMapping {
	if len(input) == 0 {
		return nil
	}

	result := make([]alertrules.EntityMapping, 0)

	for _, e := range input {
		b := e.(map[string]interface{})
		t := alertrules.EntityMappingType(b["entity_type"].(string))
		result = append(result, alertrules.EntityMapping{
			EntityType:    &t,
			FieldMappings: expandAlertRuleScheduledFieldMapping(b["field_mapping"].([]interface{})),
		})
	}

	return &result
}

func flattenAlertRuleScheduledEntityMapping(input *[]alertrules.EntityMapping) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		output = append(output, map[string]interface{}{
			"entity_type":   e.EntityType,
			"field_mapping": flattenAlertRuleScheduledFieldMapping(e.FieldMappings),
		})
	}

	return output
}

func expandAlertRuleScheduledFieldMapping(input []interface{}) *[]alertrules.FieldMapping {
	if len(input) == 0 {
		return nil
	}

	result := make([]alertrules.FieldMapping, 0)

	for _, e := range input {
		b := e.(map[string]interface{})
		result = append(result, alertrules.FieldMapping{
			Identifier: utils.String(b["identifier"].(string)),
			ColumnName: utils.String(b["column_name"].(string)),
		})
	}

	return &result
}

func flattenAlertRuleScheduledFieldMapping(input *[]alertrules.FieldMapping) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		var identifier string
		if e.Identifier != nil {
			identifier = *e.Identifier
		}

		var columnName string
		if e.ColumnName != nil {
			columnName = *e.ColumnName
		}

		output = append(output, map[string]interface{}{
			"identifier":  identifier,
			"column_name": columnName,
		})
	}

	return output
}
