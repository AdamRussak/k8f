// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"encoding/json"
	"fmt"
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

func mainAKS() {
	var allResouces []rgAndResouce
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
		r := listResource(*resource.Name, cred)
		if len(r) != 0 {
			allResouces = append(allResouces, rgAndResouce{RGName: *resource.Name, Resources: r})
		}
	}
	kJson, _ := json.Marshal(allResouces)
	fmt.Println(string(kJson))

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

func listResource(rg string, cred azcore.TokenCredential) []resource {
	var lR []resource
	resourceClient, err := armresources.NewClient(subscriptionID, cred, nil)
	onErrorFail(err, "Failed to create Resource Client")
	s := "resourceType eq 'Microsoft.ContainerService/ManagedClusters'"
	resp := resourceClient.NewListByResourceGroupPager(rg, &armresources.ClientListByResourceGroupOptions{nil, &s, nil})
	for resp.More() {
		pageResp, err := resp.NextPage(context.Background())
		for _, res := range pageResp.ResourceListResult.Value {
			ListResources := resource{Id: *res.ID, Location: *res.Location, Name: *res.Name, Type: *res.Type}
			lR = append(lR, ListResources)
		}
		onErrorFail(err, "Next Page Failed")
	}
	return lR
}
