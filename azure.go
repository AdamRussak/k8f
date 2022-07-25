// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

var (
	subscriptionID  string
	groupsClient    resources.GroupsClient
	resourcesClient resources.Client
)

func main() {
	subscriptionID = getEnvVarOrExit("AZURE_SUBSCRIPTION_ID")
	cred, err := azidentity.NewAzureCLICredential(nil)
	onErrorFail(err, "Authentication Failed")
	ctx := context.Background()
	resourceGroups, err := listResourceGroup(ctx, cred)
	onErrorFail(err, "ResourceGroup List Failed")
	var rgList []ResourceGroup
	for _, resource := range resourceGroups {
		log.Printf("Resource Group Name: %s, Location: %s", *resource.Name, *resource.Location)
		rgList = append(rgList, ResourceGroup{resource.Location, resource.ManagedBy, resource.Tags, resource.ID, resource.Name, resource.Type})
		listResource(*resource.Name, cred)
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

//works!! need to start orgenizing the input and output
func listResource(rg string, cred azcore.TokenCredential) {
	resourceClient, err := armresources.NewClient(subscriptionID, cred, nil)
	onErrorFail(err, "Failed to create Resource Client")
	// s := "$filter=eq(Microsoft.ContainerService/ManagedClusters, resourceType)"
	// s := "substringof('dev', name)"
	// op := armresources.ClientListByResourceGroupOptions{nil, &s, nil}
	resp := resourceClient.NewListByResourceGroupPager(rg, &armresources.ClientListByResourceGroupOptions{nil, nil, nil})
	resources := make([]*armresources.GenericResourceExpanded, 0)
	for resp.More() {
		pageResp, err := resp.NextPage(context.Background())
		onErrorFail(err, "Next Page Failed")
		resources = append(resources, pageResp.ResourceListResult.Value...)
		kJson, _ := json.Marshal(pageResp.ResourceListResult.Value)
		log.Println(string(kJson))
	}
	log.Println(resources)
}
