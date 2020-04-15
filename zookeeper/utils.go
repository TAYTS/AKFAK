package zookeeper

import (
	"AKFAK/proto/clustermetadatapb"
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
