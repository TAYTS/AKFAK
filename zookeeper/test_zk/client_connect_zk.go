package main

import (
	"AKFAK/proto/zkpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

func main() {
	opts := grpc.WithInsecure()
	broker := zkpb.Broker{
		Id:   2,
		Host: "127.0.0.1",
		Port: 3002,
	}
	cSock, err := grpc.Dial("127.0.0.1:3002", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	client := zkpb.NewZookeeperServiceClient(cSock)
	req := zkpb.ServiceDiscoveryRequest{
		Request: zkpb.ServiceDiscoveryRequest_BROKER,
		Query:   "/home/yijie/go/src/AKFAK/broker_config.json",
		Broker:  &broker,
	}
	stream, err := client.GetBrokers(context.Background(), &req)
	if err == nil {
		in, ok := stream.Recv()
		brokers := in.GetBrokerList()
		fmt.Println(brokers)
		fmt.Println(ok)
	} else {
		fmt.Println(err)
	}
	var input string
	fmt.Scanln(&input)
}
