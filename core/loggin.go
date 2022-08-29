package core

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"github.com/spf13/cobra"
)

var Verbosity bool

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {

	return []byte(fmt.Sprintf("%s\n", entry.Message)), nil
}
func ToggleDebug(cmd *cobra.Command, args []string) {
	formatter := new(prefixed.TextFormatter)
	formatter.ForceColors = true
	// Set specific colors for prefix and timestamp
	formatter.SetColorScheme(&prefixed.ColorScheme{
		InfoLevelStyle:  "green",
		WarnLevelStyle:  "yellow",
		ErrorLevelStyle: "red",
		FatalLevelStyle: "red",
		PanicLevelStyle: "red",
		DebugLevelStyle: "blue",
		PrefixStyle:     "cyan",
		TimestampStyle:  "black+h",
	})
	formatter.DisableTimestamp = true
	formatter.DisableUppercase = true
	if Verbosity {
		log.SetFormatter(formatter)
		log.Info("Debug logs enabled")
		log.SetLevel(log.TraceLevel)
		// log.SetReportCaller(true)
	} else {
		log.SetFormatter(formatter)
		log.SetLevel(log.InfoLevel)
	}
}
