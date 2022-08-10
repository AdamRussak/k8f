package provider

import (
	"context"
	"k8-upgrade/core"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
)

var (
	ctx = context.Background()
)

func MainAKS() string {
	var list []Account
	subs := listSubscriptions()
	c1 := make(chan Account)
	for _, s := range subs {
		log.Println("starting: ", s.Name)
		go getAllAKS(s, c1)
	}
	for i := 0; i < len(subs); i++ {
		res := <-c1
		list = append(list, res)
	}
	var count int
	for _, c := range list {
		count = count + c.TotalCount
	}
	return runResult(Provider{"azure", list, count})
}
func auth() *azidentity.AzureCLICredential {
	cred, err := azidentity.NewAzureCLICredential(nil)
	core.OnErrorFail(err, "Authentication Failed")
	return cred
}
func listSubscriptions() []subs {
	var res []subs
	client, err := armsubscriptions.NewClient(auth(), nil)
	core.OnErrorFail(err, "Failed to Auth")
	r := client.NewListPager(nil)
	for r.More() {
		nextResult, err := r.NextPage(ctx)
		core.OnErrorFail(err, "failed to advance page")
		for _, v := range nextResult.Value {
			res = append(res, subs{*v.DisplayName, *v.SubscriptionID})

		}
	}
	return res
}

// this is the only path we need to get the aks, now need to get latest version.
func getAllAKS(subscription subs, c1 chan Account) {
	var r []Cluster
	client, err := armcontainerservice.NewManagedClustersClient(subscription.Id, auth(), nil)
	core.OnErrorFail(err, "failed to create client")
	pager := client.NewListPager(nil)
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		core.OnErrorFail(err, "failed to advance page")
		for _, v := range nextResult.Value {
			s := strings.Split(*v.ID, "/")
			l := getUpgrade(s[4], *v.Name, subscription.Id)
			r = append(r, Cluster{*v.Name, *v.Properties.KubernetesVersion, l, *v.Location})
		}
	}
	c1 <- Account{subscription.Name, r, len(r)}
}

func getUpgrade(resourceGroup string, resourceName string, subscription string) string {
	var supportList []string
	client, err := armcontainerservice.NewManagedClustersClient(subscription, auth(), nil)
	profile, err := client.GetUpgradeProfile(ctx, resourceGroup, resourceName, nil)
	core.OnErrorFail(err, "Update Profile Failed")
	for _, a := range profile.Properties.ControlPlaneProfile.Upgrades {
		supportList = append(supportList, *a.KubernetesVersion)
	}
	return evaluateVersion(supportList)
}
