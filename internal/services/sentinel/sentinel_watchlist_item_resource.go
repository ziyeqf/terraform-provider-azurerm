package sentinel

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-07-01-preview/watchlistitems"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type WatchlistItemResource struct{}

var _ sdk.ResourceWithUpdate = WatchlistItemResource{}

type WatchlistItemModel struct {
	Name        string                 `tfschema:"name"`
	WatchlistID string                 `tfschema:"watchlist_id"`
	Properties  map[string]interface{} `tfschema:"properties"`
}

func (r WatchlistItemResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},
		"watchlist_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WatchlistID,
		},
		"properties": {
			Type:     pluginsdk.TypeMap,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r WatchlistItemResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r WatchlistItemResource) ResourceType() string {
	return "azurerm_sentinel_watchlist_item"
}

func (r WatchlistItemResource) ModelObject() interface{} {
	return &WatchlistItemModel{}
}

func (r WatchlistItemResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.WatchlistItemID
}

func (r WatchlistItemResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.WatchlistItemsClient

			var model WatchlistItemModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			// Generate a random UUID as the resource name if the user doesn't specify it.
			if model.Name == "" {
				model.Name = uuid.New().String()
			}

			watchlistId, err := watchlistitems.ParseWatchlistID(model.WatchlistID)
			if err != nil {
				return err
			}

			id := watchlistitems.NewWatchlistItemID(watchlistId.SubscriptionId, watchlistId.ResourceGroupName, watchlistId.WorkspaceName, watchlistId.WatchlistAlias, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			params := watchlistitems.WatchlistItem{
				Properties: &watchlistitems.WatchlistItemProperties{
					ItemsKeyValue: model.Properties,
				},
			}

			if _, err = client.CreateOrUpdate(ctx, id, params); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r WatchlistItemResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.WatchlistItemsClient
			id, err := watchlistitems.ParseWatchlistItemID(metadata.ResourceData.Id())
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

			watchlistId := watchlistitems.NewWatchlistID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.WatchlistAlias)

			var name string
			var properties map[string]interface{}
			if model := resp.Model; model != nil {
				if model.Name != nil {
					name = *model.Name
				}
				if props := model.Properties; props != nil {
					if itemsKV := props.ItemsKeyValue; itemsKV != nil {
						properties = itemsKV.(map[string]interface{})
					}
				}
			}

			model := WatchlistItemModel{
				WatchlistID: watchlistId.ID(),
				Name:        name,
				Properties:  properties,
			}

			return metadata.Encode(&model)
		},
	}
}

func (r WatchlistItemResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.WatchlistItemsClient

			id, err := watchlistitems.ParseWatchlistItemID(metadata.ResourceData.Id())
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

func (r WatchlistItemResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.WatchlistItemsClient

			var model WatchlistItemModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			watchlistId, err := watchlistitems.ParseWatchlistID(model.WatchlistID)
			if err != nil {
				return err
			}
			id := watchlistitems.NewWatchlistItemID(watchlistId.SubscriptionId, watchlistId.ResourceGroupName, watchlistId.WorkspaceName, watchlistId.WatchlistAlias, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			update := watchlistitems.WatchlistItem{}

			if t := existing.Model; t != nil {
				update.Properties = t.Properties
			}

			if metadata.ResourceData.HasChange("properties") {
				if update.Properties == nil {
					update.Properties = &watchlistitems.WatchlistItemProperties{}
				}
				update.Properties.ItemsKeyValue = model.Properties
			}

			if _, err = client.CreateOrUpdate(ctx, id, update); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
