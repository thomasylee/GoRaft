package main

import (
	"io/ioutil"
	"log"

	"github.com/go-yaml/yaml"
)

/**
 * The main() method is the initial method that gets run, so it will start
 * the necessary goroutines to become a functional Raft node.
 *
 * Currently, all main() does is read config.yaml, parse the results into a
 * map, and print out the map.
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
}
