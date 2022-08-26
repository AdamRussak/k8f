/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var (
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get a Specific K8S in Azure/AWS or Both",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("o:", o)
			kJson, _ := json.Marshal(o)
			fmt.Println(string(kJson))
			fmt.Println("get called")
		},
	}
)

func init() {
	// getCmd.Flags().StringVarP(&o.Output, "output", "o", configYAML, "Merged kubeconfig output type(json or yaml)")
	// getCmd.Flags().StringVar(&o.Path, "path", confPath, "Merged kubeconfig output name and path")
	// getCmd.Flags().BoolVar(&o.Overwrite, "overwrite", false, "If true, force merge kubeconfig")
	// getCmd.Flags().BoolVar(&o.Backup, "backup", true, "If true, backup $HOME/.kube/config file to $HOME/.kube/config.bk")
	rootCmd.AddCommand(getCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
