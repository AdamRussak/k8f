package core

import (
	"errors"
	"fmt"
	"sort"

	"github.com/bndr/gotabulate"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// PromptUI output prompt ui
func PromptUI(label string, name string) string {
	validate := func(input string) error {
		if len(input) < 1 {
			return errors.New("context name must have more than 1 characters")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
		Default:  name,
	}
	result, err := prompt.Run()

	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

// PrintTable generate table
func PrintTable(config *clientcmdapi.Config) error {
	var table [][]string
	sortedKeys := make([]string, 0)
	for k := range config.Contexts {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	ctx := config.Contexts
	for _, k := range sortedKeys {
		namespace := "default"
		head := ""
		if config.CurrentContext == k {
			head = "*"
		}
		if ctx[k].Namespace != "" {
			namespace = ctx[k].Namespace
		}
		if config.Clusters == nil {
			continue
		}
		cluster, ok := config.Clusters[ctx[k].Cluster]
		if !ok {
			continue
		}
		conTmp := []string{head, k, ctx[k].Cluster, ctx[k].AuthInfo, cluster.Server, namespace}
		table = append(table, conTmp)
	}

	if table != nil {
		tabulate := gotabulate.Create(table)
		tabulate.SetHeaders([]string{"CURRENT", "NAME", "CLUSTER", "USER", "SERVER", "Namespace"})
		// Turn On String Wrapping
		tabulate.SetWrapStrings(true)
		// Render the table
		tabulate.SetAlign("center")
		fmt.Println(tabulate.Render("grid", "left"))
	} else {
		return errors.New("context not found")
	}
	return nil
}
