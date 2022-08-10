package provider

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-version"
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

func runResult(p Provider) string {
	kJson, _ := json.Marshal(p)
	return string(kJson)
}
