package topicutils

import (
	"AKFAK/proto/adminclientpb"
	"AKFAK/proto/adminpb"
	"context"
	"log"

	"google.golang.org/grpc"
)

// CreateNewTopic create new KafkaAdminClient instance
func CreateNewTopic(ctrlConnectionStr string, userInput CommandInput) {
	gCon, err := grpc.Dial(ctrlConnectionStr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect to controller: %v\n", err)
	}
	adminCon := adminpb.NewAdminServiceClient(gCon)

	req := &adminclientpb.AdminClientNewTopicRequest{
		Topic:             userInput.Topic,
		NumPartitions:     int32(userInput.Partitions),
		ReplicationFactor: int32(userInput.ReplicaFactor),
	}

	// Send the topic operation request
	res, err := adminCon.AdminClientNewTopic(context.Background(), req)
	if err != nil {
		log.Fatalf("Fail to send operation request to Kafka controller: %v", err)
	}
	log.Printf("Response from Kafka controller: %v\n", res.GetResponse().GetMessage())
}
