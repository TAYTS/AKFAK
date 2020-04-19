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
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:5001", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := clientpb.NewClientServiceClient(cc)

	req := &consumepb.GetAssignmentRequest{
		GroupID:	c.GroupID,
		TopicName:	c.TopicName,
	}


	res, err := c.GetAssignments(context.Background(), req)
	if err != nil {
		log.Fatalf("Error whil calling controller RPC: %v\n", err)
	}

	fmt.Println("Result:", res.GetResponse())
}
