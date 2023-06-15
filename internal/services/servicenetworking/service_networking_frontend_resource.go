package servicenetworking

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-05-01-preview/frontendsinterface"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type FrontendsResource struct{}

type FrontendsModel struct {
	Name                string
	TrafficControllerId string
	Location            string
}

var _ sdk.ResourceWithUpdate = FrontendsResource{}

func (f FrontendsResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (f FrontendsResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (f FrontendsResource) ModelObject() interface{} {
	return &FrontendsModel{}
}

func (f FrontendsResource) ResourceType() string {
	return "azurerm_service_networking_frontend"
}

func (f FrontendsResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan FrontendsModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %v", err)
			}

			client := metadata.Client.ServiceNetworking.ServiceNetworkingClient.FrontendsInterface

			trafficControllerId, err := frontendsinterface.ParseTrafficControllerID(plan.TrafficControllerId)
			if err != nil {
				return fmt.Errorf("parsing traffic controller id %v", err)
			}

			id := frontendsinterface.NewFrontendID(trafficControllerId.SubscriptionId, trafficControllerId.ResourceGroupName, trafficControllerId.TrafficControllerName, plan.Name)

			frontend := frontendsinterface.Frontend{
				Location:   location.Normalize(plan.Location),
				Properties: &frontendsinterface.FrontendProperties{},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, frontend); err != nil {
				return fmt.Errorf("creating frontend %v", err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (f FrontendsResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.ServiceNetworkingClient.FrontendsInterface

			id, err := frontendsinterface.ParseFrontendID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing id %v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %v", err)
			}

			trafficControllerId := frontendsinterface.NewTrafficControllerID(id.SubscriptionId, id.ResourceGroupName, id.TrafficControllerName)
			state := FrontendsModel{
				Name:                id.FrontendName,
				TrafficControllerId: trafficControllerId.ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = model.Location
			}

			return metadata.Encode(&state)
		},
	}
}

func (f FrontendsResource) Delete() sdk.ResourceFunc {
	//TODO implement me
	panic("implement me")
}

func (f FrontendsResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	//TODO implement me
	panic("implement me")
}

func (f FrontendsResource) Update() sdk.ResourceFunc {
	//TODO implement me
	panic("implement me")
}
