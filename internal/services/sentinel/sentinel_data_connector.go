package sentinel

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-07-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func importSentinelDataConnector(expectKind dataconnectors.DataConnectorKind) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := dataconnectors.ParseDataConnectorID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).Sentinel.DataConnectorsClient
		resp, err := client.DataConnectorsGet(ctx, *id)
		if err != nil {
			return nil, fmt.Errorf("retrieving Sentinel Alert Rule %q: %+v", id, err)
		}

		if err := assertDataConnectorKind(resp.Model, expectKind); err != nil {
			return nil, err
		}
		return []*pluginsdk.ResourceData{d}, nil
	}
}

func assertDataConnectorKind(dc dataconnectors.DataConnector, expectKind dataconnectors.DataConnectorKind) error {
	var kind dataconnectors.DataConnectorKind
	switch dc.(type) {
	case dataconnectors.AADDataConnector:
		kind = dataconnectors.DataConnectorKindAzureActiveDirectory
	case dataconnectors.AATPDataConnector:
		kind = dataconnectors.DataConnectorKindAzureAdvancedThreatProtection
	case dataconnectors.ASCDataConnector:
		kind = dataconnectors.DataConnectorKindAzureSecurityCenter
	case dataconnectors.MCASDataConnector:
		kind = dataconnectors.DataConnectorKindMicrosoftCloudAppSecurity
	case dataconnectors.TIDataConnector:
		kind = dataconnectors.DataConnectorKindThreatIntelligence
	case dataconnectors.OfficeDataConnector:
		kind = dataconnectors.DataConnectorKindOfficeThreeSixFive
	case dataconnectors.OfficeATPDataConnector:
		kind = dataconnectors.DataConnectorKindOfficeATP
	case dataconnectors.AwsCloudTrailDataConnector:
		kind = dataconnectors.DataConnectorKindAmazonWebServicesCloudTrail
	case dataconnectors.MDATPDataConnector:
		kind = dataconnectors.DataConnectorKindMicrosoftDefenderAdvancedThreatProtection
	}
	if expectKind != kind {
		return fmt.Errorf("Sentinel Data Connector has mismatched kind, expected: %q, got %q", expectKind, kind)
	}
	return nil
}
