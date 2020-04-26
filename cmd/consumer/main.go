package main

import (
	"AKFAK/consumer"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {
	// log setup
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// get the user input for initialise the Consumer
	contactServer := flag.String(
		"kafka-server",
		"",
		"Kafka server IP address and port number (e.g. 127.0.0.1:9092)")
	topicPtr := flag.String(
		"topic",
		"",
		"Topic to pull the message from")

	// print usage if not all fields for first round of qns provided
	if len(os.Args) < 4 {
		fmt.Println("usage: consumer -kafka-server <server_address:port> -topic <topic_name>")
		os.Exit(2)
	}
	flag.Parse()

	// initialise the Consumer
	log.Println("Initialising the Consumer...")
	c, availablePartitions := consumer.InitConsumer(*topicPtr, *contactServer)

	// construct the available partition index string
	availablePartitionsStr := "["
	for idx, partIdx := range availablePartitions {
		availablePartitionsStr += fmt.Sprint(partIdx)
		if idx != len(availablePartitions)-1 {
			availablePartitionsStr += ", "
		}
	}
	availablePartitionsStr += "]"

	// prompt the user to select the partition index to pull the message
	var partitionIdx int
	fmt.Printf("Which partition do you want to pull from?\nPartitions available: %v\n", availablePartitionsStr)

	fmt.Scanln(&partitionIdx) // get partition chosen

	// setup the partition to pull the message
	c.SetPartitionIdx(partitionIdx)
	log.Printf("Set up consumer to start pulling from topic %v partition %v\n", *topicPtr, partitionIdx)

	// Wait for Ctrl-C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// start consume message
	go c.Consume()

	<-ch

	// clear up resource
	c.CleanupResources()
}
