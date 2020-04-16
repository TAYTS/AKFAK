package broker

import (
	"AKFAK/broker/partition"
	"AKFAK/proto/adminclientpb"
	"AKFAK/proto/clientpb"
	"AKFAK/proto/recordpb"
	"errors"
	"fmt"
)

// WriteRecordBatchToLocal is a helper function for Produce request handler to save the RecordBatch to the local log file
func WriteRecordBatchToLocal(topicName string, partitionID int, fileHandlerMapping map[int]*recordpb.FileRecord, recordBatch *recordpb.RecordBatch) {
	fHandler, exist := fileHandlerMapping[partitionID]
	if exist {
		fHandler.WriteToFile(recordBatch)
	} else {
		filePath := fmt.Sprintf("./%v/%v", partition.ConstructPartitionDirName(topicName, partitionID), partition.ContructPartitionLogName(topicName))
		fileRecordHandler, _ := recordpb.InitialiseFileRecordFromFile(filePath)
		fileHandlerMapping[partitionID] = fileRecordHandler
		fileRecordHandler.WriteToFile(recordBatch)
	}
}

// CleanupProducerResource help to clean up the Producer resources
func CleanupProducerResource(replicaConn map[int]clientpb.ClientService_ProduceClient, fileHandlerMapping map[int]*recordpb.FileRecord) {
	for _, rCon := range replicaConn {
		err := rCon.CloseSend()
		if err != nil {
			fmt.Printf("Closing connection error: %v\n", err)
		}
	}

	for _, fileHandler := range fileHandlerMapping {
		err := fileHandler.CloseFile()
		if err != nil {
			fmt.Printf("Closing file error: %v\n", err)
		}
	}
}

// generateNewPartitionRequestData create a mapping of the brokerID and the create new partition RPC request
func (n *Node) generateNewPartitionRequestData(topicName string, numPartitions int, replicaFactor int) (map[int]*adminclientpb.AdminClientNewPartitionRequest, error) {
	// check if the topic exist
	topicExist := false
	if n.ClusterMetadata.GetPartitionsByTopic(topicName) != nil {
		topicExist = true
	}

	numBrokers := len(n.ClusterMetadata.GetBrokers())

	if !topicExist {
		newTopicPartitionRequests := make(map[int]*adminclientpb.AdminClientNewPartitionRequest)

		for partID := 0; partID < numPartitions; partID++ {
			// distribute the partitions among the brokers
			brokerID := int(n.ClusterMetadata.GetBrokers()[partID%numBrokers].GetID())

			for replicaIdx := 0; replicaIdx < replicaFactor; replicaIdx++ {
				// distribute the replicas among the brokers
				replicaBrokerID := (brokerID + replicaIdx + partID) % numBrokers

				request, exist := newTopicPartitionRequests[replicaBrokerID]
				if exist {
					request.PartitionID = append(request.PartitionID, int32(partID))
				} else {
					request := &adminclientpb.AdminClientNewPartitionRequest{
						Topic:       topicName,
						PartitionID: []int32{int32(partID)},
						ReplicaID:   int32(replicaBrokerID),
					}
					newTopicPartitionRequests[replicaBrokerID] = request
				}
			}
		}

		return newTopicPartitionRequests, nil
	}
	return nil, errors.New("Topic already exist")
}

// createLocalPartitionFromReq take the create new partition RPC request and create the local partition directory and log file
func (n *Node) createLocalPartitionFromReq(req *adminclientpb.AdminClientNewPartitionRequest) error {
	topicName := req.GetTopic()
	partitionID := req.GetPartitionID()

	fmt.Printf("Node %v: Create partition %v\n", n.ID, partitionID)

	rootPath := n.config.LogDir
	for _, partID := range partitionID {
		err := partition.CreatePartitionDir(rootPath, topicName, int(partID))
		if err != nil {
			return err
		}
	}
	return nil
}
