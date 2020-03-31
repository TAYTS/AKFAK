package main

import (
	"AKFAK/broker"
	"os"
	"os/signal"
)

func main() {
	node1 := broker.InitNode(0, "0.0.0.0", 5001)
	node2 := broker.InitNode(1, "0.0.0.0", 5002)
	node3 := broker.InitNode(2, "0.0.0.0", 5003)
	nodes := []*broker.Node{node1, node2, node3}

	for _, node := range nodes {
		go node.InitAdminListener()
	}

	// Wait for Ctrl-C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
}
