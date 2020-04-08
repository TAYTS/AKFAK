package main

import (
	"AKFAK/broker"
	"flag"
	"fmt"
	"os"
)

func main() {
	// get the user input for initialise the Broker
	// TODO: Read the broker from config file from the given path
	// flag.String(
	// 	"-server-config",
	// 	"",
	// 	"Config to setup the Kafka broker")
	brkIDPtr := flag.Int(
		"brokerID",
		0,
		"Broker unique ID")
	hostPtr := flag.String(
		"host",
		"0.0.0.0",
		"Address for the broker to listen for connection")
	portPtr := flag.Int(
		"port",
		5000,
		"Port for the broker to listen for connection")

	// print usage if user does not provide sufficient infomation to start the broker
	if len(os.Args) < 6 {
		fmt.Println("usage: broker -brokerID <broker_ID> -host <ip_address> -port <port_number>")
		os.Exit(2)
	}
	flag.Parse()

	// start the broker
	node := broker.InitNode(*brkIDPtr, *hostPtr, *portPtr)
	node.InitAdminListener()
}
