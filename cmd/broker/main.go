package main

import (
	"AKFAK/broker"
	"AKFAK/config"
	"AKFAK/utils"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// log setup
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// get the user input for initialise the Broker
	configPath := flag.String(
		"server-config",
		"",
		"Config to setup the Kafka broker")

	// print usage if user does not provide sufficient infomation to start the broker
	if len(os.Args) < 2 {
		fmt.Println("usage: broker -server-config <server config filepath>")
		os.Exit(2)
	}
	flag.Parse()

	// parse the JSON byte into structs
	var brkConfigJSON config.BrokerConfig
	err := utils.LoadJSONData(*configPath, &brkConfigJSON)
	if err != nil {
		panic(err)
	}

	// start the broker
	node := broker.InitNode(brkConfigJSON)
	node.InitAdminListener()
}
