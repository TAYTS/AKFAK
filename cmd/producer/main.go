package main

import (
	"AKFAK/producer"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {
	// get the user input for initialise the Producer
	contactServer := flag.String(
		"kafka-server",
		"",
		"Kafka server IP address and port number (e.g. 127.0.0.1:9092)")
	topicPtr := flag.String(
		"topic",
		"",
		"Topic to send the message")

	// print usage if user does not provide kafka server address and the topic
	if len(os.Args) < 4 {
		fmt.Println("usage: producer -kafka-server <server_address:port> -topic <topic_name>")
		os.Exit(2)
	}
	flag.Parse()

	log.Println("Initialise the Producer...")

	// initialise the Producer
	p := producer.InitProducer(0, *topicPtr, *contactServer)

	// Wait for Ctrl-C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	log.Println("Done initialise the Producer")

	// routine to pass user input to the Producer
	go func(p *producer.Producer) {
		var inputStr string
		for {
			fmt.Println("Press 'Enter' to send record")
			fmt.Scanln(&inputStr) // wait for something to be entered
			go p.Send(inputStr)
		}
	}(p)

	<-ch
}
