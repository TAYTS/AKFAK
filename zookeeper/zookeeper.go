package zookeeper

import (
	"AKFAK/proto/zkpb"
	"AKFAK/utils"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Zookeeper struct {
	ID   int
	Host string
	Port int
}

func (z *Zookeeper) GetBrokers(context.Context, *zkpb.ServiceDiscoveryRequest) (*zkpb.ServiceDiscoveryResponse, error) {
	panic("implement me")
}

func (z *Zookeeper) LoadBrokerConfig(configFilePath string) []utils.Broker {
	brokers := utils.GetBrokers(configFilePath)
	return brokers
}

func (z *Zookeeper) InitBrokerListener(broker utils.Broker) {
	fmt.Printf("Listening to Host: %v on port %v\n", broker.Host, broker.Port)

	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", broker.Host, broker.Port))
	//listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:4000"))

	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)
	zkpb.RegisterZookeeperServiceServer(server, z)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
	server.Serve(listener)
}
