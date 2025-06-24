// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AppServicePlansClient   *web.AppServicePlansClient
	AppServicesClient       *web.AppsClient
	CertificatesClient      *web.CertificatesClient
	CertificatesOrderClient *web.AppServiceCertificateOrdersClient
}

func NewClient(o *common.ClientOptions) *Client {
	appServicePlansClient := web.NewAppServicePlansClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&appServicePlansClient.Client, o.ResourceManagerAuthorizer)

	appServicesClient := web.NewAppsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&appServicesClient.Client, o.ResourceManagerAuthorizer)

	certificatesClient := web.NewCertificatesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&certificatesClient.Client, o.ResourceManagerAuthorizer)

	certificatesOrderClient := web.NewAppServiceCertificateOrdersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&certificatesOrderClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AppServicePlansClient:   &appServicePlansClient,
		AppServicesClient:       &appServicesClient,
		CertificatesClient:      &certificatesClient,
		CertificatesOrderClient: &certificatesOrderClient,
	}
}
