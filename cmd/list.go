/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
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
		//TODO: check interfaces to replace this block
		if args[0] == "azure" {
			fmt.Println(provider.MainAKS())
		} else if args[0] == "aws" {
			fmt.Println(provider.MainAWS())
		} else if args[0] == "all" {
			c0 := make(chan string)
			for _, s := range cProviders {
				log.Println("starting: ", s)
				go func(c0 chan string, s string) {
					var r string
					if s == "azure" {
						r = provider.MainAKS()
					} else if s == "aws" {
						r = provider.MainAWS()
					}
					c0 <- r
				}(c0, s)
			}
			for i := 0; i < len(cProviders); i++ {
				res := <-c0
				list = append(list, res)
			}
			fmt.Println(list)
		} else {
			core.OnErrorFail(errors.New("no Provider Selected"), "Selected Provider Not avilable (yet)")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
