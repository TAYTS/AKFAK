package main

import (
	"AKFAK/producer"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	brokerAddrs := map[int]string{
		0: "0.0.0.0:5001",
		1: "0.0.0.0:5002",
		2: "0.0.0.0:5003",
	}
	topicName := "topic1"

	p := producer.InitProducer(0, topicName, brokerAddrs)

	// Wait for Ctrl-C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

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
