package config

// ZKConfig represent the zookeeper config JSON data
type ZKConfig struct {
	DataDir string `json:"dataDir"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
}
