/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"k8-upgrade/core"
	"k8-upgrade/provider"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all K8S in Azure/AWS or Both",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		c, _ := cmd.Flags().GetString("cloud")
		if c == "" {
			core.OnErrorFail(errors.New("no Provider Selected"), "No Provider Selected")
		} else if c == "aws" {
			provider.MainAWS()
		} else if c == "azure" {
			provider.MainAKS()
		} else {
			core.OnErrorFail(errors.New("no Provider Selected"), "Selected Provider Not avilable (yet)")
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().StringP("cloud", "c", "", "Select cloud provider")
}
