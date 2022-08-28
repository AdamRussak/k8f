package core

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func OnErrorFail(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s\n", message, err)
	}
}

// getEnvVarOrExit returns the value of specified environment variable or terminates if it's not defined.
func GetEnvVarOrExit(varName string) string {
	value := os.Getenv(varName)
	if value == "" {
		fmt.Printf("Missing environment variable %s\n", varName)
		os.Exit(1)
	}

	return value
}
func IfXinY(x string, y []string) bool {
	for _, t := range y {
		if x == t {
			return true
		}
	}
	return false
}

func BoolCombine(arg string, supportedProvider []string) bool {
	return !IfXinY(arg, supportedProvider)
}
