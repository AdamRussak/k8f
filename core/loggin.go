package core

import (
	"fmt"
	"os"

	"github.com/fatih/structs"
	"github.com/shiena/ansicolor"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var Verbosity bool

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s\n", entry.Message)), nil
}

func ToggleDebug(cmd *cobra.Command, args []string) {
	formatter := new(prefixed.TextFormatter)
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
	formatter.ForceColors = true
	formatter.ForceFormatting = true

	if Verbosity {
		log.SetFormatter(formatter)
		log.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
		log.Info("Debug logs enabled")
		log.SetLevel(log.TraceLevel)
		// log.SetReportCaller(true)
	} else {
		log.SetFormatter(formatter)
		log.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
		log.SetLevel(log.InfoLevel)
	}
}

func DebugWithInfo(input interface{}) map[string]interface{} {
	return structs.Map(input)
}
