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

func (z *Zookeeper) GetBrokers(ctx context.Context, req *zkpb.ServiceDiscoveryRequest) (*zkpb.ServiceDiscoveryResponse, error) {
	if req.Request == zkpb.ServiceDiscoveryRequest_BROKER {
		brokerList := z.LoadBrokerConfig(req.GetQuery())
		res := zkpb.ServiceDiscoveryResponse{
			BrokerList: brokerList,
		}
		return &res, nil
	} else {
		return &zkpb.ServiceDiscoveryResponse{}, nil
	}
}

func (z *Zookeeper) LoadBrokerConfig(configFilePath string) []*zkpb.Broker {
	brokers := utils.GetBrokers(configFilePath)
	return brokers
}

func (z *Zookeeper) InitBrokerListener(broker *zkpb.Broker) {
	fmt.Printf("Listening to Host: %v on port %v\n", broker.Host, broker.Port)

	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", broker.Host, broker.Port))
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
