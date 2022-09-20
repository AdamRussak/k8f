package provider

import (
	"context"
	"fmt"
	"k8f/core"
	"strings"

	log "github.com/sirupsen/logrus"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // register GCP auth provider
	"k8s.io/client-go/tools/clientcmd/api"
)

// main process for List Command
func (c CommandOptions) GcpMain() Provider {
	log.Info("Starting GCP List")
	var clusters []Account
	Projects := gcpProjects()
	chanel := make(chan Account)
	for _, p := range Projects {
		go func(chanel chan Account, p subs) {
			log.Info("Starting GCP project: " + p.Id)
			var err error
			projclusters, err := c.getK8sClusterConfigs(context.Background(), p.Id)
			chanel <- Account{Name: p.Id, Clusters: projclusters, TotalCount: len(projclusters)}
			log.Error(err)
		}(chanel, p)

	}
	for i := 0; i < len(Projects); i++ {
		accoutn := <-chanel
		clusters = append(clusters, accoutn)
	}
	return Provider{Provider: "gcp", Accounts: clusters, TotalCount: countTotal(clusters)}
}

// lists all GCP projects in current orgenization
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

func (c CommandOptions) getK8sClusterConfigs(ctx context.Context, projectId string) ([]Cluster, error) {
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
		clustserss = append(clustserss, Cluster{Name: a.Name, Version: a.CurrentMasterVersion, Region: a.Zone, Latest: c.latestGCP(a)})
	}
	return clustserss, nil
}

// func to get latest version
func (c CommandOptions) latestGCP(k *container.Cluster) string {
	svc, err := container.NewService(context.Background())
	core.OnErrorFail(err, "failed to create container service")
	output := svc.Projects.Zones.GetServerconfig(k.Name, k.Zone)
	ver, err := output.Do()
	core.OnErrorFail(err, "failed to get versions")
	for _, v := range ver.Channels {
		if strings.Contains(k.ReleaseChannel.Channel, v.Channel) {
			return v.ValidVersions[0]
		}
	}
	return ""
}

// FIXME: keep setting the config: https://github.com/carlpett/gke-config-helper/blob/master/main.go
// main func for Kubeconfig
func (c CommandOptions) GetK8sClusterConfigs() *api.Config {
	svc, err := container.NewService(ctx)
	core.OnErrorFail(err, "failed to create Container Serivce")

	// Basic config structure
	ret := api.Config{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters:   map[string]*api.Cluster{},  // Clusters is a map of referencable names to cluster configs
		AuthInfos:  map[string]*api.AuthInfo{}, // AuthInfos is a map of referencable names to user configs
		Contexts:   map[string]*api.Context{},  // Contexts is a map of referencable names to context configs
	}
	//get all Projects in org
	projList := gcpProjects()
	for _, p := range projList {
		resp, err := svc.Projects.Zones.Clusters.List(p.Id, "-").Context(ctx).Do()
		if err != nil {
			log.Error("clusters list project=%s: %w", p.Id, err)
			continue
		}

		for _, f := range resp.Clusters {
			name := fmt.Sprintf("gke_%s_%s_%s", p.Id, f.Zone, f.Name)
			// cert, err := base64.StdEncoding.DecodeString(f.MasterAuth.ClusterCaCertificate)
			// core.OnErrorFail(err, "erro in certificate format")
			// example: gke_my-project_us-central1-b_cluster-1 => https://XX.XX.XX.XX
			ret.Clusters[name] = &api.Cluster{
				CertificateAuthorityData: []byte(f.MasterAuth.ClusterCaCertificate),
				Server:                   "https://" + f.Endpoint,
			}
			// Just reuse the context name as an auth name.
			ret.Contexts[name] = &api.Context{
				Cluster:  name,
				AuthInfo: name,
			}
			// GCP specific configation; use cloud platform scope.
			ret.AuthInfos[name] = &api.AuthInfo{
				AuthProvider: &api.AuthProviderConfig{
					Name: "gcp",
					Config: map[string]string{
						"scopes": "https://www.googleapis.com/auth/cloud-platform",
					},
				},
			}
		}
	}
	// Ask Google for a list of all kube clusters in the given project.

	return &ret
}
