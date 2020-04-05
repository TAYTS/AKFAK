package main

import (
	"AKFAK/proto/metadatapb"
	"google.golang.org/grpc"
	"log"
	"net"
)
type metadataServer struct{}

// client = brokers
// server = zookeeper

//https://zookeeper.apache.org/doc/r3.3.3/zookeeperAdmin.html#sc_configuration// tickTime=2000
// the location where ZooKeeper will store the in-memory database snapshots and, unless specified otherwise, the transaction log of updates to the database.
// dataDir=/var/zookeeper/

// Start up - Send metadata response to the brokers so they can update their state
// When the gRPC connection of a controller is down, assign another broker to be the controller
// Define message specification
// Before Producer's first message, run go bin to create topics, partitions, replicas


// /
// dataDir
// clientPort=2181
// Amount of time, in ticks (see tickTime), to allow followers to connect and sync to a leader. Increased this value as needed, if the amount of data managed by ZooKeeper is large.
// initLimit=5

// Amount of time, in ticks (see tickTime), to allow followers to sync with ZooKeeper. If followers fall too far behind a leader, they will be dropped.
// syncLimit=2

// server.1=zoo1:2888:3888
// server.2=zoo2:2888:3888
// server.3=zoo3:2888:3888

//create : creates a node at a location in the tree
//delete : deletes a node
//exists : tests if a node exists at a location
//get data : reads the data from a node
//set data : writes data to a node
//get children : retrieves a list of children of a node
//sync : waits for data to be propagated


func main() {
	listen, err := net.Listen("tcp", "0.0.0.0:2182")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		return
	}
	gRPCmetadataServer := grpc.NewServer()
	admin
	metadatapb.RegisterMetadataServiceServer(gRPCmetadataServer, &metadataServer{})

	if err := gRPCmetadataServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}