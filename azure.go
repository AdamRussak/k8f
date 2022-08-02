// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice"
)

var (
	subscriptionID  string
	groupsClient    resources.GroupsClient
	resourcesClient resources.Client
)

func mainAKS() {
	subscriptionID = getEnvVarOrExit("AZURE_SUBSCRIPTION_ID")
	kJson, _ := json.Marshal(getAllAKS())
	fmt.Println(string(kJson))

}

// this is the only path we need to get the aks, now need to get latest version.
func getAllAKS() []resource {
	var r []resource
	subscriptionID = getEnvVarOrExit("AZURE_SUBSCRIPTION_ID")
	cred, err := azidentity.NewAzureCLICredential(nil)
	onErrorFail(err, "Authentication Failed")
	ctx := context.Background()
	client, err := armcontainerservice.NewManagedClustersClient(subscriptionID, cred, nil)
	onErrorFail(err, "failed to create client")
	pager := client.NewListPager(nil)
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		onErrorFail(err, "failed to advance page")
		for _, v := range nextResult.Value {
			s := strings.Split(*v.ID, "/")
			l := getUpgrade(s[4], *v.Name)
			r = append(r, resource{*v.ID, *v.Location, *v.Name, *v.Type, *v.Properties.KubernetesVersion, l})
		}

	}
	return r
}

func getUpgrade(resourceGroup string, resourceName string) string {
	var supportList []string
	subscriptionID = getEnvVarOrExit("AZURE_SUBSCRIPTION_ID")
	cred, err := azidentity.NewAzureCLICredential(nil)
	onErrorFail(err, "Authentication Failed")
	ctx := context.Background()
	client, err := armcontainerservice.NewManagedClustersClient(subscriptionID, cred, nil)
	profile, err := client.GetUpgradeProfile(ctx, resourceGroup, resourceName, nil)
	onErrorFail(err, "Update Profile Failed")
	for _, a := range profile.Properties.ControlPlaneProfile.Upgrades {
		supportList = append(supportList, *a.KubernetesVersion)
	}
	return evaluateVersion(supportList)
}
