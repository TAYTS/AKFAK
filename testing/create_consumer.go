package main

import (
	fmt "fmt"
	consumer "AKFAK/consumer" 
	"math/rand"
	"strconv"
)


func main() {

	num_con := 5
	for i := 1; i <= num_con; i++ {
		min := 0
		max := 5
		
		topic := rand.Intn(max - min)
		topicStr := strconv.Itoa(topic)
		partition := rand.Intn(max - min)
		broker := rand.Intn(max - min)
		consumer.InitConsumer(topicStr, partition, broker)
		fmt.Printf("Consumer %d initiated\n", i)
		fmt.Printf("===" + topicStr + "===\n" + "P%d \nB%d \n", partition, broker)

	}
	fmt.Printf("====================================== \n Total Consumers = %d \n", num_con)
}
