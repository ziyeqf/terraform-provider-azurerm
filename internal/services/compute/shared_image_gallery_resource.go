package compute

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-01-03/galleries"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSharedImageGallery() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSharedImageGalleryCreateUpdate,
		Read:   resourceSharedImageGalleryRead,
		Update: resourceSharedImageGalleryCreateUpdate,
		Delete: resourceSharedImageGalleryDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SharedImageGalleryID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SharedImageGalleryName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"tags": commonschema.Tags(),

			"unique_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"permissions": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(galleries.GallerySharingPermissionTypesCommunity),
					string(galleries.GallerySharingPermissionTypesPrivate),
					string(galleries.GallerySharingPermissionTypesGroups),
				}, false),
			},

			"publisher": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"email": {
							Optional: true,
							Type:     pluginsdk.TypeString,
						},
						"uri": {
							Optional: true,
							Type:     pluginsdk.TypeString,
						},
					},
				},
			},

			"public_name_prefix": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"eula": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceSharedImageGalleryCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleriesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Image Gallery creation.")

	id := galleries.NewGalleriesID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	description := d.Get("description").(string)
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id, galleries.DefaultGetOperationOptions())
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_shared_image_gallery", id.ID())
		}
	}

	gallery := galleries.Gallery{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &galleries.GalleryProperties{
			Description: utils.String(description),
		},
		Tags: tags.Expand(t),
	}

	prop := *gallery.Properties

	if p,ok := d.GetOk("permissions"); ok {
		permission := galleries.GallerySharingPermissionTypes(p.(string)) 
		prop.SharingProfile.Permissions = &permission
	}

	if p,ok := d.GetOk("publisher");ok {
		publisher := p.(map[string]interface{})
		email := publisher["email"].(string)
		uri := publisher["uri"].(string)
		prop.SharingProfile.CommunityGalleryInfo.PublisherContact = utils.String(email)
		prop.SharingProfile.CommunityGalleryInfo.PublisherUri = utils.String(uri)
	}

	err := client.CreateOrUpdateThenPoll(ctx, id, gallery)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSharedImageGalleryRead(d, meta)
}

func resourceSharedImageGalleryRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleriesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := galleries.ParseGalleriesID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, galleries.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Shared Image Gallery %q (Resource Group %q) was not found - removing from state", id.GalleryName, id.ResourceGroupName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on Shared Image Gallery %q (Resource Group %q): %+v", id.GalleryName, id.ResourceGroupName, err)
	}

	d.Set("name", id.GalleryName)
	d.Set("resource_group_name", id.ResourceGroupName)

	model := resp.Model

	if model == nil {
		return fmt.Errorf("reading Shared Image Gallery %q (Resource Group %q): empty response", id.GalleryName, id.ResourceGroupName)
	}

	if location := model.Location; location != "" {
		d.Set("location", azure.NormalizeLocation(location))
	}

	if props := model.Properties; props != nil {
		d.Set("description", props.Description)
		if identifier := props.Identifier; identifier != nil {
			d.Set("unique_name", identifier.UniqueName)
		}
	}

	return tags.FlattenAndSet(d, model.Tags)
}

func resourceSharedImageGalleryDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleriesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := galleries.ParseGalleriesID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting Shared Image Gallery %q (Resource Group %q): %+v", id.GalleryName, id.ResourceGroupName, err)
	}

	return nil
}
