package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-07-01-preview/alertrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-07-01-preview/alertruletemplates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-07-01-preview/automationrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-07-01-preview/dataconnectors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-07-01-preview/watchlistitems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-07-01-preview/watchlists"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AlertRulesClient         *alertrules.AlertRulesClient
	AlertRuleTemplatesClient *alertruletemplates.AlertRuleTemplatesClient
	AutomationRulesClient    *automationrules.AutomationRulesClient
	DataConnectorsClient     *dataconnectors.DataConnectorsClient
	WatchlistsClient         *watchlists.WatchlistsClient
	WatchlistItemsClient     *watchlistitems.WatchlistItemsClient
}

func NewClient(o *common.ClientOptions) *Client {
	alertRulesClient := alertrules.NewAlertRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&alertRulesClient.Client, o.ResourceManagerAuthorizer)

	alertRuleTemplatesClient := alertruletemplates.NewAlertRuleTemplatesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&alertRuleTemplatesClient.Client, o.ResourceManagerAuthorizer)

	automationRulesClient := automationrules.NewAutomationRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&automationRulesClient.Client, o.ResourceManagerAuthorizer)

	dataConnectorsClient := dataconnectors.NewDataConnectorsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&dataConnectorsClient.Client, o.ResourceManagerAuthorizer)

	watchListsClient := watchlists.NewWatchlistsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&watchListsClient.Client, o.ResourceManagerAuthorizer)

	watchListItemsClient := watchlistitems.NewWatchlistItemsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&watchListItemsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AlertRulesClient:         &alertRulesClient,
		AlertRuleTemplatesClient: &alertRuleTemplatesClient,
		AutomationRulesClient:    &automationRulesClient,
		DataConnectorsClient:     &dataConnectorsClient,
		WatchlistsClient:         &watchListsClient,
		WatchlistItemsClient:     &watchListItemsClient,
	}
}
