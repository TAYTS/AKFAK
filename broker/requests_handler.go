package broker

import (
	"AKFAK/broker/partition"
	"AKFAK/proto/adminpb"
	"context"
	"fmt"
)

// ControllerElection used for the ZK to inform the broker to start the controller routine
func (n *Node) ControllerElection(ctx context.Context, req *adminpb.ControllerElectionRequest) (*adminpb.ControllerElectionResponse, error) {
	// get the selected brokerID to be the controller from ZK
	brokerID := int(req.GetBrokerID())

	fmt.Printf("Node %v received controller election request for broker %v\n", n.ID, req.GetBrokerID())

	// if the broker got selected start the controller routine
	if brokerID == n.ID {
		n.InitControllerRoutine()
	}

	return &adminpb.ControllerElectionResponse{Response: adminpb.Response_SUCCESS}, nil
}

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

	// store the partition leader and isr
	partitionLeader := make(map[int]int)
	partitionISR := make(map[int][]int)

	// send request to each broker to create the partition
	for brokerID, req := range newPartitionReqMap {
		// local; create partition directly directly
		// !!! Not a good way of doing this
		if brokerID == n.ID {
			topicName := req.GetTopic()
			partitionID := req.GetPartitionID()
			fmt.Printf("Node %v: Create partition %v\n", n.ID, partitionID)
			for _, partID := range partitionID {
				// TODO: Get the log root directory
				err := partition.CreatePartitionDir(".", topicName, int(partID))
				if err != nil {
					return &adminpb.AdminClientNewTopicResponse{
						Response: adminpb.Response_FAIL}, err
				}
			}
		} else {
			res, err := n.peerCon[brokerID].AdminClientNewPartition(context.Background(), req)
			if err != nil && res.GetResponse() == adminpb.Response_FAIL {
				// Terminate the partition creation
				// TODO: Clean up partition if the process does not complete fully (nobody care in this school project anyway)
				return &adminpb.AdminClientNewTopicResponse{
					Response: adminpb.Response_FAIL}, err
			}
		}

		// update the partition leader and isr mapping
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

	// update ZK with the leader and isr

	// send LeaderAndIsrRequest to every live replica

	// send UpdateMetadata request to every live broker

	// response
	return &adminpb.AdminClientNewTopicResponse{
		Response: adminpb.Response_SUCCESS}, nil
}

// AdminClientNewPartition create new partition
func (*Node) AdminClientNewPartition(ctx context.Context, req *adminpb.AdminClientNewPartitionRequest) (*adminpb.AdminClientNewPartitionResponse, error) {
	topicName := req.GetTopic()
	partitionID := req.GetPartitionID()

	for _, partID := range partitionID {
		// TODO: Get the log root directory
		err := partition.CreatePartitionDir(".", topicName, int(partID))
		if err != nil {
			return &adminpb.AdminClientNewPartitionResponse{Response: adminpb.Response_FAIL}, err
		}
	}

	return &adminpb.AdminClientNewPartitionResponse{Response: adminpb.Response_SUCCESS}, nil
}

// LeaderAndIsr update the state of the local replica
func (*Node) LeaderAndIsr(ctx context.Context, req *adminpb.LeaderAndIsrRequest) (*adminpb.LeaderAndIsrResponse, error) {
	// TODO: Add update the leader and isr handler function
	return &adminpb.LeaderAndIsrResponse{Response: adminpb.Response_SUCCESS}, nil
}

// UpdateMetadata update the Metadata state of the broker
func (*Node) UpdateMetadata(ctx context.Context, req *adminpb.UpdateMetadatRequest) (*adminpb.UpdateMetadataResponse, error) {
	// TODO: Add update metatdata handler function
	return &adminpb.UpdateMetadataResponse{Response: adminpb.Response_SUCCESS}, nil
}
