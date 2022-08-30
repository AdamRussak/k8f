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

// connectCmd represents the connect command
var (
	connectCmd = &cobra.Command{
		Use:   "connect",
		Short: "Connect to the clusters of a provider or all Supported Providers",
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
		PreRun: core.ToggleDebug,
		Run: func(cmd *cobra.Command, args []string) {
			options := provider.CommandOptions{Path: o.Path, Output: o.Output, Overwrite: o.Overwrite, Combined: core.BoolCombine(args[0], supportedProvider), Backup: o.Backup, DryRun: o.DryRun, Version: o.Version, AwsRegion: AwsRegion}
			log.WithField("CommandOptions", log.Fields{"struct": core.DebugWithInfo(options)}).Debug("CommandOptions Struct Keys and Values: ")
			if !options.Overwrite && core.Exists(options.Path) {
				core.OnErrorFail(errors.New("flags error"), "Cant Run command as path exist and Overwrite is set to FALSE")
			}
			if !core.Exists(options.Path) {
				log.Warn("Path Created")
				core.CreatDIrectoryt(options.Path)
			}
			if args[0] == "azure" {
				options.ConnectAllAks()
			} else if args[0] == "aws" {
				options.ConnectAllEks()
			} else if args[0] == "all" {
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
	connectCmd.Flags().BoolVar(&o.Backup, "Backup", false, "If true, backup config file to $HOME/.kube/config.bk")
	// rootCmd.MarkFlagsRequiredTogether("username", "password")
	connectCmd.MarkFlagsMutuallyExclusive("dry-run", "overwrite")
	connectCmd.MarkFlagsMutuallyExclusive("dry-run", "Backup")
	rootCmd.AddCommand(connectCmd)
}
