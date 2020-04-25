package config

// ZKConfig represent the zookeeper config JSON data
type ZKConfig struct {
	DataDir         string `json:"dataDir"`
	ConsumerDataDir string `json:"consumerDataDir"`
	Host            string `json:"host"`
	Port            int    `json:"port"`
}
