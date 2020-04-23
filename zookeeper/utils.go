package zookeeper

import (
	"AKFAK/proto/clustermetadatapb"
	"AKFAK/proto/consumermetadatapb"
	"AKFAK/utils"
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
