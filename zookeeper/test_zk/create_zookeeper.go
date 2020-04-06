package main

import (
	"AKFAK/zookeeper"
	"fmt"
)

func main() {
	zk := zookeeper.Zookeeper{
		ID:   1,
		Host: "0.0.0.0",
		Port: 3000,
	}

	brokers := zk.LoadBrokerConfig("broker_config_test.json")
	for _, broker := range brokers {
		go zk.InitBrokerListener(*broker)
	}
	var input string
	fmt.Scanln(&input)
}