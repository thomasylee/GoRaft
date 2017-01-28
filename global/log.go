package global

import (
	"os"

	"github.com/op/go-logging"
)

var Log *logging.Logger

func SetUpLogger() {
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	format := logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05 -07:00} %{shortfunc} [%{level}]%{color:reset} %{message}`,
	)
	logging.SetBackend(backend)
	logging.SetFormatter(format)

	Log = logging.MustGetLogger("GoRaft")
}

/**
 * Sets the log level based on the config.yaml log_level value.
 */
func SetLogLevel(config map[interface{}]interface{}) {
	logLevel, err := logging.LogLevel(config["log_level"].(string))
	if err != nil {
		Log.Panic(err)
	}
	logging.SetLevel(logLevel, "GoRaft")
}
