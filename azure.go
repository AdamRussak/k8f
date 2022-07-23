// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

var (
	subscriptionID string
)

func main() {
	subscriptionID = getEnvVarOrExit("AZURE_SUBSCRIPTION_ID")

	cred, err := azidentity.NewAzureCLICredential(nil)
	onErrorFail(err, "Authentication Failed")
	ctx := context.Background()
	resourceGroups, err := listResourceGroup(ctx, cred)
	onErrorFail(err, "ResourceGroup List Failed")
	for _, resource := range resourceGroups {
		log.Printf("Resource Group Name: %s, Location: %s", *resource.Name, *resource.Location)
	}
}

func listResourceGroup(ctx context.Context, cred azcore.TokenCredential) ([]*armresources.ResourceGroup, error) {
	resourceGroupClient, err := armresources.NewResourceGroupsClient(subscriptionID, cred, nil)
	onErrorFail(err, "Resource Client Authentication Failed")

	resultPager := resourceGroupClient.NewListPager(nil)

	resourceGroups := make([]*armresources.ResourceGroup, 0)
	for resultPager.More() {
		pageResp, err := resultPager.NextPage(ctx)
		onErrorFail(err, "Next Page Failed")
		resourceGroups = append(resourceGroups, pageResp.ResourceGroupListResult.Value...)
	}
	return resourceGroups, nil
}
