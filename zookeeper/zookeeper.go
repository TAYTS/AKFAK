package zookeeper

import (
	"AKFAK/cluster"
	"AKFAK/config"
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
	clusterMetadata *cluster.Cluster
	mux             sync.Mutex
	config          config.ZKConfig
}

// InitZookeeper create the Zookeeper instance and load the cluster state data
func InitZookeeper(config config.ZKConfig) *Zookeeper {
	// load the cluster state file into in-memory data
	clusterMetadata := LoadClusterStateFromFile(config.DataDir)

	return &Zookeeper{
		Host:            config.Host,
		Port:            config.Port,
		clusterMetadata: cluster.InitCluster(&clusterMetadata),
		config:          config,
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
