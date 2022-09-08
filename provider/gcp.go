package provider

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"k8f/core"
	"log"

	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // register GCP auth provider
)

func gcpAuth() {
	// resource manager auth
	cloudresourcemanagerService, err := cloudresourcemanager.NewService(ctx, option.WithScopes(cloudresourcemanager.CloudPlatformReadOnlyScope))
	core.OnErrorFail(err, "Failed to create Auth client")
	// get list of orginization Projects
	projList := cloudresourcemanagerService.Projects.List()
	resp, err := projList.Do()
	core.OnErrorFail(err, "Failed to get projects list")
	kJson, _ := json.Marshal(resp)
	fmt.Println(string(kJson))
}

// TODO: add process to check and list all Projects user can see
func listProjects() {}

// TODO: add process to check what projects you got and scan all of them
// TODO: add process to list all GKE
// TODO: add process to map Latest Version per channel
// TODO: add process to Create GKE Config
var fProjectId = "playground-s-11-8be6d443"

func GcpMain() {
	flag.Parse()
	gcpAuth()
	if fProjectId == "" {
		log.Fatal("must specific -projectId")
	}

	err := getK8sClusterConfigs(context.Background(), fProjectId)
	if err != nil {
		core.OnErrorFail(err, "error in GCP")
	}
}

func getK8sClusterConfigs(ctx context.Context, projectId string) error {
	svc, err := container.NewService(ctx)
	if err != nil {
		return fmt.Errorf("container.NewService: %w", err)
	}

	// Ask Google for a list of all kube clusters in the given project.
	resp, err := svc.Projects.Zones.Clusters.List(projectId, "-").Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("clusters list project=%s: %w", projectId, err)
	}
	fmt.Println("============================")
	output := svc.Projects.Zones.GetServerconfig(resp.Clusters[0].Name, resp.Clusters[0].Zone)
	fmt.Println("============================")
	test, _ := output.Do()
	kJson, _ := json.Marshal(test)
	fmt.Println(string(kJson))
	return nil
}
