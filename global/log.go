package global

import (
	"os"

	"github.com/op/go-logging"
)

// Log is the logger used throughout the node's process.
var Log *logging.Logger

// SetUpLogger sets up the logger Log.
func SetUpLogger() {
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	format := logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05 -07:00} %{shortfunc} [%{level}]%{color:reset} %{message}`,
	)
	logging.SetBackend(backend)
	logging.SetFormatter(format)

	Log = logging.MustGetLogger("GoRaft")
}

// SetLogLevel sets the log level based on the config.yaml log_level value.
func SetLogLevel(level string) {
	logLevel, err := logging.LogLevel(level)
	if err != nil {
		Log.Panic(err)
	}
	logging.SetLevel(logLevel, "GoRaft")
}
