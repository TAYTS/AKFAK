package broker

import (
	"AKFAK/broker/partition"
	"AKFAK/proto/clientpb"
	"AKFAK/proto/recordpb"
	"fmt"
	"errors"
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

// ReadRecordBatchFromLocal is a helper function for Consume request handler to read the Record from local log file
func ReadRecordBatchFromLocal(topicName string, partitionID int, fileHandlerMapping map[int]*recordpb.FileRecord) (*recordpb.RecordBatch, error) {
	fHandler, exist := fileHandlerMapping[partitionID]
	if exist {
		recordBatch, err := fHandler.ReadNextRecordBatch()
		if err != nil {
			return nil, err
		}
		return recordBatch, nil
	}
	return nil, errors.New("No such file")
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