package main

import (
	"AKFAK/zookeeper"
	"fmt"
)

func main() {
	data := zookeeper.LoadClusterStateFromFile("cluster_state.json")

	fmt.Println("Loaded the config data:", data)

	zookeeper.WriteClusterStateToFile("hello.json", data)
}
