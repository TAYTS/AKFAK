package zookeeper

import (
	"AKFAK/proto/adminclientpb"
	"AKFAK/proto/adminpb"
	"AKFAK/proto/clustermetadatapb"
	"AKFAK/proto/consumermetadatapb"
	"AKFAK/utils"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

///////////////////////////////////
// 		     Public Methods		     //
///////////////////////////////////

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

// LoadConsumerStateFromFile parse the Consumer state JSON and return in-memory cache of the Consumer metadata
func LoadConsumerStateFromFile(path string) consumermetadatapb.MetadataConsumerState {
	// parse the JSON byte into structs
	var consumerDataJSON consumermetadatapb.MetadataConsumerState
	err := utils.LoadJSONData(path, &consumerDataJSON)
	if err != nil {
		// load JSON data should not fail at ZK in our case
		panic(err)
	}

	return consumerDataJSON
}

// WriteConsumerStateToFile flush the consumer metadata to file
func WriteConsumerStateToFile(path string, metadata consumermetadatapb.MetadataConsumerState) error {
	err := utils.FlushJSONData(path, metadata)
	if err != nil {
		panic(err)
	}
	return nil
}

///////////////////////////////////
// 	 	    Private Methods	    	 //
///////////////////////////////////

// updateControllerMetadata used to update the controller when there is new LiveBroker
func (zk *Zookeeper) updateControllerMetadata() {
	// wait for controller to be ready
	zk.waitCtrl.Wait()

	// setup gRPC connection to controller
	ctrl := zk.clusterMetadata.GetController()
	ctrlConn, err := grpc.Dial(fmt.Sprintf("%v:%v", ctrl.GetHost(), ctrl.GetPort()), grpc.WithInsecure())
	defer ctrlConn.Close()
	if err != nil {
		log.Fatalf("Fail to connect to controller: %v\n", err)
		return
	}

	// setup RPC service
	ctrlClient := adminpb.NewAdminServiceClient(ctrlConn)

	// send RPC call
	_, err = ctrlClient.UpdateMetadata(context.Background(), &adminclientpb.UpdateMetadataRequest{
		NewClusterInfo: zk.clusterMetadata.MetadataCluster,
	})
	if err != nil {
		log.Println("ZK failed to update controller")
		return
	}

	log.Println("ZK to controller cluster update successfull")
}

// handleControllerFailure is used to select the next available broker to be controller
// send ControllerElection rpc to the new controller
func (zk *Zookeeper) handleControllerFailure() {
	zk.mux.Lock()
	log.Println("ZK detect controller failure")

	// get the fail controller ID
	ctrlID := zk.clusterMetadata.GetController().GetID()

	// reset the controller
	zk.clusterMetadata.UpdateController(&clustermetadatapb.MetadataBroker{
		ID:   -1,
		Host: "",
		Port: -1,
	})

	for {
		// remove controller from live broker
		zk.clusterMetadata.MoveBrkToOfflineAndElectLeader(ctrlID)

		// select the first live broker to be controller
		availableLiveBrks := zk.clusterMetadata.GetLiveBrokers()

		if len(availableLiveBrks) > 0 {
			zk.clusterMetadata.UpdateController(availableLiveBrks[0])
			ctrlID = zk.clusterMetadata.GetController().GetID()

			zk.mux.Unlock()
			err := zk.sendControllerElection()
			if err == nil {
				break
			} else {
				zk.mux.Lock()
			}
		} else {
			zk.mux.Unlock()
			break
		}
	}
}

// sendControllerElection send ControllerElection RPC request to the selected broker
func (zk *Zookeeper) sendControllerElection() error {
	zk.mux.Lock()
	// setup gRPC connection to controller
	ctrl := zk.clusterMetadata.GetController()

	log.Printf("ZK elect Broker %v as new controller\n", ctrl.GetID())

	ctrlConn, err := grpc.Dial(fmt.Sprintf("%v:%v", ctrl.GetHost(), ctrl.GetPort()), grpc.WithInsecure())
	defer ctrlConn.Close()
	if err != nil {
		log.Fatalf("Fail to connect to controller: %v\n", err)
		return err
	}

	// setup RPC service
	ctrlClient := adminpb.NewAdminServiceClient(ctrlConn)

	// send RPC call
	_, err = ctrlClient.ControllerElection(context.Background(), &adminclientpb.ControllerElectionRequest{
		NewClusterInfo: zk.clusterMetadata.MetadataCluster,
	})
	if err != nil {
		log.Println("ZK failed to elect new controller", err)
		return err
	}
	log.Println("ZK successfully elect new controller")

	// set flag to wait to controller to be ready
	zk.waitCtrl.Add(1)

	zk.mux.Unlock()

	return nil
}
