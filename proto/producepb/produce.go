package producepb

import "AKFAK/proto/recordpb"

// InitEmptyProduceRequest initialise the ProduceRequest with default value and return its pointer type
func InitEmptyProduceRequest(topic string) *ProduceRequest {
	return &ProduceRequest{
		TopicName:    topic,
		RequiredAcks: 0,
		TopicData:    []*TopicData{},
	}
}

// InitProduceRequest initialise the ProduceRequest with provided data and return its pointer type
func InitProduceRequest(topic string, partition int, records ...*recordpb.Record) *ProduceRequest {
	prdReq := InitEmptyProduceRequest(topic)
	prdReq.AddRecord(partition, records...)

	return prdReq
}

// AddRecord append new record to the RecordBatch for a partition
func (prd *ProduceRequest) AddRecord(partition int, records ...*recordpb.Record) {
	// iterate to find the corresponding TopicData
	for _, topicData := range prd.GetTopicData() {
		if int(topicData.GetPartition()) == partition {
			topicData.RecordSet.AppendRecord(records...)
			return
		}
	}

	// create new TopicData if cannot find it
	prd.TopicData = append(prd.TopicData, &TopicData{
		Partition: int32(partition),
		RecordSet: recordpb.InitialiseRecordBatchWithData(records...),
	})
}
