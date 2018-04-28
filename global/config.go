package global

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

// NodeHost represents a host node in the Raft cluster.
type NodeHost struct {
	Url     string `yaml:"url"`
	ApiPort uint32 `yaml:"api_port"`
	RpcPort uint32 `yaml:"rpc_port"`
}

// ConfigMap contains all the configurations loaded from the config file.
type ConfigMap struct {
	LogLevel              string              `yaml:"log_level"`
	ElectionTimeout       uint32              `yaml:"election_timeout"`
	ElectionTimeoutJitter uint32              `yaml:"election_timeout_jitter"`
	LeaderHeartbeatPeriod uint32              `yaml:"leader_heartbeat_period"`
	NodeId                string              `yaml:"node_id"`
	Nodes                 map[string]NodeHost `yaml:"node_hosts"`
}

// ConfigMap contains the loaded configurations.
var Config ConfigMap

// LoadConfig loads the config map from config.yaml into Config.
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
