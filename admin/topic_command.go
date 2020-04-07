package admin

import (
	"flag"
	"fmt"
	"log"
)

type CommandInput struct {
	KafkaServer   string
	Topic         string
	Partitions    int
	ReplicaFactor int
	Operation     string
}

// ParseTopicCommandInput parse the command line inputs for managing the Topic
func ParseTopicCommandInput() CommandInput {
	kafkaSvrPtr := flag.String(
		"kafka-server",
		"",
		"Kafka server IP address and port number (e.g. 127.0.0.1:9092)")

	topicPtr := flag.String(
		"topic",
		"",
		"Topic to create, describe or delete")

	partitionPtr := flag.Int(
		"partitions",
		1,
		"Number of partitions for the topic. Default to 1.")

	replicaFactorPtr := flag.Int(
		"replica-factor",
		1,
		"The replication factor for each partition in the topic being created. Default to 1.")

	listPtr := flag.Bool(
		"list",
		false,
		"List all the topics.")

	describePtr := flag.Bool(
		"describe",
		false,
		"Describe a topic.")

	createPtr := flag.Bool(
		"create",
		false,
		"Create a new topic.")

	deletePtr := flag.Bool(
		"delete",
		false,
		"Delete a topic")

	flag.Parse()

	if *kafkaSvrPtr == "" {
		printErrorMessage("-kafka-server")
		log.Fatalf("Invalid argument for kafka server")
	}
	if *topicPtr == "" {
		printErrorMessage("-topic")
		log.Fatalf("Invalid argument for topic")
	}

	operation := ""
	switch {
	case *listPtr:
		operation = "list"
	case *describePtr:
		operation = "describe"
	case *createPtr:
		operation = "create"
	case *deletePtr:
		operation = "delete"
	}

	return CommandInput{
		KafkaServer:   *kafkaSvrPtr,
		Topic:         *topicPtr,
		Partitions:    *partitionPtr,
		ReplicaFactor: *replicaFactorPtr,
		Operation:     operation,
	}
}

func printErrorMessage(arg string) {
	fmt.Println("Invalid argument for", arg)
	fmt.Println("Options:")
	flag.PrintDefaults()
}
