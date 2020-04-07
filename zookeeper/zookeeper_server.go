package zookeeper

import (
	"AKFAK/proto/zkpb"
	"AKFAK/utils"
	"errors"
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

func (z *Zookeeper) GetBrokers(req *zkpb.ServiceDiscoveryRequest, server zkpb.ZookeeperService_GetBrokersServer) error {
	if req.Request == zkpb.ServiceDiscoveryRequest_BROKER {
		brokerList := z.LoadBrokerConfig(req.GetQuery())
		res := zkpb.ServiceDiscoveryResponse{
			BrokerList: brokerList,
		}
		if err := server.Send(&res); err != nil {
			return err
		}
	} else {
		return errors.New("Invalid request type")
	}
	return nil
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
