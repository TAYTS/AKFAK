package broker

import (
	"AKFAK/cluster"
	"AKFAK/config"
	"AKFAK/proto/adminpb"
	"AKFAK/proto/clientpb"
	"AKFAK/proto/clustermetadatapb"
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
	ClusterMetadata     *cluster.Cluster
	adminServiceClient  map[int]adminpb.AdminServiceClient
	clientServiceClient map[int]clientpb.ClientServiceClient
	zkClient            zookeeperpb.ZookeeperServiceClient
	config              config.BrokerConfig
}

// InitNode create new broker node instance
func InitNode(config config.BrokerConfig) *Node {
	return &Node{
		ID:     config.ID,
		Host:   config.Host,
		Port:   config.Port,
		config: config,
	}
}

// InitAdminListener create a server listening connection
func (n *Node) InitAdminListener() {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", n.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)

	// bind the rpc server service
	adminpb.RegisterAdminServiceServer(server, n)
	clientpb.RegisterClientServiceServer(server, n)

	// setup the cluster metadata cache
	n.initClusterMetadataCache()

	// TODO [Fault tolerance]: check if the broker is not insync

	// start controller routine if the broker is select as the controller
	if int(n.ClusterMetadata.GetController().GetID()) == n.ID {
		go n.initControllerRoutine()
	}

	// setup ClientService peer connection
	go n.establishClientServicePeerConn()

	log.Printf("Broker %v start listening\n", n.ID)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

// InitControllerRoutine start the controller routine
func (n *Node) initControllerRoutine() {
	log.Printf("Broker %v start controller routine\n", n.ID)

	// setup the peer connection mapping
	if n.adminServiceClient == nil {
		n.adminServiceClient = make(map[int]adminpb.AdminServiceClient)
	}

	// Connect to all brokers
	n.updateAdminPeerConnection()

	// connect and store ZK rpc client
	n.zkClient = getZKClient(n.config.ZKConn)
}

// establishClientServicePeerConn start the ClientService peer connection
func (n *Node) establishClientServicePeerConn() {
	opts := grpc.WithInsecure()
	if n.clientServiceClient == nil {
		n.clientServiceClient = make(map[int]clientpb.ClientServiceClient)
	}

	// Connect to all brokers
	for _, brk := range n.ClusterMetadata.GetLiveBrokers() {
		peerID := int(brk.GetID())
		if peerID != n.ID {
			peerAddr := fmt.Sprintf("%v:%v", brk.GetHost(), brk.GetPort())
			ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			clientCon, err := grpc.DialContext(ctx, peerAddr, opts)
			if err != nil {
				log.Printf("Fail to connect to %v: %v\n", peerAddr, err)
				// TODO: Update the ZK about the fail node
				cancel()
				continue
			}
			clientServiceClient := clientpb.NewClientServiceClient(clientCon)
			n.clientServiceClient[peerID] = clientServiceClient
			cancel()
		}
	}

}

// InitClusterMetadataCache call the ZK to get the Cluster Metadata
func (n *Node) initClusterMetadataCache() {
	// connect to ZK
	zkClient := getZKClient(n.config.ZKConn)

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
	n.ClusterMetadata = cluster.InitCluster(res.GetClusterInfo())
}

func getZKClient(zkAddress string) zookeeperpb.ZookeeperServiceClient {
	// set up grpc dial
	opts := grpc.WithInsecure()

	zkCon, err := grpc.Dial(zkAddress, opts)
	if err != nil {
		log.Fatalf("Fail to connect to ZK: %v\n", err)
	}

	// return rpc client
	return zookeeperpb.NewZookeeperServiceClient(zkCon)
}
