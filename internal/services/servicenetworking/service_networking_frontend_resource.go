package servicenetworking

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-05-01-preview/frontendsinterface"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type FrontendsResource struct{}

type FrontendsModel struct {
	Name                          string                 `tfschema:"name"`
	ContainerApplicationGatewayId string                 `tfschema:"container_application_gateway_id"`
	Location                      string                 `tfschema:"location"`
	Fqdn                          string                 `tfschema:"fully_qualified_domain_name"`
	Tags                          map[string]interface{} `tfschema:"tags"`
}

var _ sdk.Resource = FrontendsResource{}

func (f FrontendsResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"container_application_gateway_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: frontendsinterface.ValidateTrafficControllerID,
		},

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (f FrontendsResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"fully_qualified_domain_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (f FrontendsResource) ModelObject() interface{} {
	return &FrontendsModel{}
}

func (f FrontendsResource) ResourceType() string {
	return "azurerm_service_networking_frontend"
}

func (f FrontendsResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return frontendsinterface.ValidateFrontendID
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

			trafficControllerId, err := frontendsinterface.ParseTrafficControllerID(plan.ContainerApplicationGatewayId)
			if err != nil {
				return fmt.Errorf("parsing traffic controller id %v", err)
			}

			id := frontendsinterface.NewFrontendID(trafficControllerId.SubscriptionId, trafficControllerId.ResourceGroupName, trafficControllerId.TrafficControllerName, plan.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("checking presence of existing %s: %v", id.ID(), err)
				}
			}

			if !response.WasNotFound(resp.HttpResponse) {
				return tf.ImportAsExistsError(f.ResourceType(), id.ID())
			}

			frontend := frontendsinterface.Frontend{
				Location:   location.Normalize(plan.Location),
				Properties: &frontendsinterface.FrontendProperties{},
				Tags:       tags.Expand(plan.Tags),
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
				Name:                          id.FrontendName,
				ContainerApplicationGatewayId: trafficControllerId.ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = model.Location
				state.Tags = tags.Flatten(model.Tags)

				if prop := model.Properties; prop != nil {
					state.Fqdn = pointer.From(prop.Fqdn)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (f FrontendsResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.ServiceNetworkingClient.FrontendsInterface
			id, err := frontendsinterface.ParseFrontendID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing ID %q: %+v", metadata.ResourceData.Id(), err)
			}

			var plan FrontendsModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retiring `azurerm_service_networking_frontend` %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retiring `azurerm_service_networking_frontend` %s: Model was nil", *id)
			}

			model := *resp.Model

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = tags.Expand(plan.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, model); err != nil {
				return fmt.Errorf("updating `azurerm_service_networking_frontend` %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (f FrontendsResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.ServiceNetworkingClient.FrontendsInterface

			id, err := frontendsinterface.ParseFrontendID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %q: %+v", metadata.ResourceData.Id(), err)
			}

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %q: %+v", id.ID(), err)
			}

			return nil
		},
	}
}
