package config

// ZKConfig represent the zookeeper config JSON data
type ZKConfig struct {
	DataDir string `json:"dataDir"`
	Port    int    `json:"port"`
}
