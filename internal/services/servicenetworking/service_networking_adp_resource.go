package servicenetworking

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-05-01-preview/trafficcontrollerinterface"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type TrafficControllerResource struct{}

type TrafficControllerModel struct {
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Name              string            `tfschema:"name"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`
}

var _ sdk.ResourceWithUpdate = TrafficControllerResource{}

func (t TrafficControllerResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"tags": tags.Schema(),
	}
}

func (t TrafficControllerResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (t TrafficControllerResource) ModelObject() interface{} {
	return &TrafficControllerModel{}
}

func (t TrafficControllerResource) ResourceType() string {
	return "azurerm_service_networking_adp"
}

func (t TrafficControllerResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return trafficcontrollerinterface.ValidateTrafficControllerID
}

func (t TrafficControllerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan TrafficControllerModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %v", err)
			}

			client := metadata.Client.ServiceNetworking.ServiceNetworkingClient
			SubscriptionId := metadata.Client.Account.SubscriptionId

			id := trafficcontrollerinterface.NewTrafficControllerID(SubscriptionId, plan.ResourceGroupName, plan.Name)

			existing, err := client.TrafficControllerInterface.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing Traffic Controller %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(t.ResourceType(), id)
			}

			controller := trafficcontrollerinterface.TrafficController{
				Location: location.Normalize(plan.Location),
				Tags:     pointer.To(plan.Tags),
			}

			err = client.TrafficControllerInterface.CreateOrUpdateThenPoll(ctx, id, controller)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (t TrafficControllerResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.ServiceNetworkingClient

			id, err := trafficcontrollerinterface.ParseTrafficControllerID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			resp, err := client.TrafficControllerInterface.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", metadata.ResourceData.Id(), err)
			}

			state := TrafficControllerModel{
				Name:              id.TrafficControllerName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = model.Location
				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (t TrafficControllerResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan TrafficControllerModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %v", err)
			}
			client := metadata.Client.ServiceNetworking.ServiceNetworkingClient

			id, err := trafficcontrollerinterface.ParseTrafficControllerID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			existing, err := client.TrafficControllerInterface.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retreiving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("existing Traffic Controller %s has no model", id)
			}

			controller := existing.Model

			if metadata.ResourceData.HasChange("tags") {
				controller.Tags = pointer.To(plan.Tags)
			}

			err = client.TrafficControllerInterface.CreateOrUpdateThenPoll(ctx, *id, *controller)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (t TrafficControllerResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := trafficcontrollerinterface.ParseTrafficControllerID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			client := metadata.Client.ServiceNetworking.ServiceNetworkingClient
			err = client.TrafficControllerInterface.DeleteThenPoll(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", metadata.ResourceData.Id(), err)
			}

			return nil
		},
	}
}
