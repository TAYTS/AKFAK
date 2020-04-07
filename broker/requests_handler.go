package broker

import (
	"AKFAK/broker/partition"
	"AKFAK/proto/adminclientpb"
	"AKFAK/proto/clientpb"
	"AKFAK/proto/commonpb"
	"AKFAK/proto/messagepb"
	"AKFAK/proto/metadatapb"
	"context"
	"fmt"
	"io"
	"log"
)

// MessageBatch used to send message batch to the kafka cluster
func (*Node) MessageBatch(stream clientpb.ClientService_MessageBatchServer) error {
	for {
		// TODO: implement the message batch logic
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		recordBatch := req.GetRecords()
		fmt.Println(recordBatch)

		sendErr := stream.Send(&messagepb.MessageBatchResponse{Response: &commonpb.Response{Status: commonpb.ResponseStatus_SUCCESS, Message: "Thank you"}})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", sendErr)
			return sendErr
		}
	}
}

// WaitOnMetadata get the metadata about the kafka cluster
func (*Node) WaitOnMetadata(ctx context.Context, req *metadatapb.MetadataRequest) (*metadatapb.MetadataResponse, error) {
	// TODO: implement the fetch metadata logic
	return &metadatapb.MetadataResponse{
		Brokers: []*metadatapb.Broker{
			&metadatapb.Broker{NodeID: 1, Host: "0.0.0.0", Port: 5001},
		},
	}, nil
}

// ControllerElection used for the ZK to inform the broker to start the controller routine
func (n *Node) ControllerElection(ctx context.Context, req *adminclientpb.ControllerElectionRequest) (*adminclientpb.ControllerElectionResponse, error) {
	// get the selected brokerID to be the controller from ZK
	brokerID := int(req.GetBrokerID())

	fmt.Printf("Node %v received controller election request for broker %v\n", n.ID, req.GetBrokerID())

	// if the broker got selected start the controller routine
	if brokerID == n.ID {
		n.InitControllerRoutine()
	}

	return &adminclientpb.ControllerElectionResponse{Response: &commonpb.Response{Status: commonpb.ResponseStatus_SUCCESS}}, nil
}

// AdminClientNewTopic create new topic
func (n *Node) AdminClientNewTopic(ctx context.Context, req *adminclientpb.AdminClientNewTopicRequest) (*adminclientpb.AdminClientNewTopicResponse, error) {
	// get request data
	topicName := req.GetTopic()
	numPartitions := int(req.GetNumPartitions())
	replicaFactor := int(req.GetReplicationFactor())

	// handling request
	newPartitionReqMap, err := n.newPartitionRequestData(topicName, numPartitions, replicaFactor)
	if err != nil {
		// topic existed
		return &adminclientpb.AdminClientNewTopicResponse{
			Response: &commonpb.Response{Status: commonpb.ResponseStatus_SUCCESS}}, nil
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
					return &adminclientpb.AdminClientNewTopicResponse{
						Response: &commonpb.Response{Status: commonpb.ResponseStatus_SUCCESS}}, err
				}
			}
		} else {
			res, err := n.peerCon[brokerID].AdminClientNewPartition(context.Background(), req)
			if err != nil && res.GetResponse().GetStatus() == commonpb.ResponseStatus_FAIL {
				// Terminate the partition creation
				// TODO: Clean up partition if the process does not complete fully (nobody care in this school project anyway)
				return &adminclientpb.AdminClientNewTopicResponse{
					Response: &commonpb.Response{Status: commonpb.ResponseStatus_FAIL}}, err
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
	return &adminclientpb.AdminClientNewTopicResponse{
		Response: &commonpb.Response{Status: commonpb.ResponseStatus_SUCCESS}}, nil
}

// AdminClientNewPartition create new partition
func (*Node) AdminClientNewPartition(ctx context.Context, req *adminclientpb.AdminClientNewPartitionRequest) (*adminclientpb.AdminClientNewPartitionResponse, error) {
	topicName := req.GetTopic()
	partitionID := req.GetPartitionID()

	for _, partID := range partitionID {
		// TODO: Get the log root directory
		err := partition.CreatePartitionDir(".", topicName, int(partID))
		if err != nil {
			return &adminclientpb.AdminClientNewPartitionResponse{Response: &commonpb.Response{Status: commonpb.ResponseStatus_FAIL}}, err
		}
	}

	return &adminclientpb.AdminClientNewPartitionResponse{Response: &commonpb.Response{Status: commonpb.ResponseStatus_SUCCESS}}, nil
}

// LeaderAndIsr update the state of the local replica
func (*Node) LeaderAndIsr(ctx context.Context, req *adminclientpb.LeaderAndIsrRequest) (*adminclientpb.LeaderAndIsrResponse, error) {
	// TODO: Add update the leader and isr handler function
	return &adminclientpb.LeaderAndIsrResponse{Response: &commonpb.Response{Status: commonpb.ResponseStatus_SUCCESS}}, nil
}

// UpdateMetadata update the Metadata state of the broker
func (*Node) UpdateMetadata(ctx context.Context, req *adminclientpb.UpdateMetadataRequest) (*adminclientpb.UpdateMetadataResponse, error) {
	// TODO: Add update metatdata handler function
	return &adminclientpb.UpdateMetadataResponse{Response: &commonpb.Response{Status: commonpb.ResponseStatus_SUCCESS}}, nil
}

// GetMetadata gets metadata state of broker
func (*Node) GetMetadata(ctx context.Context, req *adminclientpb.UpdateMetadataRequest) (*adminclientpb.GetMetadataResponse, error) {
	// TODO: Add get metatdata handler function
	return &adminclientpb.GetMetadataResponse{ControllerID: 1, TopicStates: []*adminclientpb.UpdateMetadataTopicState{}, LiveBrokers: []*adminclientpb.UpdateMetadataBroker{}}, nil
}