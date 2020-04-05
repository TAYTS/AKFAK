package main

import (
	"AKFAK/proto/zkpb"
	"AKFAK/utils"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

type GrpcClient struct {
	conn *grpc.ClientConn
	client zkpb.ZookeeperServiceClient
}

func main() {
	opts := grpc.WithInsecure()
	conn, err := grpc.Dial("127.0.0.1:3000", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	client := zkpb.NewZookeeperServiceClient(conn)
	brokers := utils.GetBrokers("broker_config_test.json")
	var responseBrokers []*zkpb.Broker
	for _, v := range brokers {
		broker := zkpb.Broker{
			Id:                   v.Id,
			Host:                 v.Host,
			Port:                 v.Port,
			IsCoordinator:        false,
			XXX_NoUnkeyedLiteral: struct{}{},
			XXX_unrecognized:     nil,
			XXX_sizecache:        0,
		}
		responseBrokers = append(responseBrokers, &broker)
	}
	req := zkpb.ServiceDiscoveryRequest{
		BrokerList: responseBrokers,
	}
	res, err := client.GetBrokers(context.Background(), &req)
	if err == nil {
		fmt.Println(res)
	} else {
		fmt.Println(err)
	}
}