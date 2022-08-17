package sentinel

import (
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-07-01-preview/dataconnectors"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	loganalyticsParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/parse"
	loganalyticsValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSentinelDataConnectorAzureAdvancedThreatProtection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSentinelDataConnectorAzureAdvancedThreatProtectionCreate,
		Read:   resourceSentinelDataConnectorAzureAdvancedThreatProtectionRead,
		Delete: resourceSentinelDataConnectorAzureAdvancedThreatProtectionDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.DataConnectorID(id)
			return err
		}, importSentinelDataConnector(dataconnectors.DataConnectorKindAzureAdvancedThreatProtection)),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
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

			"tenant_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func resourceSentinelDataConnectorAzureAdvancedThreatProtectionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.DataConnectorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := dataconnectors.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	id := dataconnectors.NewDataConnectorID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, name)

	if d.IsNewResource() {
		resp, err := client.DataConnectorsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_sentinel_data_connector_azure_advanced_threat_protection", id.ID())
		}
	}

	tenantId := d.Get("tenant_id").(string)
	if tenantId == "" {
		tenantId = meta.(*clients.Client).Account.TenantId
	}

	param := dataconnectors.AATPDataConnector{
		Name: &name,
		Properties: &dataconnectors.AATPDataConnectorProperties{
			DataTypes: &dataconnectors.AlertsDataTypeOfDataConnector{
				Alerts: dataconnectors.DataConnectorDataTypeCommon{
					State: dataconnectors.DataTypeStateEnabled,
				},
			},
		},
	}

	if _, err = client.DataConnectorsCreateOrUpdate(ctx, id, param); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelDataConnectorAzureAdvancedThreatProtectionRead(d, meta)
}

func resourceSentinelDataConnectorAzureAdvancedThreatProtectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.DataConnectorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataconnectors.ParseDataConnectorID(d.Id())
	if err != nil {
		return err
	}
	workspaceId := loganalyticsParse.NewLogAnalyticsWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)

	resp, err := client.DataConnectorsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("retrieving %s: model is nil", id)
	}

	dc, ok := (*model).(dataconnectors.AATPDataConnector)
	if !ok {
		return fmt.Errorf("%s was not an Azure Advanced Threat Protection Data Connector", id)
	}

	d.Set("name", dc.Name)
	d.Set("log_analytics_workspace_id", workspaceId.ID())
	if dc.Properties != nil {
		d.Set("tenant_id", dc.Properties.TenantId)
	}

	return nil
}

func resourceSentinelDataConnectorAzureAdvancedThreatProtectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.DataConnectorsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataconnectors.ParseDataConnectorID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.DataConnectorsDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
