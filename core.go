package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func onErrorFail(err error, message string) {
	if err != nil {
		fmt.Printf("%s: %s\n", message, err)
		os.Exit(1)
	}
}

// getEnvVarOrExit returns the value of specified environment variable or terminates if it's not defined.
func getEnvVarOrExit(varName string) string {
	value := os.Getenv(varName)
	if value == "" {
		fmt.Printf("Missing environment variable %s\n", varName)
		os.Exit(1)
	}

	return value
}
func printOutStruct(input []string) {
	kJson, err := json.Marshal(input)
	onErrorFail(err, "Json Marshal Failed")
	fmt.Println(string(kJson))
}
