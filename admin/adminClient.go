package admin

import (
	"AKFAK/proto/adminpb"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

// Admin is the administrative client for Kafka which support managing and inspecting topics, brokers
type Admin struct{}

// CreateAdminClient create new KafkaAdminClient instance
func CreateNew() {
	// Parse the command input with sub command: topic, partition, replica
	command := ParseTopicCommandInput()
	fmt.Println("Command", command)

	// TODO: Initialise the gRPC connection to the Kafka server

	// TODO: Fetch the metadata for cluster

	// Initialise the gRPC connection with the Kafka controller
	fmt.Println("Connecting to one of the Kafka server")

	// TODO: Get the controller host and port
	gCon, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect to Kafka server: %v\n", err)
	}

	adminCon := adminpb.NewAdminServiceClient(gCon)

	req := &adminpb.AdminClientNewTopicRequest{
		Topic:             command.Topic,
		NumPartitions:     int32(command.Partitions),
		ReplicationFactor: int32(command.ReplicaFactor),
	}

	// Send the topic operation request
	res, err := adminCon.AdminClientNewTopic(context.Background(), req)
	if err != nil {
		log.Fatalf("Fail to sending the operation request to Kafka controller")
	}
	log.Printf("Response from Kafka controller: %v\n", res.GetResponse())
}
