package kusto

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2022-02-01/clusters"
	newCluster "github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2022-07-07/clusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKustoCluster() *pluginsdk.Resource {
	s := &pluginsdk.Resource{
		Create: resourceKustoClusterCreateUpdate,
		Read:   resourceKustoClusterRead,
		Update: resourceKustoClusterCreateUpdate,
		Delete: resourceKustoClusterDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KustoAttachedClusterV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := clusters.ParseClusterID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ClusterName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"sku": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							// using API version later than 2022-07-07 to unblock user access to latest skus.
							// next PR will to API version upgrade to version 2022-12-10 and align the api version in kusto.
							ValidateFunc: validation.StringInSlice(newCluster.PossibleValuesForAzureSkuName(),
								false),
						},

						"capacity": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(1, 1000),
						},
					},
				},
			},

			"allowed_fqdns": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"allowed_ip_ranges": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"trusted_external_tenants": {
				Type:       pluginsdk.TypeList,
				Optional:   true,
				Computed:   true,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.Any(validation.IsUUID, validation.StringIsEmpty, validation.StringInSlice([]string{"*"}, false)),
				},
			},

			"optimized_auto_scale": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"minimum_instances": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
						"maximum_instances": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
					},
				},
			},

			"virtual_network_configuration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"engine_public_ip_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"data_management_public_ip_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"language_extensions": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice(clusters.PossibleValuesForLanguageExtensionName(), false),
				},
			},

			"engine": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(clusters.PossibleValuesForEngineType(), false),
			},

			"uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"data_ingestion_uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"public_ip_type": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      string(clusters.PublicIPTypeIPvFour),
				ValidateFunc: validation.StringInSlice(clusters.PossibleValuesForPublicIPType(), false),
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"outbound_network_access_restricted": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"double_encryption_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"auto_stop_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"disk_encryption_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"streaming_ingestion_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"purge_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"zones": commonschema.ZonesMultipleOptionalForceNew(),

			"tags": commonschema.Tags(),
		},
	}

	if features.FourPointOhBeta() {
		s.Schema["engine"].Default = string(clusters.EngineTypeVThree)
	} else {
		s.Schema["engine"].Default = string(clusters.EngineTypeVTwo)
	}

	return s
}

func resourceKustoClusterCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Kusto Cluster creation.")

	id := clusters.NewClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil && !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_kusto_cluster", id.ID())
		}
	}

	locks.ByName(id.ClusterName, "azurerm_kusto_cluster")
	defer locks.UnlockByName(id.ClusterName, "azurerm_kusto_cluster")

	sku, err := expandKustoClusterSku(d.Get("sku").([]interface{}))
	if err != nil {
		return err
	}

	optimizedAutoScale := expandOptimizedAutoScale(d.Get("optimized_auto_scale").([]interface{}))

	if optimizedAutoScale != nil && optimizedAutoScale.IsEnabled {
		if sku.Capacity == nil {
			return fmt.Errorf("sku.capacity could not be empty")
		}
		// Ensure that requested Capcity is always between min and max to support updating to not overlapping autoscale ranges
		if *sku.Capacity < optimizedAutoScale.Minimum {
			sku.Capacity = utils.Int64(optimizedAutoScale.Minimum)
		}
		if *sku.Capacity > optimizedAutoScale.Maximum {
			sku.Capacity = utils.Int64(optimizedAutoScale.Maximum)
		}

		// Capacity must be set for the initial creation when using OptimizedAutoScaling but cannot be updated
		if d.HasChange("sku.0.capacity") && !d.IsNewResource() {
			return fmt.Errorf("cannot change `sku.capacity` when `optimized_auto_scaling.enabled` is set to `true`")
		}

		if optimizedAutoScale.Minimum > optimizedAutoScale.Maximum {
			return fmt.Errorf("`optimized_auto_scaling.maximum_instances` must be >= `optimized_auto_scaling.minimum_instances`")
		}
	}

	engine := clusters.EngineType(d.Get("engine").(string))

	publicNetworkAccess := clusters.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = clusters.PublicNetworkAccessDisabled
	}

	publicIPType := clusters.PublicIPType(d.Get("public_ip_type").(string))

	clusterProperties := clusters.ClusterProperties{
		OptimizedAutoscale:     optimizedAutoScale,
		EnableAutoStop:         utils.Bool(d.Get("auto_stop_enabled").(bool)),
		EnableDiskEncryption:   utils.Bool(d.Get("disk_encryption_enabled").(bool)),
		EnableDoubleEncryption: utils.Bool(d.Get("double_encryption_enabled").(bool)),
		EnableStreamingIngest:  utils.Bool(d.Get("streaming_ingestion_enabled").(bool)),
		EnablePurge:            utils.Bool(d.Get("purge_enabled").(bool)),
		EngineType:             &engine,
		PublicNetworkAccess:    &publicNetworkAccess,
		PublicIPType:           &publicIPType,
		TrustedExternalTenants: expandTrustedExternalTenants(d.Get("trusted_external_tenants").([]interface{})),
	}

	if v, ok := d.GetOk("virtual_network_configuration"); ok {
		vnet := expandKustoClusterVNET(v.([]interface{}))
		clusterProperties.VirtualNetworkConfiguration = vnet
	}

	if v, ok := d.GetOk("allowed_fqdns"); ok {
		clusterProperties.AllowedFqdnList, _ = expandKustoListString(v.([]interface{}))
	}

	if v, ok := d.GetOk("allowed_ip_ranges"); ok {
		clusterProperties.AllowedIPRangeList, _ = expandKustoListString(v.([]interface{}))
	}

	restrictOutboundNetworkAccess := clusters.ClusterNetworkAccessFlagDisabled
	if v, ok := d.GetOk("outbound_network_access_restricted"); ok {
		if v.(bool) {
			restrictOutboundNetworkAccess = clusters.ClusterNetworkAccessFlagEnabled
		}
	}
	clusterProperties.RestrictOutboundNetworkAccess = &restrictOutboundNetworkAccess

	expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	kustoCluster := clusters.Cluster{
		Name:       utils.String(id.ClusterName),
		Location:   location.Normalize(d.Get("location").(string)),
		Identity:   expandedIdentity,
		Sku:        *sku,
		Properties: &clusterProperties,
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		kustoCluster.Zones = &zones
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, kustoCluster, clusters.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if v, ok := d.GetOk("language_extensions"); ok {
		languageExtensions := expandKustoClusterLanguageExtensions(v.([]interface{}))

		currentLanguageExtensions, err := client.ListLanguageExtensions(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving the language extensions on %s: %+v", id, err)
		}

		languageExtensionsToAdd := diffLanguageExtensions(*languageExtensions.Value, *currentLanguageExtensions.Model.Value)
		if len(languageExtensionsToAdd) > 0 {
			languageExtensionsListToAdd := clusters.LanguageExtensionsList{
				Value: &languageExtensionsToAdd,
			}

			resp, err := client.AddLanguageExtensions(ctx, id, languageExtensionsListToAdd)
			if err != nil {
				return fmt.Errorf("adding language extensions to %s: %+v", id, err)
			}
			if err = resp.Poller.PollUntilDone(); err != nil {
				return fmt.Errorf("waiting for the addition of language extensions on %s: %+v", id, err)
			}
		}

		languageExtensionsToRemove := diffLanguageExtensions(*currentLanguageExtensions.Model.Value, *languageExtensions.Value)
		if len(languageExtensionsToRemove) > 0 {
			languageExtensionsListToRemove := clusters.LanguageExtensionsList{
				Value: &languageExtensionsToRemove,
			}

			removeLanguageExtensionsFuture, err := client.RemoveLanguageExtensions(ctx, id, languageExtensionsListToRemove)
			if err != nil {
				return fmt.Errorf("removing language extensions from %s: %+v", id, err)
			}
			if err = removeLanguageExtensionsFuture.Poller.PollUntilDone(); err != nil {
				return fmt.Errorf("waiting for the removal of language extensions from %s: %+v", id, err)
			}
		}
	}

	return resourceKustoClusterRead(d, meta)
}

func resourceKustoClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(&model.Location))
		d.Set("zones", zones.FlattenUntyped(model.Zones))

		identity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %s", err)
		}

		if err := d.Set("sku", flattenKustoClusterSku(&model.Sku)); err != nil {
			return fmt.Errorf("setting `sku`: %+v", err)
		}

		if props := model.Properties; props != nil {
			if props.PublicNetworkAccess != nil {
				d.Set("public_network_access_enabled", *props.PublicNetworkAccess == clusters.PublicNetworkAccessEnabled)
			}

			if props.RestrictOutboundNetworkAccess != nil {
				d.Set("outbound_network_access_restricted", *props.RestrictOutboundNetworkAccess == clusters.ClusterNetworkAccessFlagEnabled)
			}

			if err := d.Set("optimized_auto_scale", flattenOptimizedAutoScale(props.OptimizedAutoscale)); err != nil {
				return fmt.Errorf("setting `optimized_auto_scale`: %+v", err)
			}
			d.Set("allowed_fqdns", props.AllowedFqdnList)
			d.Set("allowed_ip_ranges", props.AllowedIPRangeList)
			d.Set("double_encryption_enabled", props.EnableDoubleEncryption)
			d.Set("trusted_external_tenants", flattenTrustedExternalTenants(props.TrustedExternalTenants))
			d.Set("auto_stop_enabled", props.EnableAutoStop)
			d.Set("disk_encryption_enabled", props.EnableDiskEncryption)
			d.Set("streaming_ingestion_enabled", props.EnableStreamingIngest)
			d.Set("purge_enabled", props.EnablePurge)
			d.Set("virtual_network_configuration", flattenKustoClusterVNET(props.VirtualNetworkConfiguration))
			d.Set("language_extensions", flattenKustoClusterLanguageExtensions(props.LanguageExtensions))
			d.Set("uri", props.Uri)
			d.Set("data_ingestion_uri", props.DataIngestionUri)
			d.Set("engine", props.EngineType)
			d.Set("public_ip_type", props.PublicIPType)

		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceKustoClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	return nil
}

