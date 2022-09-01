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
		Short: "A ClI tool to List, Connect, Search and Monitor K8S Clusters",
		Long: `A longer description that spans multiple lines and likely contains
		examples and usage of using your application. For example:

		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.`,
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
