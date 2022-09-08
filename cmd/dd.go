// /*
// Copyright © 2022 NAME HERE <EMAIL ADDRESS>

// */
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ddCmd represents the dd command
var ddCmd = &cobra.Command{
	Use:   "dd",
	Short: "Send Metrics to Data Dog",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dd called")
	},
}

func init() {
	rootCmd.AddCommand(ddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
