/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"k8-upgrade/core"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

// rootCmd represents the base command when called without any subcommands
var (
	supportedProvider = []string{"azure", "aws"}
	AwsRegion         = "eu-west-1"
	configYAML        = "yml"
	listOutput        = "json"
	confPath          = clientcmd.RecommendedHomeFile
	o                 FlagsOptions
	rootCmd           = &cobra.Command{
		Use:   "k8s-upgrade",
		Short: "A CLI tool to List, Connect, Search and check version for K8S Clusters in all your resources at once",
		Long: `A CLI tool to find, list, connect, search and check version for K8S Clusters in all your resources at once,
this tool supports Azure AKS and AWS EKS . For example:
	to get List of all EKS:
		k8-upgrade  list aws
	to connect to all K8S:
		k8-upgrade  connect all`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&o.DryRun, "dry-run", false, "Run the task as Dry-run, no action is done")
	rootCmd.PersistentFlags().BoolVarP(&core.Verbosity, "verbose", "v", false, "verbose logging")
	rootCmd.PersistentFlags().StringVar(&AwsRegion, "aws-region", AwsRegion, "Set Default AWS Region")

	rootCmd.Flags().BoolVar(&o.Version, "version", false, "Show Cli version")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
