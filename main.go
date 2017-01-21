package main

import (
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/op/go-logging"

	"github.com/thomasylee/GoRaft/api"
	"github.com/thomasylee/GoRaft/state"
)

/**
 * The leveled Logger to use throughout GoRaft.
 */
var Log *logging.Logger

/**
 * Returns a timeout value between average - jitter and average + jitter.
 */
func calcTimeout(average int, jitter int) int {
	if jitter == 0 {
		return average
	}

	return average - jitter + rand.Intn(2 * jitter)
}

/**
 * Returns the Logger to use throughout GoRaft.
 */
func setUpLogger() *logging.Logger {
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	format := logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05 -07:00} %{shortfunc} [%{level}]%{color:reset} %{message}`,
	)
	logging.SetBackend(backend)
	logging.SetFormatter(format)

	return logging.MustGetLogger("GoRaft")
}

/**
 * Sets the log level based on the config.yaml log_level value.
 */
func setLogLevel(config map[interface{}]interface{}) {
	logLevel, err := logging.LogLevel(config["log_level"].(string))
	if err != nil {
		Log.Panic(err)
	}
	logging.SetLevel(logLevel, "GoRaft")
}

/**
 * Load the config map from config.yaml.
 */
func loadConfig() map[interface{}]interface{} {
	configData, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		Log.Panic(err)
	}
	configYaml := string(configData)

	config := make(map[interface{}]interface{})
	yaml.Unmarshal([]byte(configYaml), &config)
	return config
}

/**
 * Runs the infinite loop that keeps the node active.
 */
func runNode(config map[interface{}]interface{}) {
	timeoutChannel := make(chan bool)
	go api.RunServer(Log, timeoutChannel, config["api_port"].(int))

	// Randomize the election timeout to minimize the risk of two nodes
	// initiating an election at the same time.
	electionTimeout := config["election_timeout"].(int)
	electionTimeoutJitter := config["election_timeout_jitter"].(int)
	for {
		timeout := calcTimeout(electionTimeout, electionTimeoutJitter)
		select {
		case <-timeoutChannel:
			// Do nothing since we didn't time out.
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			// TODO: Start leader election process.
		}
	}
}

/**
 * The main() method is the initial method that gets run, so it will start
 * the necessary goroutines to become a functional Raft node.
 */
func main() {
	Log = setUpLogger()
	Log.Info("GoRaft starting... Logger initialized.")

	config := loadConfig()

	Log.Info("Loaded config:", config)
	setLogLevel(config)

	// Check if state was loaded correctly from previous run.
	Log.Debug(state.GetNodeState(Log))

	runNode(config)
}
