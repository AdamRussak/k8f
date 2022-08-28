package core

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Verbosity bool

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s\n", entry.Message)), nil
}
func ToggleDebug(cmd *cobra.Command, args []string) {
	if Verbosity {
		log.SetFormatter(&log.JSONFormatter{})
		log.Info("Debug logs enabled")
		log.SetLevel(log.TraceLevel)
		log.SetReportCaller(true)
	} else {
		log.SetFormatter(&log.TextFormatter{
			DisableColors:   true,
			PadLevelText:    true,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   false,
		})
		log.SetLevel(log.InfoLevel)
	}
}
