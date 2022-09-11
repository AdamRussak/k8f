package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"k8f/core"

	log "github.com/sirupsen/logrus"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // register GCP auth provider
)

func gcpProjects() []subs {
	// resource manager auth
	var projStruct []subs
	cloudresourcemanagerService, err := cloudresourcemanager.NewService(ctx, option.WithScopes(cloudresourcemanager.CloudPlatformReadOnlyScope))
	core.OnErrorFail(err, "Failed to create Auth client")
	// get list of orginization Projects
	projList := cloudresourcemanagerService.Projects.List()
	resp, err := projList.Do()
	core.OnErrorFail(err, "Failed to get projects list")
	for _, a := range resp.Projects {
		projStruct = append(projStruct, subs{Name: a.Name, Id: a.ProjectId})
	}
	return projStruct
}

// TODO: add process to Create GKE Config
func (c CommandOptions) GcpMain() {
	var clusters []Account
	Projects := gcpProjects()
	for _, p := range Projects {
		var err error
		projclusters, err := getK8sClusterConfigs(context.Background(), p.Id)
		clusters = append(clusters, Account{Name: p.Id, Clusters: projclusters, TotalCount: len(projclusters)})
		log.Error(err)
	}
	kJson, _ := json.Marshal(clusters)
	log.Info(string(kJson))
}

func getK8sClusterConfigs(ctx context.Context, projectId string) ([]Cluster, error) {
	var clustserss []Cluster
	svc, err := container.NewService(ctx)
	if err != nil {
		return []Cluster{}, fmt.Errorf("container.NewService: %w", err)
	}

	// Ask Google for a list of all kube clusters in the given project.

	resp, err := svc.Projects.Zones.Clusters.List(projectId, "-").Context(ctx).Do()
	if err != nil {
		return []Cluster{}, fmt.Errorf("clusters list project=%s: %w", projectId, err)
	}

	for _, a := range resp.Clusters {
		log.Info("the Cluster name is: " + a.Name + " and its in zone " + a.Zone)
		clustserss = append(clustserss, Cluster{Name: a.Name, Version: a.CurrentMasterVersion, CluserChannel: a.ReleaseChannel.Channel, Region: a.Zone, Latest: ""})
	}
	return clustserss, nil
}

// func to get latest version
// TODO: add process to map Latest Version per channel
func latestGCP() {
	//versions
	// fmt.Println("============================")
	// output := svc.Projects.Zones.GetServerconfig(resp.Clusters[0].Name, resp.Clusters[0].Zone)
	// fmt.Println("============================")
	// test, _ := output.Do()
	// kJson, _ := json.Marshal(test)
	// fmt.Println(string(kJson))
}
