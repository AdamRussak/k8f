/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
// TODO: add verbosity to the logs (info,debug,error): https://github.com/sirupsen/logrus#:~:text=Note%20that%20it%27s%20completely%20api%2Dcompatible%20with%20the%20stdlib%20logger%2C%20so%20you%20can%20replace%20your%20log%20imports%20everywhere%20with%20log%20%22github.com/sirupsen/logrus%22%20and%20you%27ll%20now%20have%20the%20flexibility%20of%20Logrus.%20You%20can%20customize%20it%20all%20you%20want%3A
// TODO: add auto configuration of Regions (if not in env var or if requested auto descovery)
// TODO: cli commands to be used with DD & Nagios (or other monitoring system)
// TODO: add all EKS/AKS to kubeconfig
// TODO: add recomendation to upgrade / everyting is ok for each cluster
// TODO: add output flag for: json,csv,toTerminal,pdf
// TODO: investigate on GOlang AUTOMAPER

package main

import (
	"k8-upgrade/cmd"
)

func main() {
	cmd.Execute()
}
