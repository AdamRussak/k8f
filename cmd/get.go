/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"k8-upgrade/core"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var (
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get a Specific K8S in Azure/AWS or Both",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		PreRun: core.ToggleDebug,
		Run: func(cmd *cobra.Command, args []string) {
			log.Debug("Get start on Debug Mode")
			log.Info("Get Command Starting")
			fmt.Println("get called")
		},
	}
)

func init() {
	rootCmd.AddCommand(getCmd)
}
