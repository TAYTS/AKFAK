package main

import (
	"AKFAK/proto/clientpb"
	"AKFAK/proto/consumepb"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

func main() {
	// mock consumer
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("0.0.0.0:5000", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := clientpb.NewClientServiceClient(cc)

	req := &consumepb.GetAssignmentRequest{
		GroupID:	2,
		TopicName:	"topic1",
	}


	res, err := c.GetAssignment(context.Background(), req)
	if err != nil {
		log.Fatalf("Error whil calling controller RPC: %v\n", err)
	}

	fmt.Println("Result:", res)
}
