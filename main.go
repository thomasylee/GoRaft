package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/go-yaml/yaml"

	"github.com/thomasylee/GoRaft/api"
)

/**
 * Returns a timeout value between average - jitter and average + jitter.
 */
func calc_timeout(average int, jitter int) int {
	if jitter == 0 {
		return average
	}

	return average - jitter + rand.Intn(2 * jitter)
}

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

	// Randomize the election timeout to minimize the risk of two nodes
	// initiating an election at the same time.
	election_timeout := config["election_timeout"].(int)
	election_timeout_jitter := config["election_timeout_jitter"].(int)
	for {
		timeout := calc_timeout(election_timeout, election_timeout_jitter)
		select {
		case <-timeoutChannel:
			// Do nothing since we didn't time out.
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			// TODO: Start leader election process.
		}
	}
}
