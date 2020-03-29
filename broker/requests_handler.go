package broker

import (
	"AKFAK/broker/partition"
	"AKFAK/proto/adminpb"
	"context"
)

// AdminClientNewTopic create new topic
func (n *Node) AdminClientNewTopic(ctx context.Context, req *adminpb.AdminClientNewTopicRequest) (*adminpb.AdminClientNewTopicResponse, error) {

	// get request data
	topicName := req.GetTopic()
	numPartitions := int(req.GetNumPartitions())
	replicaFactor := int(req.GetReplicationFactor())

	// handling request
	newPartitionReqMap, err := n.newPartitionRequestData(topicName, numPartitions, replicaFactor)
	if err != nil {
		// topic existed
		return &adminpb.AdminClientNewTopicResponse{
			Response: adminpb.Response_FAIL}, nil
	}

	partitionLeader := make(map[int]int)
	partitionISR := make(map[int][]int)

	for brokerID, req := range newPartitionReqMap {
		res, err := n.peerCon[brokerID].AdminClientNewPartition(context.Background(), req)
		if err != nil && res.GetResponse() == adminpb.Response_FAIL {
			// Terminate the partition creation
			// TODO: Clean up partition if the process does not complete fully (nobody care in this school project anyway)
			return &adminpb.AdminClientNewTopicResponse{
				Response: adminpb.Response_FAIL}, err
		}

		for _, partID := range req.GetPartitionID() {
			partIDInt := int(partID)
			if _, exist := partitionLeader[partIDInt]; !exist {
				partitionLeader[partIDInt] = brokerID
				continue
			}
			if _, exist := partitionISR[partIDInt]; !exist {
				partitionISR[partIDInt] = []int{brokerID}
				continue
			}
			partitionISR[partIDInt] = append(partitionISR[partIDInt], brokerID)
		}
	}

	// response
	return &adminpb.AdminClientNewTopicResponse{
		Response: adminpb.Response_SUCCESS}, nil
}

// AdminClientNewPartition create new partition
func (*Node) AdminClientNewPartition(ctx context.Context, req *adminpb.AdminClientNewPartitionRequest) (*adminpb.AdminClientNewPartitionResponse, error) {
	topicName := req.GetTopic()
	partitionID := req.GetPartitionID()

	for _, partID := range partitionID {
		err := partition.CreatePartitionDir(".", topicName, int(partID))
		if err != nil {
			return &adminpb.AdminClientNewPartitionResponse{Response: adminpb.Response_FAIL}, err
		}
	}

	return &adminpb.AdminClientNewPartitionResponse{Response: adminpb.Response_SUCCESS}, nil
}
