package core

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

func FailOnError(err error, message string) {
	if err != nil {
		log.Errorf("%s: %s\n", message, err)
		os.Exit(1)
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
		FailOnError(err, "Issue setting the 'AWS_REGION' Enviroment Variable")
	}
}

func PrintOutStirng(arrayOfStrings []string) string {
	var s string
	for _, t := range arrayOfStrings {
		s = s + " " + t
	}
	return s
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
		FailOnError(err, "Failed to Create Directory")
	}

}
func MergeINIFiles(inputPaths []string) (*bytes.Reader, error) {
	// Create a buffer to store the merged result
	outputBuffer := bytes.Buffer{}

	// Iterate over input INI files
	for _, inputPath := range inputPaths {
		// Open the input INI file
		inputFile, err := ini.Load(inputPath)
		FailOnError(err, "failed to load INI")

		// Iterate over sections in the input file
		for _, section := range inputFile.Sections() {
			outputBuffer.WriteString(fmt.Sprintf("[%s]\n", section.Name()))

			// Iterate over keys in the section
			for _, key := range section.Keys() {
				outputBuffer.WriteString(fmt.Sprintf("%s = %s\n", key.Name(), key.Value()))
			}

			outputBuffer.WriteString("\n")
		}

	}
	return bytes.NewReader(outputBuffer.Bytes()), nil
}
