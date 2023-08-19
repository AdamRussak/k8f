package core

import (
	"fmt"
	"os"

	"github.com/fatih/structs"
	"github.com/shiena/ansicolor"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var Verbosity bool
var ErrorLevel bool

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s\n", entry.Message)), nil
}

func ToggleDebug() {
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
	log.SetFormatter(formatter)
	log.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
	if Verbosity {
		log.Info("Debug logs enabled")
		log.SetLevel(log.TraceLevel)
	} else if ErrorLevel {
		log.SetLevel(log.ErrorLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func DebugWithInfo(input interface{}) map[string]interface{} {
	return structs.Map(input)
}
