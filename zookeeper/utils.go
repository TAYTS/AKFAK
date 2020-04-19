package zookeeper

import (
	"AKFAK/proto/adminclientpb"
	"AKFAK/proto/adminpb"
	"AKFAK/proto/clustermetadatapb"
	"AKFAK/utils"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

// LoadClusterStateFromFile parse the cluster state JSON and return in-memory cache of the cluster metadata
func LoadClusterStateFromFile(path string) clustermetadatapb.MetadataCluster {
	// parse the JSON byte into structs
	var clusterDataJSON clustermetadatapb.MetadataCluster
	err := utils.LoadJSONData(path, &clusterDataJSON)
	if err != nil {
		// load JSON data should not fail at ZK in our case
		panic(err)
	}

	// Set the Controller to be invalid
	clusterDataJSON.Controller = &clustermetadatapb.MetadataBroker{
		ID:   -1,
		Host: "",
		Port: -1,
	}

	// Clear the LiveBrokers
	clusterDataJSON.LiveBrokers = []*clustermetadatapb.MetadataBroker{}

	return clusterDataJSON
}

// WriteClusterStateToFile flush the cluster metadata to file
func WriteClusterStateToFile(path string, metadata clustermetadatapb.MetadataCluster) error {
	err := utils.FlushJSONData(path, metadata)
	// flush JSON data should not fail at ZK in our case
	if err != nil {
		panic(err)
	}

	return nil
}

// updateControllerMetadata used to update the controller when there is new LiveBroker
func (zk *Zookeeper) updateControllerMetadata() {
	// setup gRPC connection to controller
	ctrl := zk.clusterMetadata.GetController()
	ctrlConn, err := grpc.Dial(fmt.Sprintf("%v:%v", ctrl.GetHost(), ctrl.GetPort()), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Fail to connect to controller: %v\n", err)
		// TODO: select neew controller
	}

	// setup RPC service
	ctrlClient := adminpb.NewAdminServiceClient(ctrlConn)

	// send RPC call
	_, err = ctrlClient.UpdateMetadata(context.Background(), &adminclientpb.UpdateMetadataRequest{
		NewClusterInfo: zk.clusterMetadata.MetadataCluster,
	})
	if err != nil {
		log.Println("ZK failed to update controller")
		// TODO: select new controller
	}

	log.Println("ZK to controller cluster update successfull")
}
