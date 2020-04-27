package topicutils

import (
	"flag"
	"fmt"
	"os"
)

// CommandInput consists of all the data required for the
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

	// CREATE COMMAND //
	createCommand := flag.NewFlagSet(string(CREATE_TOPIC), flag.ExitOnError)
	createTopicPtr := createCommand.String(
		"topic",
		"",
		"Topic to create")
	partitionPtr := createCommand.Int(
		"partitions",
		0,
		"Number of partitions for the topic")
	replicaFactorPtr := createCommand.Int(
		"replica-factor",
		0,
		"The replication factor for each partition in the topic being created")
	// CREATE COMMAND //

	// LIST COMMAND //
	// flag.NewFlagSet(string(LIST_TOPIC), flag.ExitOnError)
	// LIST COMMAND //

	// DELETE COMMAND //
	// deleteCommand := flag.NewFlagSet(string(DELETE_TOPIC), flag.ExitOnError)
	// deleteTopicPtr := deleteCommand.String(
	// 	"topic",
	// 	"",
	// 	"Topic to delete")
	// DELETE COMMAND //

	// print usage if user does not provide kafka server address and the entry command
	if len(os.Args) <= 3 {
		fmt.Println("usage: admin-topic -kafka-server <server_address:port> <command> [<args>]")
		fmt.Println("Command:")
		fmt.Println("create       : Create new topic")
		createCommand.PrintDefaults()
		// fmt.Println("list         : List all topic")
		// fmt.Println("delete       : Delete a topic")
		// deleteCommand.PrintDefaults()
		os.Exit(2)
	}

	// parse and check kafka server address
	flag.Parse()
	if *kafkaSvrPtr == "" {
		printErrorMessage("-kafka-server")
		os.Exit(2)
	}

	// parse command and check parameter
	operation := ""
	topic := ""
	switch os.Args[3] {
	case string(CREATE_TOPIC):
		createCommand.Parse(os.Args[4:])
		operation = string(CREATE_TOPIC)
		if *createTopicPtr == "" {
			printErrorMessage("-topic")
			createCommand.PrintDefaults()
			os.Exit(2)
		} else if *partitionPtr == 0 {
			printErrorMessage("-partitions")
			createCommand.PrintDefaults()
			os.Exit(2)
		} else if *replicaFactorPtr == 0 {
			printErrorMessage("-replica-factor")
			createCommand.PrintDefaults()
			os.Exit(2)
		}
		topic = *createTopicPtr
	// case string(LIST_TOPIC):
	// 	operation = string(LIST_TOPIC)
	// case string(DELETE_TOPIC):
	// 	deleteCommand.Parse(os.Args[4:])
	// 	operation = string(DELETE_TOPIC)
	// 	if *deleteTopicPtr == "" {
	// 		printErrorMessage("-topic")
	// 		deleteCommand.PrintDefaults()
	// 		os.Exit(2)
	// 	}
	// 	topic = *deleteTopicPtr
	default:
		fmt.Printf("%q is not a valid command.\n", os.Args[3])
		os.Exit(2)
	}

	return CommandInput{
		KafkaServer:   *kafkaSvrPtr,
		Topic:         topic,
		Partitions:    *partitionPtr,
		ReplicaFactor: *replicaFactorPtr,
		Operation:     operation,
	}
}

func printErrorMessage(arg string) {
	fmt.Printf(`Missing value for argument "%v"`, arg)
}
