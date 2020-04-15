package main

import (
	consumer "AKFAK/consumer"
	fmt "fmt"
)

func main() {

	NUM_CONSUMERGROUP := 2
	for i := 1; i <= NUM_CONSUMERGROUP; i++ {
		consumer.InitGroupConsumer(i)
	}
	fmt.Printf("Total Consumer Group created: [%d] \n", NUM_CONSUMERGROUP)

	// Assignment(topicName, partitionIdx, brokerIdx)
	c := consumer.Consumer{}
	c.Assignment("hello", 1, 1)
	c.Assignment("bye", 2, 2)

}
