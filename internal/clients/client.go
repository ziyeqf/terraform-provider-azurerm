// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/subscription/mgmt/2020-09-01/subscription"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/validation"
	aadb2c_v2021_04_01_preview "github.com/hashicorp/go-azure-sdk/resource-manager/aadb2c/2021-04-01-preview"
	analysisservices_v2017_08_01 "github.com/hashicorp/go-azure-sdk/resource-manager/analysisservices/2017-08-01"
	azurestackhci_v2023_03_01 "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01"
	storagecache_2023_05_01 "github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-05-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	resource "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/client"
)

type Client struct {
	autoClient

	// StopContext is used for propagating control from Terraform Core (e.g. Ctrl/Cmd+C)
	StopContext context.Context

	Account  *ResourceManagerAccount
	Features features.UserFeatures

	AadB2c                       *aadb2c_v2021_04_01_preview.Client
	AnalysisServices             *analysisservices_v2017_08_01.Client
	AzureManagedLustreFileSystem *storagecache_2023_05_01.Client
	AzureStackHCI                *azurestackhci_v2023_03_01.Client
	Resource                     *resource.Client
	Subscription                 *subscription.Client
}

// NOTE: it should be possible for this method to become Private once the top level Client's removed

func (client *Client) Build(ctx context.Context, o *common.ClientOptions) error {
	autorest.Count429AsRetry = false
	// Disable the Azure SDK for Go's validation since it's unhelpful for our use-case
	validation.Disabled = true

	if err := buildAutoClients(&client.autoClient, o); err != nil {
		return fmt.Errorf("building auto-clients: %+v", err)
	}

	client.Features = o.Features
	client.StopContext = ctx

	var err error
	if client.Resource, err = resource.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Resource: %+v", err)
	}

	return nil
}
