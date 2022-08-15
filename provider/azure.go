package provider

import (
	"context"
	"k8-upgrade/core"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"gopkg.in/yaml.v2"
)

var (
	ctx = context.Background()
)

func FullAzureList() Provider {
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
	return Provider{"azure", list, countTotal(list)}
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
			l := getAksConfig(SplitAzIDAndGiveItem(*v.ID, 4), *v.Name, subscription.Id)
			r = append(r, Cluster{*v.Name, *v.Properties.KubernetesVersion, l, *v.Location, *v.ID})
		}
	}
	c1 <- Account{subscription.Name, r, len(r)}
}

// Getting a single
func getAksConfig(resourceGroup string, resourceName string, subscription string) string {
	var supportList []string
	client, err := armcontainerservice.NewManagedClustersClient(subscription, auth(), nil)
	core.OnErrorFail(err, "Create Client Failed")
	profile, err := client.GetUpgradeProfile(ctx, resourceGroup, resourceName, nil)
	core.OnErrorFail(err, "Update Profile Failed")
	for _, a := range profile.Properties.ControlPlaneProfile.Upgrades {
		supportList = append(supportList, *a.KubernetesVersion)
	}
	return evaluateVersion(supportList)
}

func getAksProfile(client *armcontainerservice.ManagedClustersClient, resourceGroupName string, resourceName string) AllConfig {
	l, err := client.ListClusterUserCredentials(ctx, resourceGroupName, resourceName, nil)
	core.OnErrorFail(err, "get user creds Failed")
	y := Config{}
	for _, c := range l.Kubeconfigs {
		b64 := c.Value
		err := yaml.Unmarshal(b64, &y)
		core.OnErrorFail(err, "Failed To Unmarshal Config")
	}
	return AllConfig{auth: y.Users, context: y.Contexts, clusters: y.Clusters}
}

func ConnectAllAks(combined string) (AllConfig, string) {
	var authe []Users
	var context []Contexts
	var clusters []Clusters
	var arnContext string
	p := FullAzureList()
	for _, a := range p.Accounts {
		for _, c := range a.Clusters {
			client, err := armcontainerservice.NewManagedClustersClient(SplitAzIDAndGiveItem(c.Id, 2), auth(), nil)
			core.OnErrorFail(err, "get user creds Failed")
			aks := getAksProfile(client, SplitAzIDAndGiveItem(c.Id, 4), c.Name)
			arnContext = c.Name
			authe = append(authe, aks.auth...)
			context = append(context, aks.context...)
			clusters = append(clusters, aks.clusters...)
		}
	}
	if combined == "azure" {
		Merge(AllConfig{authe, context, clusters}, arnContext)
		return AllConfig{}, ""
	} else {
		return AllConfig{authe, context, clusters}, arnContext
	}

}

func SplitAzIDAndGiveItem(input string, out int) string {
	s := strings.Split(input, "/")
	return s[out]
}
