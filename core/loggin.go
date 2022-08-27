package core

import (
	"fmt"

	"github.com/fatih/structs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Verbosity bool

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s\n", entry.Message)), nil
}
func ToggleDebug(cmd *cobra.Command, args []string) {
	log.SetFormatter(&log.JSONFormatter{})
	if Verbosity {
		log.Info("Debug logs enabled")
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func LoggerCostumeFields(p interface{}) map[string]interface{} {
	return structs.Map(p)
}
