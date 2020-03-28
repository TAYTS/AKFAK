package broker

import (
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

	for brokerID, req := range newPartitionReqMap {
		res, err := n.peerCon[brokerID].AdminClientNewPartition(context.Background(), req)
		if err != nil || res.GetResponse() == adminpb.Response_FAIL {
			// Terminate the partition creation
			// TODO: Clean up partition if the process does not complete fully (nobody care in this school project anyway)
			return &adminpb.AdminClientNewTopicResponse{
				Response: adminpb.Response_FAIL}, nil
		}
	}

	// response
	return &adminpb.AdminClientNewTopicResponse{
		Response: adminpb.Response_SUCCESS}, nil
}

// AdminClientNewPartition create new partition
func (*Node) AdminClientNewPartition(ctx context.Context, req *adminpb.AdminClientNewPartitionRequest) (*adminpb.AdminClientNewPartitionResponse, error) {

}
