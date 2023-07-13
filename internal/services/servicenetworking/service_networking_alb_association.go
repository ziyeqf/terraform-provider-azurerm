package servicenetworking

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-05-01-preview/associationsinterface"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type TrafficControllerAssociationResource struct{}

type AssociationModel struct {
	Name                          string            `tfschema:"name"`
	ContainerApplicationGatewayId string            `tfschema:"container_application_gateway_id"`
	SubnetId                      string            `tfschema:"subnet_id"`
	Location                      string            `tfschema:"location"`
	Tags                          map[string]string `tfschema:"tags"`
}

var _ sdk.ResourceWithUpdate = TrafficControllerAssociationResource{}

func (t TrafficControllerAssociationResource) Arguments() map[string]*schema.Schema {
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
			ValidateFunc: associationsinterface.ValidateTrafficControllerID,
		},

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (t TrafficControllerAssociationResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (t TrafficControllerAssociationResource) ModelObject() interface{} {
	return &AssociationModel{}
}

func (t TrafficControllerAssociationResource) ResourceType() string {
	return "azurerm_alb_association"
}

func (t TrafficControllerAssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return associationsinterface.ValidateAssociationID
}
func (t TrafficControllerAssociationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan AssociationModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %v", err)
			}

			client := metadata.Client.ServiceNetworking.ServiceNetworkingClient.AssociationsInterface

			parsedTrafficControllerId, err := associationsinterface.ParseTrafficControllerID(plan.ContainerApplicationGatewayId)
			if err != nil {
				return fmt.Errorf("parsing traffic controller id %v", err)
			}

			id := associationsinterface.NewAssociationID(parsedTrafficControllerId.SubscriptionId, parsedTrafficControllerId.ResourceGroupName, parsedTrafficControllerId.TrafficControllerName, plan.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of exisiting %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(t.ResourceType(), id.ID())
			}

			association := associationsinterface.Association{
				Location: location.Normalize(plan.Location),
				Properties: &associationsinterface.AssociationProperties{
					Subnet: &associationsinterface.AssociationSubnet{
						Id: plan.SubnetId,
					},
					AssociationType: associationsinterface.AssociationTypeSubnets,
				},
			}

			if len(plan.Tags) > 0 {
				association.Tags = &plan.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, association); err != nil {
				return fmt.Errorf("creating %v", err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (t TrafficControllerAssociationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.ServiceNetworkingClient.AssociationsInterface

			id, err := associationsinterface.ParseAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing id %v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retreiving %s: %v", id.ID(), err)
			}

			trafficControllerId := associationsinterface.NewTrafficControllerID(id.SubscriptionId, id.ResourceGroupName, id.TrafficControllerName)
			state := AssociationModel{
				Name:                          id.AssociationName,
				ContainerApplicationGatewayId: trafficControllerId.ID(),
			}

			if model := resp.Model; model != nil {
				state.Tags = pointer.From(model.Tags)
				state.Location = location.Normalize(model.Location)

				if prop := model.Properties; prop != nil {
					if prop.Subnet != nil {
						state.SubnetId = prop.Subnet.Id
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (t TrafficControllerAssociationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan AssociationModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %v", err)
			}

			id, err := associationsinterface.ParseAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing id %v", err)
			}

			client := metadata.Client.ServiceNetworking.ServiceNetworkingClient.AssociationsInterface
			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %v", id.ID(), err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model is nil", id.ID())
			}
			association := *resp.Model

			if metadata.ResourceData.HasChange("tags") {
				association.Tags = &plan.Tags
			}

			if err = client.CreateOrUpdateThenPoll(ctx, *id, association); err != nil {
				return fmt.Errorf("updating %s: %v", id.ID(), err)
			}
			return nil
		},
	}
}

func (t TrafficControllerAssociationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.ServiceNetworkingClient.AssociationsInterface

			id, err := associationsinterface.ParseAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing id %s: %v", metadata.ResourceData.Id(), err)
			}

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", id.ID(), err)
			}

			return nil
		},
	}
}
