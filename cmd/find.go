// /*
// Copyright © 2022 NAME HERE <EMAIL ADDRESS>

// */
package cmd

import (
	"errors"
	"fmt"
	"k8f/core"
	"k8f/provider"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:     "find",
	Short:   "Find if a specific K8S exist in Azure or AWS",
	Example: `k8f find {aws/azure/all} my-k8s-cluster`,
	PreRun:  core.ToggleDebug,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			fmt.Println(len(args))
			return errors.New("requires both cloud provider & cluster name")
		}
		argouments = append(argouments, supportedProvider...)
		if core.IfXinY(args[0], argouments) {
			return nil
		}
		return fmt.Errorf("invalid cloud provider specified: %s", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		var p provider.Cluster
		options := provider.CommandOptions{AwsRegion: AwsRegion, Path: o.Path, Output: o.Output, Overwrite: o.Overwrite, Combined: core.BoolCombine(args[0], supportedProvider), Backup: o.Backup, DryRun: o.DryRun, AwsAuth: o.AwsAuth, AwsRoleString: o.AwsRoleString, AwsEnvProfile: o.AwsEnvProfile}
		log.WithField("CommandOptions", log.Fields{"struct": core.DebugWithInfo(options)}).Debug("CommandOptions Struct Keys and Values: ")
		log.Info("find called")
		if args[0] == "azure" {
			p = options.GetSingleAzureCluster(args[1])
		} else if args[0] == "aws" {
			p = options.GetSingleAWSCluster(args[1])
			// TODO: add find single cluster to GCP
			// TODO: add find single cluster to all (in case we know name but not platform)
		} else if args[0] == "all" {
			log.Info("Supported Platform are:" + core.PrintOutStirng(supportedProvider))
			// p = options.FullCloudConfig()
		} else {
			core.OnErrorFail(errors.New("no Provider Selected"), "Selected Provider Not avilable (yet)")
		}
		log.Debug(string("Outputing List as " + options.Output + " Format"))
		fmt.Println(provider.PrintoutResults(p, options.Output))
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
}
