package main

import (
	"AKFAK/producer"
	"flag"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	// get the user input for initialise the Producer
	flag.String(
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

	// TODO: Load all the brokers from config & rethink about this part
	// Maybe used the server that the user provide to initialise the producer instead
	brokerAddrs := map[int]string{
		0: "broker-1:5000",
		1: "broker-2:5000",
		2: "broker-3:5000",
	}

	fmt.Println("Initialise the Producer...")

	// initialise the Producer
	p := producer.InitProducer(0, *topicPtr, brokerAddrs)

	// Wait for Ctrl-C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

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
