/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	supportedProvider = []string{"azure", "aws"}
	configYAML        = "yml"
	listOutput        = "json"
	confPath          = "/tmp/test.yml"
	o                 FlagsOptions
	rootCmd           = &cobra.Command{
		Use:   "k8s-upgrade",
		Short: "A ClI tool to list, search and monitor K8S Clusters",
		// 	Long: `A longer description that spans multiple lines and likely contains
		// examples and usage of using your application. For example:

		// Cobra is a CLI library for Go that empowers applications.
		// This application is a tool to generate the needed files
		// to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
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
	rootCmd.Flags().BoolVar(&o.Version, "version", false, "Show Cli version")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
