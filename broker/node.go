package broker

import (
	"AKFAK/proto/adminpb"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

// Node represent the Kafka broker node
type Node struct {
	ID       int
	Host     string
	Port     int
	Metadata []int // Get the cluster metadata
	peerCon  map[int]adminpb.AdminServiceClient
}

// InitNode create new broker node instance
func InitNode(ID int, host string, port int) *Node {
	return &Node{
		ID:   ID,
		Host: host,
		Port: port,
	}
}

// InitAdminListener create a server listening connection
func (n *Node) InitAdminListener() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", n.Host, n.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)

	adminpb.RegisterAdminServiceServer(server, n)

	// TODO: call ZK to get the update metadata
	// TODO: initialise the internal state cache

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

// InitControllerRoutine start the controller routine
func (n *Node) InitControllerRoutine() {
	// TODO: Get the node list from metadata
	nodes := map[int]string{
		0: "0.0.0.0:5001",
		1: "0.0.0.0:5002",
		2: "0.0.0.0:5003",
	}

	opts := grpc.WithInsecure()
	if n.peerCon == nil {
		n.peerCon = make(map[int]adminpb.AdminServiceClient)
	}

	// Connect to all brokers
	// TODO: Get the correct brokerID from metadata
	for nodeID, peerAddr := range nodes {
		if nodeID != n.ID { // TODO: Need to update this after finalise the metadata
			clientCon, err := grpc.Dial(peerAddr, opts)
			if err != nil {
				fmt.Printf("Fail to connect to %v: %v\n", peerAddr, err)
				// TODO: Update the ZK about the fail node
				continue
			}
			serviceClient := adminpb.NewAdminServiceClient(clientCon)
			n.peerCon[nodeID] = serviceClient
		}
	}
}
