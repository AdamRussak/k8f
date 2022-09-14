// /*
// Copyright © 2022 NAME HERE <EMAIL ADDRESS>

// */
package cmd

// import (
// 	"fmt"
// 	"k8f/tools"

// 	"github.com/spf13/cobra"
// )

// // ddCmd represents the dd command
// var apikey string
// var appkey string
// var ddCmd = &cobra.Command{
// 	Use:   "dd",
// 	Short: "Send Metrics to Data Dog",
// 	// Args: func(cmd *cobra.Command, args []string) error {
// 	// 	if len(args) < 1 {
// 	// 		return errors.New("requires cloud provider")
// 	// 	}
// 	// 	argouments = append(argouments, supportedProvider...)

// 	// 	if core.IfXinY(args[0], argouments) {
// 	// 		return nil
// 	// 	}
// 	// 	return fmt.Errorf("invalid cloud provider specified: %s", args[0])
// 	// },
// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("apikey: " + apikey)
// 		tools.DdMain(apikey)
// 		fmt.Println("dd called")
// 	},
// }

// func init() {
// 	ddCmd.Flags().StringVar(&apikey, "apikey", "", "Set API Key for Datadog")
// 	// ddCmd.Flags().StringVar(&appkey, "appkey", "", "Set App Key for Datadog")
// 	// ddCmd.MarkFlagsRequiredTogether("appkey", "apikey")
// 	rootCmd.AddCommand(ddCmd)
// }
