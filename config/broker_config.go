package config

// BrokerConfig represent the broker config JSON data
type BrokerConfig struct {
	ID     int    `json:"ID"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
	LogDir string `json:"logDir"`
	ZKConn string `json:"zookeeperConnection"`
}
