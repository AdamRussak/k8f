package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8-upgrade/core"

	"github.com/hashicorp/go-version"
	"gopkg.in/yaml.v2"
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
	clientConfig := Config{
		Kind:           "Config",
		APIVersion:     "v1",
		Clusters:       configs.clusters,
		Contexts:       configs.context,
		CurrentContext: arn,
		Preferences:    Preferences{},
		Users:          configs.auth,
	}
	y, _ := yaml.Marshal(clientConfig)
	err := ioutil.WriteFile("testconfig/words.yaml", y, 0777)
	core.OnErrorFail(err, "failed to save config")
	// clientcmd.WriteToFile(res, "testconfig/fullConfig")
}
