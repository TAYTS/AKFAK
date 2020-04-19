package broker

import (
	"AKFAK/proto/adminpb"
	"AKFAK/proto/clientpb"
	"AKFAK/proto/clustermetadatapb"
	"AKFAK/proto/consumermetadatapb"
	"AKFAK/proto/zkmessagepb"
	"AKFAK/proto/zookeeperpb"
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

// Node represent the Kafka broker node
type Node struct {
	ID                  int
	Host                string
	Port                int
	Metadata            *clustermetadatapb.MetadataCluster
	ConsumerMetadata	*consumermetadatapb.MetadataConsumerState
	adminServiceClient  map[int]adminpb.AdminServiceClient
	clientServiceClient map[int]clientpb.ClientServiceClient
}

// TODO: Get the node list from metadata
var nodes = map[int]string{
	0: "0.0.0.0:5001",
	1: "0.0.0.0:5002",
	2: "0.0.0.0:5003",
}

// var nodes = map[int]string{
// 	0: "broker-0:5000",
// 	1: "broker-1:5000",
// 	2: "broker-2:5000",
// }

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

	// bind the rpc server service
	adminpb.RegisterAdminServiceServer(server, n)
	clientpb.RegisterClientServiceServer(server, n)

	// setup the cluster metadata cache
	n.InitClusterMetadataCache()

	// start controller routine if the broker is select as the controller
	if int(n.Metadata.GetController().GetID()) == n.ID {
		go n.InitControllerRoutine()
	}

	// setup ClientService peer connection
	go n.EstablishClientServicePeerConn()
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

// InitControllerRoutine start the controller routine
func (n *Node) InitControllerRoutine() {
	opts := grpc.WithInsecure()
	if n.adminServiceClient == nil {
		n.adminServiceClient = make(map[int]adminpb.AdminServiceClient)
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
			adminServiceClient := adminpb.NewAdminServiceClient(clientCon)
			n.adminServiceClient[nodeID] = adminServiceClient
		}
	}
}

// EstablishClientServicePeerConn start the ClientService peer connection
func (n *Node) EstablishClientServicePeerConn() {
	opts := grpc.WithInsecure()
	if n.clientServiceClient == nil {
		n.clientServiceClient = make(map[int]clientpb.ClientServiceClient)
	}

	// Connect to all brokers
	// TODO: Get the correct brokerID from metadata
	for nodeID, peerAddr := range nodes {
		if nodeID != n.ID {
			ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			clientCon, err := grpc.DialContext(ctx, peerAddr, opts)
			if err != nil {
				fmt.Printf("Fail to connect to %v: %v\n", peerAddr, err)
				// TODO: Update the ZK about the fail node
				cancel()
				continue
			}
			clientServiceClient := clientpb.NewClientServiceClient(clientCon)
			n.clientServiceClient[nodeID] = clientServiceClient
			cancel()
		}
	}
}

// InitClusterMetadataCache call the ZK to get the Cluster Metadata
func (n *Node) InitClusterMetadataCache() {
	// set up grpc dial
	opts := grpc.WithInsecure()
	zkAddress := "0.0.0.0:9092"
	zkCon, err := grpc.Dial(zkAddress, opts)
	if err != nil {
		log.Fatalf("Fail to connect to %v: %v\n", zkAddress, err)
	}

	// set rpc client
	zkClient := zookeeperpb.NewZookeeperServiceClient(zkCon)

	// create request with the current broker info
	req := &zkmessagepb.GetClusterMetadataRequest{
		Broker: &clustermetadatapb.MetadataBroker{
			ID:   int32(n.ID),
			Host: n.Host,
			Port: int32(n.Port),
		},
	}

	// send the GetClusterMetadata request
	res, err := zkClient.GetClusterMetadata(context.Background(), req)
	if err != nil {
		log.Fatalf("Fail to get the cluster metada from Zk")
	}

	// store the cluster metadata to cache
	n.Metadata = res.GetClusterInfo()
}
