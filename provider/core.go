package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"k8f/core"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/hashicorp/go-version"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func returnMinorDiff(splitcurernt, splitLatest []string) int {
	latestMinor, err := strconv.Atoi(splitLatest[1])
	core.FailOnError(err, "faild to convert string to int")
	currentMinor, err := strconv.Atoi(splitcurernt[1])
	core.FailOnError(err, "faild to convert string to int")
	getStatus := latestMinor - currentMinor
	return getStatus
}

// evaluate latest version from addon version list
func evaluateVersion(list []string) string {
	var latest string
	for _, v := range list {
		var lt *version.Version
		var err error
		v1, err := version.NewVersion(v)
		core.FailOnError(err, "Error Evaluating Version")
		if latest == "" {
			lt, err = version.NewVersion("0.0")
		} else {
			lt, err = version.NewVersion(latest)
		}
		core.FailOnError(err, "Error Evaluating Version")
		// Options availabe
		if v1.GreaterThan(lt) {
			latest = v
		} // GreaterThen
	}
	return latest
}

// Microsoft Compliance
func microsoftSupportedVersion(latest string, current string) string {
	//IMPORTANT: only supports same major at the moment!!!
	splitLatest := strings.Split(latest, ".")
	splitcurernt := strings.Split(current, ".")
	// make sure its the same major
	if splitLatest[0] == splitcurernt[0] {
		getStatus := returnMinorDiff(splitcurernt, splitLatest)
		//if its latest minor or -1, mark as ok
		if getStatus <= 1 {
			return "OK"
			//if its minor -2 show warning
		} else if getStatus > 1 && getStatus <= 2 {
			return "Warning"
			// if its minor > -2 show Critical
		} else {
			return "Critical"
		}

	}
	return "Unknown"
}

// provide version compare
func HowManyVersionsBack(versionsList []string, currentVersion string) string {
	log.Debug("current version is: " + currentVersion)
	latest := evaluateVersion(versionsList)
	splitcurernt := strings.Split(currentVersion, ".")
	splitLatest := strings.Split(latest, ".")
	//check if same major
	if splitLatest[0] == splitcurernt[0] {
		getStatus := returnMinorDiff(splitcurernt, splitLatest)
		if getStatus <= 1 {
			return "Perfect"
		} else if getStatus <= 3 {
			return "OK"
		} else {
			return "Warning"
		}
	}
	return "Critical"
}

// printout format selection
func PrintoutResults(p interface{}, output string) string {
	var kJson []byte
	log.Debug("start RunResult Func")
	if output == "json" {
		log.Info("start Json Marshal")
		kJson, _ = json.Marshal(p)
	} else if output == "yaml" {
		log.Info("start YAML Marshal")
		kJson, _ = yaml.Marshal(p)
	} else {
		return "Requested Output is not supported"
	}
	log.Info("returning Output Marshal")
	log.Debug("returning Output Marshal")
	return string(kJson)
}

// func to count ammount of Cluster in an account
func countTotal(f []Account) int {
	var count int
	for _, a := range f {
		count = count + a.TotalCount
	}
	return count
}

// func to merge kubeconfig output to singe config file
func (c CommandOptions) CombineConfigs(configs AllConfig, arn string) {
	var y []byte
	var err error
	clientConfig := Config{
		Kind:           "Config",
		APIVersion:     "v1",
		Clusters:       configs.clusters,
		Contexts:       configs.context,
		CurrentContext: arn,
		Preferences:    Preferences{},
		Users:          configs.auth,
	}
	if c.DryRun {
		log.Debugf("Dry-run will Output a %s Output", c.Output)
		fmt.Println(PrintoutResults(clientConfig, c.Output))
	} else {
		if c.Backup {
			y, _ = yaml.Marshal(clientConfig)
			log.Debug("calling copy file to bak")
			c.Configcopy()
		}
		if c.Merge {
			y, err = c.runMerge(clientConfig)
			core.FailOnError(err, "failed to merge configs")
			c.cleanFile()
		} else {
			y, _ = yaml.Marshal(clientConfig)
		}
		err := os.WriteFile(c.Path, y, 0666)
		core.FailOnError(err, "failed to save config")
		log.Infof("「 %s 」 write successful!\n", c.Path)
	}
}

func (c CommandOptions) FullCloudConfig() {
	var auth []Users
	var context []Contexts
	var clusters []Clusters
	r := make(chan AllConfig)
	for _, cloud := range []string{"azure", "aws"} {
		go func(cloud string, r chan AllConfig, c CommandOptions) {
			var res AllConfig
			if cloud == "azure" {
				res = c.ConnectAllAks()
			} else if cloud == "aws" {
				res = c.ConnectAllEks()
			}
			r <- res
		}(cloud, r, c)
	}
	for i := 0; i < len([]string{"azure", "aws"}); i++ {
		response := <-r
		auth = append(auth, response.auth...)
		context = append(context, response.context...)
		clusters = append(clusters, response.clusters...)
	}
	c.CombineConfigs(AllConfig{auth: auth, context: context, clusters: clusters}, context[0].Context.User)
}

func (c CommandOptions) Configcopy() {
	sourceFileStat, err := os.Stat(c.Path)
	core.FailOnError(err, "Issue Findign the Files in the path: "+c.Path)
	if !sourceFileStat.Mode().IsRegular() {
		core.FailOnError(err, c.Path+" is not a regular file")
	}
	source, err := os.Open(c.Path)
	core.FailOnError(err, "failed to Open target file")
	defer source.Close()
	var destination *os.File
	backupVersion := c.GetBackupFileVersion()
	if backupVersion >= 1 {
		destination, err = os.Create(c.Path + backupExtnesion + "." + fmt.Sprint(backupVersion))
		log.Debug(c.Path + backupExtnesion + "." + fmt.Sprint(backupVersion))
		core.FailOnError(err, "failed to Copy target file")
	} else {
		destination, err = os.Create(c.Path + backupExtnesion)
		core.FailOnError(err, "failed to Copy target file")
	}

	defer destination.Close()
	_, err = io.Copy(destination, source)
	core.FailOnError(err, "failed to Copy target file")
}

func (c CommandOptions) GetBackupFileVersion() int {
	dir := filepath.Dir(c.Path)
	dirFiles, err := os.ReadDir(dir)
	core.FailOnError(err, "failed to list all files in directory")
	var countBackups []string
	for _, file := range dirFiles {
		if strings.Contains(file.Name(), backupExtnesion) {
			countBackups = append(countBackups, file.Name())
			log.Debug(file.Name(), " ", file.IsDir())
		}
	}
	return len(countBackups)
}
func SplitAzIDAndGiveItem(input string, seperator string, out int) string {
	s := strings.Split(input, seperator)
	log.Debug("Split output")
	return s[out]
}

func (c CommandOptions) cleanFile() {
	// Open the file with write only mode and set the file mode to 0644
	file, err := os.OpenFile(c.Path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Truncate the file content to 0
	if err := file.Truncate(0); err != nil {
		panic(err)
	}

	// Get the file info to print the file size
	fileStat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	// Print the file size after cleaning the file
	log.Debug("File size after cleaning:", fileStat.Size())
}

func checkIfStructInit(u interface{}, key string) bool {
	v := reflect.ValueOf(u)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("yaml")
		if tag == key+",omitempty" {
			// Check if the field is set to its zero value
			if reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
				continue
			}
			// If any field has a non-zero value, return true
			return true
		}
	}
	return false
}
