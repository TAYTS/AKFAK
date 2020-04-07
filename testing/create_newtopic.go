package main

import (
	"AKFAK/proto/adminclientpb"
	"AKFAK/proto/adminpb"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

func main() {
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:5001", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := adminpb.NewAdminServiceClient(cc)

	req := &adminclientpb.AdminClientNewTopicRequest{
		Topic:             "topic1",
		NumPartitions:     3,
		ReplicationFactor: 3,
	}

	res, err := c.AdminClientNewTopic(context.Background(), req)
	if err != nil {
		log.Fatalf("Error whil calling controller RPC: %v\n", err)
	}

	fmt.Println("Result:", res.GetResponse())
}
