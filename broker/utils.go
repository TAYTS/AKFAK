package broker

import (
	"AKFAK/broker/partition"
	"AKFAK/proto/clientpb"
	"AKFAK/proto/recordpb"
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
