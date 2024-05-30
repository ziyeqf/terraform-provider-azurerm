// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package environments

func AzurePublicCanary() *Environment {
	// Canary is Azure Public with a different Microsoft Graph endpoint
	env := AzurePublic()
	env.Name = "Canary"
	env.ResourceManager = ResourceManagerAPI("https://eastus2euap.management.azure.com").WithResourceIdentifier("https://management.azure.com")
	env.MicrosoftGraph = MicrosoftGraphAPI("https://canary.graph.microsoft.com").WithResourceIdentifier("https://canary.graph.microsoft.com")
	return env
}