func expandOptimizedAutoScale(input []interface{}) *clusters.OptimizedAutoscale {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	config := input[0].(map[string]interface{})
	optimizedAutoScale := &clusters.OptimizedAutoscale{
		Version:   1,
		IsEnabled: true,
		Minimum:   int64(config["minimum_instances"].(int)),
		Maximum:   int64(config["maximum_instances"].(int)),
	}

	return optimizedAutoScale
}

func flattenOptimizedAutoScale(optimizedAutoScale *clusters.OptimizedAutoscale) []interface{} {
	if optimizedAutoScale == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"maximum_instances": int(optimizedAutoScale.Maximum),
			"minimum_instances": int(optimizedAutoScale.Minimum),
		},
	}
}

func expandKustoListString(input []interface{}) (*[]string, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("list of string is empty")
	}

	result := make([]string, 0)

	for _, v := range input {
		result = append(result, v.(string))
	}

	return &result, nil
}

func expandKustoClusterSku(input []interface{}) (*clusters.AzureSku, error) {
	sku := input[0].(map[string]interface{})
	name := sku["name"].(string)

	skuNamePrefixToTier := map[string]string{
		"Dev(No SLA)": "Basic",
		"Standard":    "Standard",
	}

	skuNamePrefix := strings.Split(sku["name"].(string), "_")[0]
	tier, ok := skuNamePrefixToTier[skuNamePrefix]
	if !ok {
		return nil, fmt.Errorf("sku name begins with invalid tier, possible are Dev(No SLA) and Standard but is: %q", skuNamePrefix)
	}
	capacity := sku["capacity"].(int)

	azureSku := clusters.AzureSku{
		Name:     clusters.AzureSkuName(name),
		Tier:     clusters.AzureSkuTier(tier),
		Capacity: utils.Int64(int64(capacity)),
	}

	return &azureSku, nil
}

func expandKustoClusterVNET(input []interface{}) *clusters.VirtualNetworkConfiguration {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	vnet := input[0].(map[string]interface{})
	subnetID := vnet["subnet_id"].(string)
	enginePublicIPID := vnet["engine_public_ip_id"].(string)
	dataManagementPublicIPID := vnet["data_management_public_ip_id"].(string)

	return &clusters.VirtualNetworkConfiguration{
		SubnetId:                 subnetID,
		EnginePublicIPId:         enginePublicIPID,
		DataManagementPublicIPId: dataManagementPublicIPID,
	}
}

func expandKustoClusterLanguageExtensions(input []interface{}) *clusters.LanguageExtensionsList {
	if len(input) == 0 {
		return nil
	}

	extensions := make([]clusters.LanguageExtension, 0)
	for _, language := range input {
		name := clusters.LanguageExtensionName(language.(string))
		v := clusters.LanguageExtension{
			LanguageExtensionName: &name,
		}
		extensions = append(extensions, v)
	}

	return &clusters.LanguageExtensionsList{
		Value: &extensions,
	}
}

func flattenKustoClusterSku(sku *clusters.AzureSku) []interface{} {
	if sku == nil {
		return []interface{}{}
	}

	s := map[string]interface{}{
		"name": string(sku.Name),
	}

	if sku.Capacity != nil {
		s["capacity"] = int(*sku.Capacity)
	}

	return []interface{}{s}
}

func flattenKustoClusterVNET(vnet *clusters.VirtualNetworkConfiguration) []interface{} {
	if vnet == nil {
		return []interface{}{}
	}

	output := map[string]interface{}{
		"subnet_id":                    vnet.SubnetId,
		"engine_public_ip_id":          vnet.EnginePublicIPId,
		"data_management_public_ip_id": vnet.DataManagementPublicIPId,
	}

	return []interface{}{output}
}

func flattenKustoClusterLanguageExtensions(extensions *clusters.LanguageExtensionsList) []interface{} {
	if extensions == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)
	if extensions.Value != nil {
		for _, v := range *extensions.Value {
			output = append(output, v.LanguageExtensionName)
		}
	}

	return output
}

func diffLanguageExtensions(a, b []clusters.LanguageExtension) []clusters.LanguageExtension {
	target := make(map[string]bool)
	for _, x := range b {
		if x.LanguageExtensionName != nil {
			target[string(*x.LanguageExtensionName)] = true
		}
	}

	diff := make([]clusters.LanguageExtension, 0)
	for _, x := range a {
		if x.LanguageExtensionName != nil {
			if _, ok := target[string(*x.LanguageExtensionName)]; !ok {
				diff = append(diff, x)
			}
		}
	}

	return diff
}
