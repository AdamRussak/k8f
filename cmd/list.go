/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"k8-upgrade/core"
	"k8-upgrade/provider"
	"log"

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
	Run: func(cmd *cobra.Command, args []string) {
		cProviders := []string{"azure", "aws"}
		var list []string
		var p interface{}
		options := provider.CommandOptions{Path: o.Path, Output: o.Output, Overwrite: o.Overwrite, Combined: core.BoolCombine(args[0], supportedProvider), Backup: o.Backup, DryRun: o.DryRun, Version: o.Version}
		if args[0] == "azure" {
			p = options.FullAzureList()
		} else if args[0] == "aws" {
			p = options.FullAwsList()
		} else if args[0] == "all" {
			c0 := make(chan provider.Provider)
			for _, s := range cProviders {
				log.Println("starting: ", s)
				go func(c0 chan provider.Provider, s string) {
					var r provider.Provider
					if s == "azure" {
						r = options.FullAzureList()
					} else if s == "aws" {
						r = options.FullAwsList()
					}
					c0 <- r
				}(c0, s)
			}
			for i := 0; i < len(cProviders); i++ {
				res := <-c0
				kJson, _ := json.Marshal(res)
				list = append(list, string(kJson))
			}
			p = list
		} else {
			core.OnErrorFail(errors.New("no Provider Selected"), "Selected Provider Not avilable (yet)")
		}
		fmt.Println(provider.RunResult(p, options.Output))

	},
}

func init() {
	listCmd.Flags().StringVarP(&o.Output, "output", "o", listOutput, "Set output type(json or yaml)")
	rootCmd.AddCommand(listCmd)

}
