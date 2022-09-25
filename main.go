/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
// TODO: cli commands to be used with DD & Nagios (or other monitoring system)

// TODO: add output flag for: json,csv,toTerminal,pdf
// PRIORITY: add FIND command to look for a specific Cluster accross Accounts (need to set cloud provder and cluster name)
// PRIORITY: fix AWS Role output format (just role name is needed to build ARN)
// PRIORITY: update ReadMe and ChangeLOG.MD
// testing GPG
package main

import (
	"k8f/cmd"
)

func main() {
	cmd.Execute()
}
