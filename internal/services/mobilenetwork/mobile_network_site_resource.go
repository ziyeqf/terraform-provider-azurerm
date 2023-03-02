package mobilenetwork

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/mobilenetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/site"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SiteModel struct {
	Name             string            `tfschema:"name"`
	MobileNetworkId  string            `tfschema:"mobile_network_id"`
	Location         string            `tfschema:"location"`
	NetworkFunctions []string          `tfschema:"network_function_ids"`
	Tags             map[string]string `tfschema:"tags"`
}

type SiteResource struct{}

var _ sdk.ResourceWithUpdate = SiteResource{}

func (r SiteResource) ResourceType() string {
	return "azurerm_mobile_network_site"
}

func (r SiteResource) ModelObject() interface{} {
	return &SiteModel{}
}

func (r SiteResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return site.ValidateSiteID
}

func (r SiteResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"mobile_network_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: mobilenetwork.ValidateMobileNetworkID,
		},

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (r SiteResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"network_function_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r SiteResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SiteModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.SiteClient
			mobileNetworkId, err := mobilenetwork.ParseMobileNetworkID(model.MobileNetworkId)
			if err != nil {
				return err
			}

			id := site.NewSiteID(mobileNetworkId.SubscriptionId, mobileNetworkId.ResourceGroupName, mobileNetworkId.MobileNetworkName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &site.Site{
				Location:   location.Normalize(model.Location),
				Properties: &site.SitePropertiesFormat{},
				Tags:       &model.Tags,
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r SiteResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.SiteClient

			id, err := site.ParseSiteID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SiteModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			properties.SystemData = nil

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r SiteResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.SiteClient

			id, err := site.ParseSiteID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := SiteModel{
				Name:            id.SiteName,
				MobileNetworkId: mobilenetwork.NewMobileNetworkID(id.SubscriptionId, id.ResourceGroupName, id.MobileNetworkName).ID(),
				Location:        location.Normalize(model.Location),
			}

			if properties := model.Properties; properties != nil {
				state.NetworkFunctions = flattenSubResourceModel(properties.NetworkFunctions)
			}
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r SiteResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.SiteClient

			id, err := site.ParseSiteID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err := resourceMobileNetworkChildWaitForDeletion(ctx, id.ID(), func() (*http.Response, error) {
				resp, err := client.Get(ctx, *id)
				return resp.HttpResponse, err
			}); err != nil {
				return err
			}

			return nil
		},
	}
}

func flattenSubResourceModel(inputList *[]site.SubResource) []string {
	var outputList []string
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := input.Id
		outputList = append(outputList, output)
	}

	return outputList
}
