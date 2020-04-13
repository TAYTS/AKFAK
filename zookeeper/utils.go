package zookeeper

import (
	"AKFAK/proto/clustermetadatapb"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// LoadClusterStateFromFile parse the cluster state JSON and return in-memory cache of the cluster metadata
func LoadClusterStateFromFile(path string) clustermetadatapb.MetadataCluster {
	// load the JSON byte data
	clusterData, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}

	// parse the JSON byte into structs
	var clusterDataJSON clustermetadatapb.MetadataCluster
	if err := json.Unmarshal([]byte(clusterData), &clusterDataJSON); err != nil {
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
	// parse the data into JSON byte
	metadataBytes, marshallErr := json.MarshalIndent(metadata, "", " ")
	if marshallErr != nil {
		log.Println("Unable to convert data into JSON:", marshallErr)
		return marshallErr
	}

	// flush the JSON byte intp file
	writeErr := ioutil.WriteFile(path, metadataBytes, 0644)
	if writeErr != nil {
		log.Println("Unable to save JSON to file:", writeErr)
		return writeErr
	}

	return nil
}
