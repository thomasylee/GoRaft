package global

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

type NodeHost struct {
	Url string `yaml:"url"`
	ApiPort uint32 `yaml:"api_port"`
	RpcPort uint32 `yaml:"rpc_port"`
}

type ConfigMap struct {
	LogLevel string `yaml:"log_level"`
	ElectionTimeout uint32 `yaml:"election_timeout"`
	ElectionTimeoutJitter uint32 `yaml:"election_timeout_jitter"`
	NodeId string `yaml:"node_id"`
	Nodes map[string]NodeHost `yaml:"node_hosts"`
}

var Config ConfigMap

/**
 * Load the config map from config.yaml.
 */
/*
func LoadConfig() map[interface{}]interface{} {
	configData, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		Log.Panic(err)
	}
	configYaml := string(configData)

	Config = make(map[interface{}]interface{})
	yaml.Unmarshal([]byte(configYaml), &Config)
	return Config
}
*/

func LoadConfig() ConfigMap {
	configData, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		Log.Panic(err)
	}

	err = yaml.Unmarshal(configData, &Config)
	if err != nil {
		Log.Panic(err)
	}

	return Config
}
