package core

import (
	"os"
	"path/filepath"

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

func Exists(path string) bool {
	log.Trace("Start Checking if Path Exist")
	_, err := os.Stat(path)
	if err == nil {
		log.Debug("Path Exist")
		return true
	}
	if os.IsNotExist(err) {
		log.Debug("Path Dose NOT Exist")
		return false
	}
	return false
}
func CreatDIrectoryt(path string) {
	// parts := strings.Split(path, string(os.PathSeparator))
	dir := filepath.Dir(path)
	var create string
	if Exists(dir) {
		log.Trace(dir + " Path Exist")
	} else {
		log.Debug("Createing Directory: " + filepath.Dir(path))
		if !filepath.IsAbs(dir) {
			log.Debug(dir + " Path is Not Absolute")
			create = "./" + dir
		} else {
			create = dir
		}
		err := os.MkdirAll(create, 0777)
		OnErrorFail(err, "Failed to Create Directory")
	}

}
