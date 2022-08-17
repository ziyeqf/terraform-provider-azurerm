package sentinel

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-07-01-preview/dataconnectors"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2021-09-01-preview/securityinsight"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	loganalyticsValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DataConnectorAwsS3Resource struct{}

var _ sdk.ResourceWithUpdate = DataConnectorAwsS3Resource{}
var _ sdk.ResourceWithCustomImporter = DataConnectorAwsS3Resource{}

type DataConnectorAwsS3Model struct {
	Name                    string   `tfschema:"name"`
	LogAnalyticsWorkspaceId string   `tfschema:"log_analytics_workspace_id"`
	AwsRoleArm              string   `tfschema:"aws_role_arn"`
	DestinationTable        string   `tfschema:"destination_table"`
	SqsUrls                 []string `tfschema:"sqs_urls"`
}

func (r DataConnectorAwsS3Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
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
			ValidateFunc: loganalyticsValidate.LogAnalyticsWorkspaceID,
		},

		"aws_role_arn": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.IsARN,
		},

		"destination_table": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"sqs_urls": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func (r DataConnectorAwsS3Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DataConnectorAwsS3Resource) ResourceType() string {
	return "azurerm_sentinel_data_connector_aws_s3"
}

func (r DataConnectorAwsS3Resource) ModelObject() interface{} {
	return &DataConnectorAwsS3Model{}
}

func (r DataConnectorAwsS3Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.DataConnectorID
}

func (r DataConnectorAwsS3Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		_, err := importSentinelDataConnector(dataconnectors.DataConnectorKindAmazonWebServicesSThree)(ctx, metadata.ResourceData, metadata.Client)
		return err
	}
}

func (r DataConnectorAwsS3Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient

			var plan DataConnectorAwsS3Model
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			workspaceId, err := dataconnectors.ParseWorkspaceID(plan.LogAnalyticsWorkspaceId)
			if err != nil {
				return err
			}

			id := dataconnectors.NewDataConnectorID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, plan.Name)
			existing, err := client.DataConnectorsGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			params := securityinsight.AwsS3DataConnector{
				Name: &plan.Name,
				AwsS3DataConnectorProperties: &securityinsight.AwsS3DataConnectorProperties{
					DestinationTable: utils.String(plan.DestinationTable),
					SqsUrls:          &plan.SqsUrls,
					RoleArn:          utils.String(plan.AwsRoleArm),
					DataTypes: &securityinsight.AwsS3DataConnectorDataTypes{
						Logs: &securityinsight.AwsS3DataConnectorDataTypesLogs{
							State: securityinsight.DataTypeStateEnabled,
						},
					},
				},
				Kind: securityinsight.KindBasicDataConnectorKindAmazonWebServicesS3,
			}
			if _, err = client.DataConnectorsCreateOrUpdate(ctx, id, params); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DataConnectorAwsS3Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient
			id, err := dataconnectors.ParseDataConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			workspaceId := dataconnectors.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)

			existing, err := client.DataConnectorsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := existing.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model is nil", id)
			}

			dc, ok := (*model).(dataconnectors.AwsS3DataConnector)
			if !ok {
				return fmt.Errorf("%s was not an AWS S3 Data Connector", id)
			}

			outModel := DataConnectorAwsS3Model{
				Name:                    *dc.Name,
				LogAnalyticsWorkspaceId: workspaceId.ID(),
			}

			if prop := dc.Properties; prop != nil {
				outModel.AwsRoleArm = prop.RoleArn
				outModel.DestinationTable = prop.DestinationTable
				outModel.SqsUrls = prop.SqsUrls
			}

			return metadata.Encode(&model)
		},
	}
}

func (DataConnectorAwsS3Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := dataconnectors.ParseDataConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var plan DataConnectorAwsS3Model
			if err := metadata.Decode(&plan); err != nil {
				return err
			}

			client := metadata.Client.Sentinel.DataConnectorsClient

			resp, err := client.DataConnectorsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model is nil", id)
			}

			params, ok := (*model).(dataconnectors.AwsS3DataConnector)
			if !ok {
				return fmt.Errorf("%s was not an AWS S3 Data Connector", id)
			}

			if props := params.Properties; props != nil {
				if metadata.ResourceData.HasChange("aws_role_arn") {
					props.RoleArn = plan.AwsRoleArm
				}
				if metadata.ResourceData.HasChange("destination_table") {
					props.DestinationTable = plan.DestinationTable
				}
				if metadata.ResourceData.HasChange("sqs_urls") {
					props.SqsUrls = plan.SqsUrls
				}
			}

			if _, err := client.DataConnectorsCreateOrUpdate(ctx, *id, params); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r DataConnectorAwsS3Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient

			id, err := dataconnectors.ParseDataConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.DataConnectorsDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
