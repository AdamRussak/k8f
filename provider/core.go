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

func (c CommandOptions) Merge(configs AllConfig, arn string) {
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

func (c CommandOptions) FullCloudConfig() {
	var auth []Users
	var context []Contexts
	var clusters []Clusters
	r := make(chan AllConfig)
	for _, cloud := range []string{"azure", "aws"} {
		go func(cloud string, r chan AllConfig, c CommandOptions) {
			var res AllConfig
			if cloud == "azure" {
				res = c.ConnectAllAks()
			} else if cloud == "aws" {
				res = c.ConnectAllEks()
			}
			r <- res
		}(cloud, r, c)
	}
	for i := 0; i < len([]string{"azure", "aws"}); i++ {
		response := <-r
		auth = append(auth, response.auth...)
		context = append(context, response.context...)
		clusters = append(clusters, response.clusters...)
	}
	c.Merge(AllConfig{auth: auth, context: context, clusters: clusters}, context[0].Context.User)
}
