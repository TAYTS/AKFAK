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
	cID := flag.Int(
		"id",
		0,
		"Consumer ID (e.g. 1)")
	contactServer := flag.String(
		"kafka-server",
		"",
		"Kafka server IP address and port number (e.g. 127.0.0.1:9092)")
	partitionNum := flag.Int(
		"partition",
		0,
		"Partition Num (e.g. 1)"
	)
	topicPtr := flag.String(
		"topic",
		"",
		"Topic to pull the message from")


	// print usage if not all fields for first round of qns provided
	if len(os.Args) < 3 {
		fmt.Println("usage: consumer -id <consumer_id> -kafka-server <server_address:port> -topic <topic_name>")
		os.Exit(2)
	}
	flag.Parse()

	// initialise the Consumer
	log.Println("Initialising the Consumer...")
	c, numPartitions := consumer.InitConsumer(*cID, *topic, *contactServer)

	// choose partition
	var partition int
	fmt.Printf("Which partition do you want to pull from?\nPartitions available: %v\n", numPartitions)
	fmt.Scanln(&partition) // get partition chosen
	c.PartitionIdx = &partition

	log.Printf("Set up consumer to start pulling from topic %v partition %v\n", *topic, &partition)
	
	// Wait for Ctrl-C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	c.Consume()

	<-ch
}
