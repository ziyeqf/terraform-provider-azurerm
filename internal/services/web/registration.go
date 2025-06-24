// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Web"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"App Service (Web Apps)",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	resources := map[string]*pluginsdk.Resource{
		"azurerm_app_service_certificate":                           resourceAppServiceCertificate(),
		"azurerm_app_service_certificate_order":                     resourceAppServiceCertificateOrder(),
		"azurerm_app_service_custom_hostname_binding":               resourceAppServiceCustomHostnameBinding(),
		"azurerm_app_service_certificate_binding":                   resourceAppServiceCertificateBinding(),
		"azurerm_app_service_managed_certificate":                   resourceAppServiceManagedCertificate(),
		"azurerm_app_service_public_certificate":                    resourceAppServicePublicCertificate(),
		"azurerm_app_service_slot_custom_hostname_binding":          resourceAppServiceSlotCustomHostnameBinding(),
		"azurerm_app_service_slot_virtual_network_swift_connection": resourceAppServiceSlotVirtualNetworkSwiftConnection(),
		"azurerm_app_service_virtual_network_swift_connection":      resourceAppServiceVirtualNetworkSwiftConnection(),
	}

	return resources
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{}
}
