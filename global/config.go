package global

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

var Config map[interface{}]interface{}

/**
 * Load the config map from config.yaml.
 */
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
