package main

import (
	consumer "AKFAK/consumer"
)

func main() {
	consumer.InitConsumerGroup(1, "fish,time", "0.0.0.0:5001")
	// Assignment(topicName, partitionIdx, brokerIdx)
	//c := consumer.Consumer{}
	//c.Assignment("hello", 1, 1)
	//c.Assignment("bye", 2, 2)

}
