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

	// get the user input for initialise the Producer
	cgID := flag.Int(
		"cg-id",
		0,
		"Consumer group ID (e.g. 1)")
	contactServer := flag.String(
		"kafka-server",
		"",
		"Kafka server IP address and port number (e.g. 127.0.0.1:9092)")
	topicPtr := flag.String(
		"topics",
		"",
		"Topics to pull the message from, separated by comma (e.g. topic1,topic2")

	// print usage if user does not provide kafka server address and the topic
	if len(os.Args) < 4 {
		fmt.Println("usage: consumer -kafka-server <server_address:port> -topics <topic_names>")
		os.Exit(2)
	}
	flag.Parse()

	log.Println("Initialising the Consumer Group...")

	// initialise the Consumer Group
	cg := consumer.InitConsumerGroup(*cgID, *topicPtr, *contactServer)

	// Wait for Ctrl-C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// routine to pass user input to the Producer
	go func(cg *consumer.ConsumerGroup) {
		var inputStr string
		for {
			fmt.Println("Enter a topic to pull from")
			fmt.Scanln(&inputStr) // wait for something to be entered
			go cg.Consume(inputStr)
		}
	}(cg)

	<-ch
}