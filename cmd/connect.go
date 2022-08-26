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

// connectCmd represents the connect command
var (
	connectCmd = &cobra.Command{
		Use:   "connect",
		Short: "Connect to a specific cluster",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires cloud provider")
			}
			if core.IfXinY(args[0], []string{"azure", "aws", "all"}) {
				return nil
			}
			return fmt.Errorf("invalid cloud provider specified: %s", args[0])
		},
		Run: func(cmd *cobra.Command, args []string) {
			options := provider.CommandOptions{Path: o.Path, Output: o.Output, Overwrite: o.Overwrite, Combined: core.BoolCombine(args[0], supportedProvider), Backup: o.Backup, DryRun: o.DryRun, Version: o.Version}
			if args[0] == "azure" {
				options.ConnectAllAks()
			} else if args[0] == "aws" {
				options.ConnectAllEks()
			} else if args[0] == "all" {
				options.FullCloudConfig()
			} else {
				core.OnErrorFail(errors.New("no Provider Selected"), "Selected Provider Not avilable (yet)")
			}
		},
	}
)

func init() {
	connectCmd.Flags().StringVarP(&o.Output, "output", "o", configYAML, "Merged kubeconfig output type(json or yaml)")
	connectCmd.Flags().StringVar(&o.Path, "path", confPath, "Merged kubeconfig output name and path")
	connectCmd.Flags().BoolVar(&o.Overwrite, "overwrite", false, "If true, force merge kubeconfig")
	connectCmd.Flags().BoolVar(&o.DryRun, "dry-run", false, "If true, backup $HOME/.kube/config file to $HOME/.kube/config.bk")
	rootCmd.AddCommand(connectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// connectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// connectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
