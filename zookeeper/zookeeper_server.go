package zookeeper

import (
	"AKFAK/proto/zkpb"
	"AKFAK/utils"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

type Zookeeper struct {
	ID   int
	Host string
	Port int
	hasInit bool
	mux sync.Mutex
}

func (z *Zookeeper) GetBrokers(req *zkpb.ServiceDiscoveryRequest, server zkpb.ZookeeperService_GetBrokersServer) error {
	if !z.hasInit {
		switch req.Request {
		case zkpb.ServiceDiscoveryRequest_BROKER:
			brokerInfo := req.GetBroker()
			brokerList := z.LoadBrokerConfig(req.GetQuery())
			// Assign the first broker to connect to ZK be the coordinator
			// Temporarily saved in mem and not flush to the config file
			for i, v := range brokerList {
				if v.Id == brokerList[i].Id && v.GetHost() == brokerInfo.GetHost() && v.GetPort() == brokerInfo.GetPort(){
					v.IsCoordinator = true
					v.Id = int32(len(brokerList))
				} else {
					v.IsCoordinator = false
					v.Id = int32(i)
				}
			}
			res := zkpb.ServiceDiscoveryResponse{
				BrokerList: brokerList,
			}
			if err := server.Send(&res); err != nil {
				return err
			}
			z.hasInit = true
		default:
			return errors.New("Invalid request type")
		}
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
