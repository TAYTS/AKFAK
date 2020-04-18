package zookeeper

import (
	"AKFAK/proto/clustermetadatapb"
	"AKFAK/proto/zookeeperpb"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

// Zookeeper represent the Zookeeper server state
type Zookeeper struct {
	Host            string
	Port            int
	clusterMetadata clustermetadatapb.MetadataCluster
	mux             sync.Mutex
}

// InitZookeeper create the Zookeeper instance and load the cluster state data
func InitZookeeper() *Zookeeper {
	// load the cluster state file into in-memory data
	clusterMetadata := LoadClusterStateFromFile("cluster_state.json")

	return &Zookeeper{
		Host:            "0.0.0.0",
		Port:            9092,
		clusterMetadata: clusterMetadata,
	}
}

// InitZKListener create Zookeeper server listener
func (zk *Zookeeper) InitZKListener() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", zk.Host, zk.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)

	zookeeperpb.RegisterZookeeperServiceServer(server, zk)

	// start the server
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}
