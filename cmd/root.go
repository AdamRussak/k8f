/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"k8f/core"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

var tversion string

// rootCmd represents the base command when called without any subcommands
var (
	o                 FlagsOptions
	supportedProvider = []string{"azure", "aws", "gcp"}
	argouments        = []string{"all"}
	AwsRegion         = "eu-west-1"
	configYAML        = "yaml"
	listOutput        = "json"
	version           = tversion
	confPath          = clientcmd.RecommendedHomeFile
	rootCmd           = &cobra.Command{
		Version: version,
		Use:     "k8f",
		Short:   "A CLI tool to List, Connect, Search and check version for K8S Clusters in all your resources at once",
		Long: `A CLI tool to find, list, connect, search and check version for K8S Clusters in all your resources at once,
this tool supports Azure AKS and AWS EKS. For example:
	to get List of all EKS:
		k8f  list aws
	to connect to all K8S:
		k8f  connect all`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	core.OnErrorFail(err, "error executing command")
}

func init() {
	rootCmd.Flags().BoolVar(&o.DryRun, "dry-run", false, "Run the task as Dry-run, no action is done")
	rootCmd.PersistentFlags().BoolVarP(&core.Verbosity, "verbose", "v", false, "verbose logging")
	rootCmd.PersistentFlags().StringVar(&AwsRegion, "aws-region", AwsRegion, "Set Default AWS Region")
	rootCmd.Flags().IntVar(&o.UiSize, "ui-size", 4, "number of list items to show in menu at once")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
