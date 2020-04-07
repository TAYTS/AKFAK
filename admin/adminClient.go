package admin

import (
	"AKFAK/proto/adminpb"
	"AKFAK/proto/adminclientpb"
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
	opts := grpc.WithInsecure()
	conn, err := grpc.Dial(command.KafkaServer, opts)
	if err != nil {
		log.Fatalf("Unable to connect to Kafka server %s: %v", command.KafkaServer, err)
	}
	adminClient := adminpb.NewAdminServiceClient(conn)
	log.Println(adminClient)
	
	// TODO: Fetch the metadata for cluster
	metadata, err := adminClient.GetMetadata(context.Background(), &adminclientpb.GetMetadataRequest{})
	if err != nil {
		log.Fatal("Failed to send getmetadata request to Kafka server")
	}

	// Initialise the gRPC connection with the Kafka controller
	fmt.Printf("Connecting to controller: %d\n", metadata.GetControllerID())

	// TODO: Get the controller host and port based on metadata.GetControllerID()
	// gCon, err := grpc.Dial(controller, opts)
	gCon, err := grpc.Dial("localhost:5001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect to Kafka server: %v\n", err)
	}

	adminCon := adminpb.NewAdminServiceClient(gCon)

	req := &adminclientpb.AdminClientNewTopicRequest{
		Topic:             command.Topic,
		NumPartitions:     int32(command.Partitions),
		ReplicationFactor: int32(command.ReplicaFactor),
	}

	// Send the topic operation request
	res, err := adminCon.AdminClientNewTopic(context.Background(), req)
	if err != nil {
		log.Fatalf("Fail to send operation request to Kafka controller: %v", err)
	}
	log.Printf("Response from Kafka controller: %v\n", res.GetResponse())
}
