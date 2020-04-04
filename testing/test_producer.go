package main

import (
	"AKFAK/producer"
	"fmt"
)

func main() {
	brokerAddrs := make(map[int]string)
	brokerAddrs[0] = "0.0.0.0:5001"
	p := producer.InitProducer(0, brokerAddrs)
	records := []producer.Record{
		producer.Record{
			Topic:     "topic1",
			Timestamp: producer.GetCurrentTimeinMs(),
			Value:     []byte("message content"),
		},
	}

	for {
		fmt.Println("Press 'Enter' to send record")
		fmt.Scanln() // wait for something to be entered
		p.Send(records)
	}
}
