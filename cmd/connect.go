/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"k8f/core"
	"k8f/provider"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// FEATURE: add flag to support aws commands per account/profile (for example: ask per-account what to use)
// connectCmd represents the connect command
var (
	connectCmd = &cobra.Command{
		Use:     connectCMD,
		Short:   connectShort,
		Example: connectExample,
		Run: func(cmd *cobra.Command, args []string) {
			var options provider.CommandOptions = newCommandStruct(cmd, o, args)
			log.WithField("CommandOptions", log.Fields{"struct": core.DebugWithInfo(options)}).Debug("CommandOptions Struct Keys and Values: ")
			if !options.Overwrite && core.Exists(options.Path) && !options.DryRun && !options.Backup && !options.Merge {
				core.FailOnError(errors.New("flags error"), "Cant Run command as path exist and Overwrite is set to FALSE")
			}
			if !core.Exists(options.Path) {
				core.CreateDirectory(options.Path)
				log.Warn("Path directorys were created")
			}
			if !core.Exists(options.Path) && (options.Backup || options.Merge) {
				fileCreateConfirm := provider.BoolUI("--backup/--merge option was used but "+options.Path+
					" does not exist. Choose True to create empty file, if False is chosen backup/merge will fail", options)
				if fileCreateConfirm == "True" {
					err := os.WriteFile(options.Path, []byte(""), 0666)
					core.FailOnError(err, "failed to save config")
				} else {
					core.FailOnError(errors.New("flags error"), "Cant Run command as path does not exist and Backup or Merge is set to True")
				}
			}
			if args[0] == "azure" {
				options.ConnectAllAks()
			} else if args[0] == "aws" {
				options.ConnectAllEks()
			} else if args[0] == "all" {
				log.Info("Supported Platform are:" + core.PrintOutStirng(supportedProvider))
				options.FullCloudConfig()
			} else {
				core.FailOnError(errors.New("no Provider Selected"), "Selected Provider Not avilable (yet)")
			}
		},
	}
)

func init() {
	connectCmd.Flags().StringVarP(&o.Output, "output", "o", defaultYAMLoutput, "kubeconfig output type format(json or yaml)")
	connectCmd.Flags().StringVarP(&o.Path, "path", "p", confPath, "Set output path")
	connectCmd.Flags().BoolVar(&o.ProfileSelector, "profile-select", false, "provides a UI to select a single profile to scan")
	connectCmd.Flags().BoolVar(&o.Overwrite, "overwrite", false, "If true, force overwrite kubeconfig")
	connectCmd.Flags().BoolVar(&o.DryRun, DryRun, false, "If true, only run a dry-run with cli output")
	connectCmd.Flags().BoolVar(&o.Backup, "backup", false, "If true, backup config file to $HOME/.kube/config.bk")
	connectCmd.Flags().BoolVar(&o.Merge, "merge", false, "If true, add new K8s to the existing kubeconfig path")
	connectCmd.Flags().BoolVar(&o.ForceMerge, "force-merge", false, "If set, all duplication will be merged without prompt, default is interactive")
	connectCmd.Flags().BoolVar(&o.AwsAuth, "auth", false, "change from AWS CLI Auth to AWS IAM Authenticator, Default set to AWS CLI")
	connectCmd.Flags().BoolVar(&o.AwsEnvProfile, "isEnv", false, "Add AWS Profile as Env to the Kubeconfig")
	connectCmd.Flags().BoolVarP(&o.AwsClusterName, "short-name", "s", false, "shorten EKS name from <account>:<region>:<cluster> to <region>:<cluster>")
	connectCmd.Flags().StringVar(&o.AwsRoleString, "role-name", "", "Set Role Name (Example: 'myRoleName')")
	connectCmd.MarkFlagsMutuallyExclusive(DryRun, "overwrite", "backup", "merge")
	rootCmd.AddCommand(connectCmd)
}
