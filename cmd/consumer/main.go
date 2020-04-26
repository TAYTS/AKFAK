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
		"Partition Num (eg. 1)")
	topicPtr := flag.String(
		"topic",
		"",
		"Topic to pull the message from")

	// print usage if not all fields provided
	if len(os.Args) < 4 {
		fmt.Println("usage: consumer -id <consumer_id> -kafka-server <server_address:port> -topic <topic_name> -part <partition>")
		os.Exit(2)
	}
	flag.Parse()

	log.Println("Initialising the Consumer...")

	// initialise the Consumer
	c := consumer.InitConsumer(*cID, *topicPtr, *partitionNum, *contactServer)

	// Wait for Ctrl-C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// routine to pass user input to the Consumer
	go func(c *consumer.Consumer) {
		var inputStr string
		for {
			fmt.Println("Press enter to pull message")
			fmt.Scanln(&inputStr) // wait for something to be entered
			go c.Consume(inputStr)
		}
	}(cg)
	<-ch
}
