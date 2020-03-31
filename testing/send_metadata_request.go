package main

import (
	"AKFAK/proto/clientpb"
	"AKFAK/proto/metadatapb"
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

	req := &metadatapb.MetadataRequest{TopicName: "topic-1"}

	res, err := c.WaitOnMetadata(context.Background(), req)
	if err != nil {
		log.Fatalf("Error whil calling controller RPC: %v\n", err)
	}

	fmt.Println("Result:", res)
}
