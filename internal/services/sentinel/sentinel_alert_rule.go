package sentinel

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-07-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func alertRuleID(rule alertrules.AlertRule) *string {
	if rule == nil {
		return nil
	}
	switch rule := rule.(type) {
	case alertrules.FusionAlertRule:
		return rule.Id
	case alertrules.MicrosoftSecurityIncidentCreationAlertRule:
		return rule.Id
	case alertrules.ScheduledAlertRule:
		return rule.Id
	case alertrules.MLBehaviorAnalyticsAlertRule:
		return rule.Id
	default:
		return nil
	}
}

func importSentinelAlertRule(expectKind alertrules.AlertRuleKind) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := alertrules.ParseAlertRuleID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).Sentinel.AlertRulesClient
		resp, err := client.AlertRulesGet(ctx, *id)
		if err != nil {
			return nil, fmt.Errorf("retrieving Sentinel Alert Rule %q: %+v", id, err)
		}

		model := *resp.Model
		if model == nil {
			return nil, fmt.Errorf("retrieving Sentinel Alert Rule %q: model is nil", id)
		}

		if err := assertAlertRuleKind(model, expectKind); err != nil {
			return nil, err
		}
		return []*pluginsdk.ResourceData{d}, nil
	}
}

func assertAlertRuleKind(rule alertrules.AlertRule, expectKind alertrules.AlertRuleKind) error {
	var kind alertrules.AlertRuleKind
	switch rule.(type) {
	case alertrules.MLBehaviorAnalyticsAlertRule:
		kind = alertrules.AlertRuleKindMLBehaviorAnalytics
	case alertrules.FusionAlertRule:
		kind = alertrules.AlertRuleKindFusion
	case alertrules.MicrosoftSecurityIncidentCreationAlertRule:
		kind = alertrules.AlertRuleKindMicrosoftSecurityIncidentCreation
	case alertrules.ScheduledAlertRule:
		kind = alertrules.AlertRuleKindScheduled
	}
	if expectKind != kind {
		return fmt.Errorf("Sentinel Alert Rule has mismatched kind, expected: %q, got %q", expectKind, kind)
	}
	return nil
}
