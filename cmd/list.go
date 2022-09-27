/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"k8f/core"
	"k8f/provider"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// listCmd represents the list command

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all K8S in Azure/AWS or Both",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires cloud provider")
		}
		argouments = append(argouments, supportedProvider...)
		if len(args) > 0 && len(args) <= len(argouments) {
			for a := range args {
				if !core.IfXinY(args[a], argouments) {
					return fmt.Errorf("invalid cloud provider specified: %s", args[a])
				}
			}
		}
		return nil
	},
	Example: `k8f list {aws/azure/all}`,
	PreRun:  core.ToggleDebug,
	Run: func(cmd *cobra.Command, args []string) {
		var list []provider.Provider
		var p interface{}
		options := provider.CommandOptions{Path: o.Path, Output: o.Output, Overwrite: o.Overwrite, Combined: core.BoolCombine(args[0], supportedProvider), Backup: o.Backup, DryRun: o.DryRun, AwsRegion: AwsRegion}
		log.WithField("CommandOptions", log.Fields{"struct": core.DebugWithInfo(options)}).Debug("CommandOptions Struct Keys and Values: ")
		log.Debug("CommandOptions Used")
		if len(args) == 1 && args[0] == "azure" {
			log.Debug("Starting Azure List")
			p = options.FullAzureList()
		} else if len(args) == 1 && args[0] == "aws" {
			log.Debug("Starting AWS List")
			p = options.FullAwsList()
		} else if len(args) == 1 && args[0] == "gcp" {
			log.Debug("Starting GCP List")
			p = options.GcpMain()
		} else if len(args) == 1 && args[0] == "all" {
			log.Debug("Starting All List")
			log.Info("Supported Platform are:" + core.PrintOutStirng(supportedProvider))

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
					} else if s == "gcp" {
						log.Trace(string("triggered " + s))
						r = options.GcpMain()
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
			log.Debug("Starting All List")
			log.Info("Supported Platform are:" + core.PrintOutStirng(supportedProvider))
			c0 := make(chan provider.Provider)
			for _, s := range args {
				log.Debug(string("Starting " + s + " Provider"))
				go func(c0 chan provider.Provider, s string) {
					var r provider.Provider
					if s == "azure" {
						log.Trace(string("triggered " + s))
						r = options.FullAzureList()
					} else if s == "aws" {
						log.Trace(string("triggered " + s))
						r = options.FullAwsList()
					} else if s == "gcp" {
						log.Trace(string("triggered " + s))
						r = options.GcpMain()
					}
					c0 <- r
				}(c0, s)
			}
			for i := 0; i < len(args); i++ {
				res := <-c0
				log.Trace(string("Recived A response from: " + args[i]))
				list = append(list, res)
			}
			p = list
		}
		log.Debug(string("Outputing List as " + options.Output + " Format"))
		fmt.Println(provider.RunResult(p, options.Output))
	},
}

func init() {
	listCmd.Flags().StringVarP(&o.Output, "output", "o", listOutput, "Set output type(json or yaml)")
	rootCmd.AddCommand(listCmd)

}
