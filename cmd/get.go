/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

// import (
// 	"fmt"
// 	"k8f/core"
// 	"k8f/provider"
// 	"os"

// 	log "github.com/sirupsen/logrus"
// 	"github.com/spf13/cobra"
// 	"gopkg.in/yaml.v2"
// )

// // getCmd represents the get command
// var (
// 	getCmd = &cobra.Command{
// 		Use:    "get",
// 		Short:  "Get a Specific K8S in Azure/AWS or Both",
// 		PreRun: core.ToggleDebug,
// 		Run: func(cmd *cobra.Command, args []string) {
// 			options := provider.CommandOptions{}
// 			log.Debug("Get start on Debug Mode")
// 			log.Info("Get Command Starting")
// 			y, _ := yaml.Marshal(options.GetK8sClusterConfigs())
// 			err := os.WriteFile("./testfiles/gcpConfig.yaml", y, 0666)
// 			core.OnErrorFail(err, "failed to save yaml")
// 			fmt.Println("get called")
// 			m := map[string]string{"1": "a", "2": "b"}
// 			fmt.Println(provider.RunResult(m, "yaml"))
// 		},
// 	}
// )

// func init() {
// 	rootCmd.AddCommand(getCmd)
// }
