// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
)

var (
	groupsClient    resources.GroupsClient
	resourcesClient resources.Client
	ctx             = context.Background()
)

func mainAKS() {
	subs := listSubscriptions()
	for n, _ := range subs {
		log.Println("starting: ", subs[n].Name)
		l := getAllAKS(subs[n].Id)
		kJson, _ := json.Marshal(l)
		fmt.Println(string(kJson))
	}
}
func auth() *azidentity.AzureCLICredential {
	cred, err := azidentity.NewAzureCLICredential(nil)
	onErrorFail(err, "Authentication Failed")
	return cred
}
func listSubscriptions() []subs {
	var res []subs
	client, err := armsubscriptions.NewClient(auth(), nil)
	onErrorFail(err, "Failed to Auth")
	r := client.NewListPager(nil)
	for r.More() {
		nextResult, err := r.NextPage(ctx)
		onErrorFail(err, "failed to advance page")
		for _, v := range nextResult.Value {
			res = append(res, subs{*v.DisplayName, *v.SubscriptionID})

		}
	}
	return res
}

// this is the only path we need to get the aks, now need to get latest version.
func getAllAKS(subscription string) []resource {
	var r []resource
	client, err := armcontainerservice.NewManagedClustersClient(subscription, auth(), nil)
	onErrorFail(err, "failed to create client")
	pager := client.NewListPager(nil)
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		onErrorFail(err, "failed to advance page")
		for _, v := range nextResult.Value {
			s := strings.Split(*v.ID, "/")
			l := getUpgrade(s[4], *v.Name, subscription)
			r = append(r, resource{*v.ID, *v.Location, *v.Name, *v.Type, *v.Properties.KubernetesVersion, l})
		}
	}
	return r
}

func getUpgrade(resourceGroup string, resourceName string, subscription string) string {
	var supportList []string
	client, err := armcontainerservice.NewManagedClustersClient(subscription, auth(), nil)
	profile, err := client.GetUpgradeProfile(ctx, resourceGroup, resourceName, nil)
	onErrorFail(err, "Update Profile Failed")
	for _, a := range profile.Properties.ControlPlaneProfile.Upgrades {
		supportList = append(supportList, *a.KubernetesVersion)
	}
	return evaluateVersion(supportList)
}
