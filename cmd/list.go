/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"k8-upgrade/core"
	"k8-upgrade/provider"

	log "github.com/sirupsen/logrus"
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
		if core.IfXinY(args[0], []string{"azure", "aws", "all"}) {
			return nil
		}
		return fmt.Errorf("invalid cloud provider specified: %s", args[0])
	},
	Example: `k8-upgrade list {aws/azure}`,
	PreRun:  core.ToggleDebug,
	Run: func(cmd *cobra.Command, args []string) {
		var list []provider.Provider
		var p interface{}
		options := provider.CommandOptions{Path: o.Path, Output: o.Output, Overwrite: o.Overwrite, Combined: core.BoolCombine(args[0], supportedProvider), Backup: o.Backup, DryRun: o.DryRun, Version: o.Version, AwsRegion: AwsRegion}
		log.Debug("CommandOptions Used")
		if args[0] == "azure" {
			log.Debug("Starting Azure List")
			p = options.FullAzureList()
		} else if args[0] == "aws" {
			log.Debug("Starting AWS List")
			p = options.FullAwsList()
		} else if args[0] == "all" {
			log.Debug("Starting All List")

			c0 := make(chan provider.Provider)
			for _, s := range supportedProvider {
				log.Debug(string("Starting " + s + " Provider"))
				go func(c0 chan provider.Provider, s string) {
					var r provider.Provider
					if s == "azure" {
						log.Trace(string("triggered " + s))
						r = options.FullAzureList()
					} else if s == "aws" {
						log.Trace(string("triggered " + s))
						r = options.FullAwsList()
					}
					c0 <- r
				}(c0, s)
			}
			for i := 0; i < len(supportedProvider); i++ {
				res := <-c0
				log.Trace(string("Recived A response from: " + supportedProvider[i]))
				list = append(list, res)
			}
			p = list
		} else {
			core.OnErrorFail(errors.New("no Provider Selected"), "Selected Provider Not avilable (yet)")
		}
		log.Debug(string("Outputing List as " + options.Output + " Format"))
		fmt.Println(provider.RunResult(p, options.Output))

	},
}

func init() {
	listCmd.Flags().StringVarP(&o.Output, "output", "o", listOutput, "Set output type(json or yaml)")
	rootCmd.AddCommand(listCmd)

}
