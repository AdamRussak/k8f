package provider

import (
	"context"
	"k8f/core"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	ctx = context.Background()
)

func (c CommandOptions) FullAzureList() Provider {
	log.Info("Starting Azure Full List")
	var list []Account
	c0 := make(chan string)
	tenant := GetTenentList()
	for _, t := range tenant {
		log.Info("Start Tenanat: " + *t.DisplayName)
		go func(c0 chan string, t armsubscriptions.TenantIDDescription) {
			subs := listSubscriptions(*t.TenantID)
			c1 := make(chan Account)
			for _, s := range subs {
				log.Info("Start Subscription: " + s.Name)
				go getAllAKS(s, c1, *t.TenantID)
			}
			for i := 0; i < len(subs); i++ {
				res := <-c1
				list = append(list, res)
				log.Debug("Finished Subscription: " + subs[i].Name)
			}
			c0 <- "Finished Tenanat:"
		}(c0, t)

	}
	for i := 0; i < len(tenant); i++ {
		res := <-c0
		log.Debug(res + " " + *tenant[i].DisplayName)
	}
	return Provider{"azure", list, countTotal(list)}
}

func auth(tenantid string) *azidentity.AzureCLICredential {
	log.Debug("Start Authentication for tenant ID: " + tenantid)
	cred, err := azidentity.NewAzureCLICredential(&azidentity.AzureCLICredentialOptions{TenantID: tenantid})
	core.OnErrorFail(err, "Authentication Failed")
	log.Debug("Finished Authentication for tenant ID: " + tenantid)
	return cred
}

// get full list of tenants user got permissions to.
func GetTenentList() []armsubscriptions.TenantIDDescription {
	log.Debug("Start getting tenant list")
	var res []armsubscriptions.TenantIDDescription
	tenants, err := armsubscriptions.NewTenantsClient(auth(""), nil)
	core.OnErrorFail(err, "Failed to get Tenants")
	tenant := tenants.NewListPager(nil)
	for tenant.More() {
		nextResult, err := tenant.NextPage(ctx)
		core.OnErrorFail(err, "failed to advance page")
		for _, v := range nextResult.Value {
			res = append(res, *v)
		}
	}
	log.Debug("Finished getting tenant list")
	return res
}

func listSubscriptions(id string) []subs {
	log.Debug("Start getting Subscription list")
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
	log.Debug("Finished getting Subscription list")
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
			supportedAKS := findSupportedAksVersions(SplitAzIDAndGiveItem(*v.ID, "/", 4), *v.Name, subscription.Id, id)
			l := getAksConfig(supportedAKS)
			r = append(r, Cluster{*v.Name, *v.Properties.KubernetesVersion, l, *v.Location, *v.ID, "", microsoftSupportedVersion(l, *v.Properties.KubernetesVersion)})
		}
	}
	c1 <- Account{subscription.Name, r, len(r), id}
}

// Getting a single
func getAksConfig(supportedList []string) string {
	return evaluateVersion(supportedList)
}

func findSupportedAksVersions(resourceGroup string, resourceName string, subscription string, id string) []string {
	var supportList []string
	log.WithField("CommandOptions", log.Fields{"subscription": subscription, "tenantID": id, "resourceName": resourceName}).Debug("getAksConfig Variables and Values: ")
	client, err := armcontainerservice.NewManagedClustersClient(subscription, auth(id), nil)
	core.OnErrorFail(err, "Create Client Failed")
	profile, err := client.GetUpgradeProfile(ctx, resourceGroup, resourceName, nil)
	core.OnErrorFail(err, "Update Profile Failed")
	for _, a := range profile.Properties.ControlPlaneProfile.Upgrades {
		supportList = append(supportList, *a.KubernetesVersion)
	}
	log.Debug("List of Supported Versions")
	log.Debug(supportList)
	return supportList
}
func getAksProfile(client *armcontainerservice.ManagedClustersClient, resourceGroupName string, resourceName string) AllConfig {
	log.WithField("CommandOptions", log.Fields{"struct": core.DebugWithInfo(client), "resourceGroupName": resourceGroupName, "resourceName": resourceName}).Debug("getAksProfile Variables and Values: ")
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
	var authe []Users
	var context []Contexts
	var clusters []Clusters
	var arnContext string
	p := c.FullAzureList()
	for _, a := range p.Accounts {
		chanel := make(chan AllConfig)
		for _, c := range a.Clusters {
			go func(chanel chan AllConfig, c Cluster, a Account) {
				log.WithField("Cluster Struct", log.Fields{"struct": core.DebugWithInfo(c), "tenentAuth": core.DebugWithInfo(a)}).Debug("Creating NewManagedClustersClient")
				client, err := armcontainerservice.NewManagedClustersClient(SplitAzIDAndGiveItem(c.Id, "/", 2), auth(a.Tenanat), nil)
				core.OnErrorFail(err, "get user creds Failed")
				chanel <- getAksProfile(client, SplitAzIDAndGiveItem(c.Id, "/", 4), c.Name)
			}(chanel, c, a)
		}
		for i := 0; i < len(a.Clusters); i++ {
			response := <-chanel
			arnContext = response.context[0].Context.User
			authe = append(authe, response.auth...)
			context = append(context, response.context...)
			clusters = append(clusters, response.clusters...)
		}
	}
	if !c.Combined {
		log.Debug("Started azure only config creation")
		c.Merge(AllConfig{authe, context, clusters}, arnContext)
		return AllConfig{}
	} else {
		log.Debug("Started azure combined config creation")
		return AllConfig{authe, context, clusters}
	}

}

// func (c CommandOptions) GetSingleAzureCluster(clusterToFind string) Cluster {
// 	log.Info("Starting Azure find cluster named: " + clusterToFind)
// 	var list Cluster
// 	c0 := make(chan string)
// 	tenant := GetTenentList()
// 	for _, t := range tenant {
// 		log.Info("Start Tenanat: " + *t.DisplayName)
// 		go func(c0 chan string, t armsubscriptions.TenantIDDescription) {
// 			subs := listSubscriptions(*t.TenantID)
// 			c1 := make(chan Account)
// 			for _, s := range subs {
// 				log.Info("Start Subscription: " + s.Name)
// 				go getAllAKS(s, c1, *t.TenantID)
// 			}
// 			for i := 0; i < len(subs); i++ {
// 				res := <-c1
// 				if condition {

// 				}
// 				list = append(list, res)
// 				log.Debug("Finished Subscription: " + subs[i].Name)
// 			}
// 			c0 <- "Finished Tenanat:"
// 		}(c0, t)

// 	}
// 	for i := 0; i < len(tenant); i++ {
// 		res := <-c0
// 		log.Debug(res + " " + *tenant[i].DisplayName)
// 	}
// 	return Provider{"azure", list, countTotal(list)}

// }
