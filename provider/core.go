package provider

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-version"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// evaluate latest version from addon version list
func evaluateVersion(list []string) string {
	var latest string
	for _, v := range list {
		var lt *version.Version
		var err error
		v1, err := version.NewVersion(v)
		if err != nil {
			fmt.Println(err)
		}
		if latest == "" {
			lt, err = version.NewVersion("0.0")
		} else {
			lt, err = version.NewVersion(latest)
		}
		if err != nil {
			fmt.Println(err)
		}
		// Options availabe
		if v1.GreaterThan(lt) {
			latest = v
		} // GreaterThen
	}
	return latest
}

func RunResult(p Provider) string {
	kJson, _ := json.Marshal(p)
	return string(kJson)
}

func countTotal(f []Account) int {
	var count int
	for _, a := range f {
		count = count + a.TotalCount
	}
	return count
}

func Merge(configs AllConfig, arn string) {
	clientConfig := clientcmdapi.Config{
		Kind:           "Config",
		APIVersion:     "v1",
		Clusters:       configs.Clusters,
		Contexts:       configs.Contexts,
		CurrentContext: arn,
		AuthInfos:      configs.Authinfos,
	}
	clientcmd.WriteToFile(clientConfig, "testconfig/fullConfig")
}
