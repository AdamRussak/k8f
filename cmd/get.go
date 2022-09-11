/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"k8f/core"
	"k8f/provider"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var (
	getCmd = &cobra.Command{
		Use:    "get",
		Short:  "Get a Specific K8S in Azure/AWS or Both",
		PreRun: core.ToggleDebug,
		Run: func(cmd *cobra.Command, args []string) {
			options := provider.CommandOptions{}
			log.Debug("Get start on Debug Mode")
			log.Info("Get Command Starting")
			options.GcpMain()
			fmt.Println("get called")
			m := map[string]string{"1": "a", "2": "b"}
			fmt.Println(provider.RunResult(m, "yaml"))
		},
	}
)

func init() {
	rootCmd.AddCommand(getCmd)
}
