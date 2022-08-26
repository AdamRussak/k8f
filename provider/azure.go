package provider

import (
	"context"
	"encoding/json"
	"fmt"
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
	c0 := make(chan string)
	tenant := GetTenentList()
	for _, t := range tenant {
		go func(c0 chan string, t armsubscriptions.TenantIDDescription) {
			subs := listSubscriptions(*t.TenantID)
			c1 := make(chan Account)
			for _, s := range subs {
				log.Println("starting: ", s.Name)
				go getAllAKS(s, c1, *t.TenantID)
			}
			for i := 0; i < len(subs); i++ {
				res := <-c1
				list = append(list, res)
			}
			c0 <- "tenant is done"
		}(c0, t)

	}
	for i := 0; i < len(tenant); i++ {
		res := <-c0
		log.Println(res)
	}
	return Provider{"azure", list, countTotal(list)}
}

func auth(tenantid string) *azidentity.AzureCLICredential {
	cred, err := azidentity.NewAzureCLICredential(&azidentity.AzureCLICredentialOptions{TenantID: tenantid})
	core.OnErrorFail(err, "Authentication Failed")
	return cred
}

// get full list of tenants user got permissions to.
// URGENT: add multi-tenant support
func GetTenentList() []armsubscriptions.TenantIDDescription {
	var res []armsubscriptions.TenantIDDescription
	tenants, err := armsubscriptions.NewTenantsClient(auth(""), nil)
	core.OnErrorFail(err, "Failed to get Tenants")
	tenant := tenants.NewListPager(nil)
	for tenant.More() {
		nextResult, err := tenant.NextPage(ctx)
		core.OnErrorFail(err, "failed to advance page")
		for _, v := range nextResult.Value {
			kJson, _ := json.Marshal(v)
			fmt.Println(string(kJson))
			res = append(res, *v)
		}
	}
	return res
}

func listSubscriptions(id string) []subs {
	var res []subs
	client, err := armsubscriptions.NewClient(auth(id), nil)
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
func getAllAKS(subscription subs, c1 chan Account, id string) {
	var r []Cluster
	client, err := armcontainerservice.NewManagedClustersClient(subscription.Id, auth(id), nil)
	core.OnErrorFail(err, "failed to create client")
	pager := client.NewListPager(nil)
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		core.OnErrorFail(err, "failed to advance page")
		for _, v := range nextResult.Value {
			l := getAksConfig(SplitAzIDAndGiveItem(*v.ID, 4), *v.Name, subscription.Id, id)
			r = append(r, Cluster{*v.Name, *v.Properties.KubernetesVersion, l, *v.Location, *v.ID})
		}
	}
	c1 <- Account{subscription.Name, r, len(r)}
}

// Getting a single
func getAksConfig(resourceGroup string, resourceName string, subscription string, id string) string {
	var supportList []string
	client, err := armcontainerservice.NewManagedClustersClient(subscription, auth(id), nil)
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

func (c CommandOptions) ConnectAllAks() AllConfig {
	kJson, _ := json.Marshal(c)
	fmt.Println(string(kJson))
	var authe []Users
	var context []Contexts
	var clusters []Clusters
	var arnContext string
	p := FullAzureList()
	for _, a := range p.Accounts {
		chanel := make(chan AllConfig)
		for _, c := range a.Clusters {
			go func(chanel chan AllConfig, c Cluster) {
				client, err := armcontainerservice.NewManagedClustersClient(SplitAzIDAndGiveItem(c.Id, 2), auth(c.Id), nil)
				core.OnErrorFail(err, "get user creds Failed")
				chanel <- getAksProfile(client, SplitAzIDAndGiveItem(c.Id, 4), c.Name)
			}(chanel, c)
		}
		for i := 0; i < len(a.Clusters); i++ {
			response := <-chanel
			arnContext = response.context[0].Context.User
			authe = append(authe, response.auth...)
			context = append(context, response.context...)
			clusters = append(clusters, response.clusters...)
		}
	}
	if c.Combined == false {
		log.Println("Started azure only config creation")
		c.Merge(AllConfig{authe, context, clusters}, arnContext)
		return AllConfig{}
	} else {
		log.Println("Started azure combined config creation")
		return AllConfig{authe, context, clusters}
	}

}

func SplitAzIDAndGiveItem(input string, out int) string {
	s := strings.Split(input, "/")
	return s[out]
}
