package sentinel

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-07-01-preview/watchlists"
	"time"

	commonValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	loganalyticsParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/parse"
	loganalyticsValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WatchlistResource struct{}

var _ sdk.Resource = WatchlistResource{}

type WatchlistModel struct {
	Name                    string   `tfschema:"name"`
	LogAnalyticsWorkspaceId string   `tfschema:"log_analytics_workspace_id"`
	DisplayName             string   `tfschema:"display_name"`
	Description             string   `tfschema:"description"`
	Labels                  []string `tfschema:"labels"`
	DefaultDuration         string   `tfschema:"default_duration"`
	ItemSearchKey           string   `tfschema:"item_search_key"`
}

func (r WatchlistResource) Arguments() map[string]*pluginsdk.Schema {
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
		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"item_search_key": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"labels": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
		"default_duration": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonValidate.ISO8601Duration,
		},
	}
}

func (r WatchlistResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r WatchlistResource) ResourceType() string {
	return "azurerm_sentinel_watchlist"
}

func (r WatchlistResource) ModelObject() interface{} {
	return &WatchlistModel{}
}

func (r WatchlistResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.WatchlistID
}

func (r WatchlistResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.WatchlistsClient

			var model WatchlistModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			workspaceId, err := loganalyticsParse.LogAnalyticsWorkspaceID(model.LogAnalyticsWorkspaceId)
			if err != nil {
				return fmt.Errorf("parsing Log Analytics Workspace ID: %w", err)
			}

			id := watchlists.NewWatchlistID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.WorkspaceName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := watchlists.Watchlist{
				Properties: &watchlists.WatchlistProperties{
					DisplayName: model.DisplayName,
					// The only supported provider for now is "Microsoft"
					Provider: "Microsoft",

					// The "source" represent the source file name which contains the watchlist items.
					// Setting them here is merely to make the API happy.
					Source: utils.String("a.csv"),

					ItemsSearchKey: model.ItemSearchKey,
				},
			}

			if model.Description != "" {
				param.Properties.Description = &model.Description
			}
			if len(model.Labels) != 0 {
				param.Properties.Labels = &model.Labels
			}
			if model.DefaultDuration != "" {
				param.Properties.DefaultDuration = &model.DefaultDuration
			}

			_, err = client.CreateOrUpdate(ctx, id, param)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r WatchlistResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.WatchlistsClient
			id, err := watchlists.ParseWatchlistID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model is nil", id)
			}

			watchlistModel := WatchlistModel{
				Name:                    *model.Name,
				LogAnalyticsWorkspaceId: loganalyticsParse.NewLogAnalyticsWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID(),
			}

			if props := model.Properties; props != nil {
				if props.DisplayName != "" {
					watchlistModel.DisplayName = props.DisplayName
				}
				if props.Description != nil {
					watchlistModel.Description = *props.Description
				}
				if props.Labels != nil {
					watchlistModel.Labels = *props.Labels
				}
				if props.DefaultDuration != nil {
					watchlistModel.DefaultDuration = *props.DefaultDuration
				}
				if props.ItemsSearchKey != "" {
					watchlistModel.ItemSearchKey = props.ItemsSearchKey
				}
			}

			return metadata.Encode(&watchlistModel)
		},
	}
}

func (r WatchlistResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.WatchlistsClient

			id, err := watchlists.ParseWatchlistID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
