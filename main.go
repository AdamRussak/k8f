/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
// TODO: create discovery mode (to find all known EKS Clustesr)
// TODO: add auto configuration of Regions (if not in env var or if requested auto descovery)
// TODO: cli commands to be used with Nagios (or other monitoring system)
// TODO: add Same funcionality to Azure

// TODO: add count of total clusters
// TODO: add recomendation to upgrade / everyting is ok for each cluster

// TODO:Investigate sending resoult to DD
package main

import (
	"k8-upgrade/cmd"
)

func main() {
	cmd.Execute()
	// mainAWS()
}
