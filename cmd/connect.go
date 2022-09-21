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

// connectCmd represents the connect command
var (
	connectCmd = &cobra.Command{
		Use:   "connect",
		Short: "Connect to all the clusters of a provider or all Supported Providers",
		Example: `k8f connect aws -p ./testfiles/config --backup -v
k8f connect aws --isEnv -p ./testfiles/config --overwrite --backup --role-name "test role" -v`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires cloud provider")
			}
			argouments = append(argouments, supportedProvider...)
			if core.IfXinY(args[0], argouments) {
				return nil
			}
			return fmt.Errorf("invalid cloud provider specified: %s", args[0])
		},
		PreRun: core.ToggleDebug,
		Run: func(cmd *cobra.Command, args []string) {
			options := provider.CommandOptions{AwsRegion: AwsRegion, Path: o.Path, Output: o.Output, Overwrite: o.Overwrite, Combined: core.BoolCombine(args[0], supportedProvider), Backup: o.Backup, DryRun: o.DryRun, AwsAuth: o.AwsAuth, AwsRoleString: o.AwsRoleString, AwsEnvProfile: o.AwsEnvProfile}
			log.WithField("CommandOptions", log.Fields{"struct": core.DebugWithInfo(options)}).Debug("CommandOptions Struct Keys and Values: ")
			if !options.Overwrite && core.Exists(options.Path) && !options.DryRun && !options.Backup {
				core.OnErrorFail(errors.New("flags error"), "Cant Run command as path exist and Overwrite is set to FALSE")
			}
			if !core.Exists(options.Path) {
				core.CreatDIrectoryt(options.Path)
				log.Warn("Path directorys were created")
			}
			if args[0] == "azure" {
				options.ConnectAllAks()
			} else if args[0] == "aws" {
				options.ConnectAllEks()
			} else if args[0] == "all" {
				log.Info("Supported Platform are:" + core.PrintOutStirng(supportedProvider))
				options.FullCloudConfig()
			} else {
				core.OnErrorFail(errors.New("no Provider Selected"), "Selected Provider Not avilable (yet)")
			}
		},
	}
)

func init() {
	connectCmd.Flags().StringVarP(&o.Output, "output", "o", configYAML, "Merged kubeconfig output type(json or yaml)")
	connectCmd.Flags().StringVarP(&o.Path, "path", "p", confPath, "Merged kubeconfig output path")
	connectCmd.Flags().BoolVar(&o.Overwrite, "overwrite", false, "If true, force merge kubeconfig")
	connectCmd.Flags().BoolVar(&o.DryRun, "dry-run", false, "If true, only run a dry-run with cli output")
	connectCmd.Flags().BoolVar(&o.Backup, "backup", false, "If true, backup config file to $HOME/.kube/config.bk")
	connectCmd.Flags().BoolVar(&o.AwsAuth, "auth", false, "change from CLI Auth to AMI Auth, Default set to CLI")
	connectCmd.Flags().BoolVar(&o.AwsEnvProfile, "isEnv", false, "Add AWS Env Profile to the AWS Config")
	connectCmd.Flags().StringVar(&o.AwsRoleString, "role-name", "", "Set Role Name (Example: '')")
	connectCmd.MarkFlagsMutuallyExclusive("dry-run", "overwrite")
	connectCmd.MarkFlagsMutuallyExclusive("dry-run", "backup")
	rootCmd.AddCommand(connectCmd)
}
