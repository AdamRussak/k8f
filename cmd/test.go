/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"k8-upgrade/provider"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "test a single K8S config creation",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Args:    func(cmd *cobra.Command, args []string) error {},
	Example: `k8-upgrade test`,
	Run:     func(cmd *cobra.Command, args []string) { provider.TetsKubeConfig() },
}

func init() {
	rootCmd.AddCommand(testCmd)
}
