package main

import (
	"AKFAK/proto/adminclientpb"
	"AKFAK/proto/adminpb"
	"context"
	"flag"
	"log"

	"google.golang.org/grpc"
)

func main() {
	// TODO: Remove this after the ZK is done!!!!
	brkPtr := flag.String(
		"-broker",
		"broker-0:5000",
		"Broker connection string")
	flag.Parse()

	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(*brkPtr, opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := adminpb.NewAdminServiceClient(cc)

	req := &adminclientpb.ControllerElectionRequest{
		BrokerID: 0,
		HostName: "broker-0",
		Port:     5000,
	}

	_, err = c.ControllerElection(context.Background(), req)
	if err != nil {
		log.Fatalf("Error whil calling controller RPC: %v\n", err)
	}
}
