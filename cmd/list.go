/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
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
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires cloud provider")
		}
		if core.IfXinY(args[0], []string{"azure", "aws"}) {
			return nil
		}
		return fmt.Errorf("invalid cloud provider specified: %s", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "azure" {
			provider.MainAKS()
		} else if args[0] == "aws" {
			provider.MainAWS()
		} else {
			core.OnErrorFail(errors.New("no Provider Selected"), "Selected Provider Not avilable (yet)")
		}
		// c, _ := cmd.Flags().GetString("cloud")
		// if c == "" {
		// 	core.OnErrorFail(errors.New("no Provider Selected"), "No Provider Selected")
		// } else if c == "aws" {
		// 	provider.MainAWS()
		// } else if c == "azure" {
		// 	provider.MainAKS()
		// } else {
		// 	core.OnErrorFail(errors.New("no Provider Selected"), "Selected Provider Not avilable (yet)")
		// }

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().StringP("cloud", "c", "", "Select cloud provider")
}
