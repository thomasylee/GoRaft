package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/go-yaml/yaml"

	"github.com/thomasylee/GoRaft/api"
)

/**
 * The main() method is the initial method that gets run, so it will start
 * the necessary goroutines to become a functional Raft node.
 */
func main() {
	config_data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	config_yaml := string(config_data)

	config := make(map[interface{}]interface{})

	yaml.Unmarshal([]byte(config_yaml), &config)
	log.Println(config)

	timeoutChannel := make(chan bool)
	go api.RunServer(timeoutChannel, config["api_port"].(int))

	timeout := time.Duration(config["heartbeat_timeout"].(int))
	for {
		select {
		case <-timeoutChannel:
			// Do nothing since we didn't time out.
		case <-time.After(timeout * time.Millisecond):
			// TODO: Start leader election process.
		}
	}
}
