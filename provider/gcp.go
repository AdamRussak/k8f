package provider

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"k8f/core"
	"log"

	"google.golang.org/api/container/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // register GCP auth provider
)

var fProjectId = "playground-s-11-8be6d443"

func GcpMain() {
	flag.Parse()
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
	// for _, f := range resp.Clusters {
	// 	name := fmt.Sprintf("gke_%s_%s_%s", projectId, f.Zone, f.Name)
	// 	cert, err := base64.StdEncoding.DecodeString(f.MasterAuth.ClusterCaCertificate)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("invalid certificate cluster=%s cert=%s: %w", name, f.MasterAuth.ClusterCaCertificate, err)
	// 	}
	// 	// example: gke_my-project_us-central1-b_cluster-1 => https://XX.XX.XX.XX
	// 	ret.Clusters[name] = &api.Cluster{
	// 		CertificateAuthorityData: cert,
	// 		Server:                   "https://" + f.Endpoint,
	// 	}
	// 	// Just reuse the context name as an auth name.
	// 	ret.Contexts[name] = &api.Context{
	// 		Cluster:  name,
	// 		AuthInfo: name,
	// 	}
	// 	// GCP specific configation; use cloud platform scope.
	// 	ret.AuthInfos[name] = &api.AuthInfo{
	// 		AuthProvider: &api.AuthProviderConfig{
	// 			Name: "gcp",
	// 			Config: map[string]string{
	// 				"scopes": "https://www.googleapis.com/auth/cloud-platform",
	// 			},
	// 		},
	// 	}
	// }

	return nil
}
