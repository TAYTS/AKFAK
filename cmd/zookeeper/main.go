package main

import (
	"AKFAK/config"
	"AKFAK/utils"
	"AKFAK/zookeeper"
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
		"Config to setup the Zookeeper")

	// print usage if user does not provide sufficient infomation to start the broker
	if len(os.Args) < 2 {
		fmt.Println("usage: broker -server-config <server config filepath>")
		os.Exit(2)
	}
	flag.Parse()

	// parse the JSON byte into structs
	var zkConfigJSON config.ZKConfig
	err := utils.LoadJSONData(*configPath, &zkConfigJSON)
	if err != nil {
		panic(err)
	}

	// start the ZK
	zk := zookeeper.InitZookeeper(brkConfigJSON)
	zk.InitZKListener()
}
