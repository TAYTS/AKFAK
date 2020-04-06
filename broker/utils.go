package broker

import (
	"AKFAK/broker/partition"
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
