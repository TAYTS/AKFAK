package main

import (
	topic "AKFAK/admin/topicutils"
	"AKFAK/proto/adminclientpb"
	"AKFAK/proto/adminpb"
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
)

func main() {
	// parse the user input opration
	topicCommandInput := topic.ParseTopicCommandInput()

	// get the controller information
	conn, err := grpc.Dial(topicCommandInput.KafkaServer, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect to Kafka server %s: %v", topicCommandInput.KafkaServer, err)
	}
	adminClient := adminpb.NewAdminServiceClient(conn)

	resp, err := adminClient.GetController(context.Background(), &adminclientpb.GetControllerRequest{})
	if err != nil {
		log.Fatal("Failed to send getmetadata request to Kafka server")
	}

	controllerHost := resp.GetHost()
	controllerPort := resp.GetPort()
	ctrlConnectionStr := fmt.Sprintf("%v:%v", controllerHost, controllerPort)

	switch topicCommandInput.Operation {
	case string(topic.CREATE_TOPIC):
		// create new topic routine
		topic.CreateNewTopic(ctrlConnectionStr, topicCommandInput)
	case string(topic.LIST_TOPIC):
		// list all topic routine
	case string(topic.DELETE_TOPIC):
		// delete topic routine
	default:
		os.Exit(2)
	}
}
