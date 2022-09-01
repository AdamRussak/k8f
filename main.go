/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
// TODO: cli commands to be used with DD & Nagios (or other monitoring system)

// TODO: add options for connect aws command: profile or assume role authentication

// TODO: add recomendation to upgrade / everyting is ok for each cluster
// TODO: add output flag for: json,csv,toTerminal,pdf

package main

import (
	"k8-upgrade/cmd"
)

func main() {
	cmd.Execute()
}
