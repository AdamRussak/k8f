package core

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func OnErrorFail(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s\n", message, err)
	}
}

// getEnvVarOrExit returns the value of specified environment variable or terminates if it's not defined.
func CheckEnvVarOrSitIt(varName string, varKey string) {
	val, present := os.LookupEnv(varName)
	if present {
		log.Debug("Variable " + varName + " Was Pre-set with Value: " + val)
	} else {
		err := os.Setenv(varName, varKey)
		val = os.Getenv(varName)
		log.Debug("Variable " + varName + " is Set with Value: " + val)
		OnErrorFail(err, "Issue setting the 'AWS_REGION' Enviroment Variable")
	}
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
