package broker

import (
	"AKFAK/broker/partition"
	"AKFAK/proto/adminclientpb"
	"AKFAK/proto/clientpb"
	"AKFAK/proto/commonpb"
	"AKFAK/proto/metadatapb"
	"AKFAK/proto/producepb"
	"AKFAK/proto/recordpb"
	"context"
	"fmt"
	"io"
	"log"
)

var topicMapping = map[string][]*partition.Partition{
	"topic1": []*partition.Partition{
		partition.InitPartition("topic1", 0, 0, []int{0, 1, 2}, []int{1, 2}, []int{}),
		partition.InitPartition("topic1", 1, 1, []int{0, 1, 2}, []int{0, 2}, []int{}),
		partition.InitPartition("topic1", 2, 2, []int{0, 1, 2}, []int{0, 1}, []int{}),
	},
}

var brkMapping = map[int]string{
	0: "0.0.0.0:5001",
	1: "0.0.0.0:5002",
	2: "0.0.0.0:5003",
}

// Produce used to receive message batch from Producer and forward other brokers in the cluster
func (n *Node) Produce(stream clientpb.ClientService_ProduceServer) error {
	// define the mapping for replicaConn & fileHandlers
	replicaConn := make(map[int]clientpb.ClientService_ProduceClient)
	fileHandlerMapping := make(map[int]*recordpb.FileRecord)

	for {
		// TODO: implement the message batch logic
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error while reading client stream: %v", err)
			CleanupProducerResource(replicaConn, fileHandlerMapping)
			return err
		}
		topicName := req.GetTopicName()
		topicData := req.GetTopicData()

		// setup connection with all the insync replicas
		for _, partInfo := range topicMapping[topicName] {
			// get all insync replicas for each partition of the specific topic
			insycBrks := partInfo.GetInSyncReplicas()
			// create stream connection to each insync replica
			for _, brkID := range insycBrks {
				if _, exist := replicaConn[brkID]; !exist && brkID != n.ID {
					// setup the stream connection
					stream, err := n.clientServiceClient[brkID].Produce(context.Background())
					if err != nil {
						// TODO: Fault handling
					}
					replicaConn[brkID] = stream
				}
			}
		}

		for _, tpData := range topicData {
			for _, partInfo := range topicMapping[topicName] {
				partID := partInfo.GetPartitionID()
				if partID == int(tpData.GetPartition()) {
					// current broker is the leader of the partition
					if partInfo.GetLeader() == n.ID {
						// save to local
						fmt.Printf("Broker %v receive message for partition %v\n", n.ID, partID)
						WriteRecordBatchToLocal(topicName, partID, fileHandlerMapping, tpData.GetRecordSet())

						fmt.Printf("Broker %v broadcast message for partition %v\n", n.ID, partID)
						// broadcast to all insync replica
						for _, brkID := range partInfo.GetInSyncReplicas() {
							err := replicaConn[brkID].Send(producepb.InitProduceRequest(topicName, partID, tpData.GetRecordSet().GetRecords()...))
							if err != nil {
								// TODO: handling send error
							}

							// TODO: Handling the response message for ACK, for now is just to clear the buffer
							_, err = replicaConn[brkID].Recv()
						}
					} else {
						// insync replica broker, save to local
						fmt.Printf("Broker %v receive replica message for partition %v\n", n.ID, partID)
						WriteRecordBatchToLocal(topicName, partID, fileHandlerMapping, tpData.GetRecordSet())
					}
				}
			}
		}

		// TODO: Update this when dealing with fault tolerance
		sendErr := stream.Send(&producepb.ProduceResponse{Response: &commonpb.Response{Status: commonpb.ResponseStatus_SUCCESS, Message: "Thank you"}})
		if sendErr != nil {
			log.Printf("Error while sending data to client: %v", sendErr)
			return sendErr
		}
	}

	// clean up resources
	CleanupProducerResource(replicaConn, fileHandlerMapping)

	return nil
}

// WaitOnMetadata get the metadata about the kafka cluster
func (*Node) WaitOnMetadata(ctx context.Context, req *metadatapb.MetadataRequest) (*metadatapb.MetadataResponse, error) {
	topic := req.GetTopicName()

	// TODO: Get the metadata from the cache
	metadataResp := &metadatapb.MetadataResponse{
		Brokers: []*metadatapb.Broker{
			&metadatapb.Broker{
				NodeID: 0,
				Host:   "0.0.0.0",
				Port:   5001,
			},
			&metadatapb.Broker{
				NodeID: 1,
				Host:   "0.0.0.0",
				Port:   5002,
			},
			&metadatapb.Broker{
				NodeID: 2,
				Host:   "0.0.0.0",
				Port:   5002,
			},
		},
		Topic: &metadatapb.Topic{
			Name: topic,
			Partitions: []*metadatapb.Partition{
				&metadatapb.Partition{
					PartitionIndex: 0,
					LeaderID:       0,
				},
				&metadatapb.Partition{
					PartitionIndex: 1,
					LeaderID:       1,
				},
				&metadatapb.Partition{
					PartitionIndex: 2,
					LeaderID:       2,
				},
			},
		},
	}

	return metadataResp, nil
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
			res, err := n.adminServiceClient[brokerID].AdminClientNewPartition(context.Background(), req)
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
func (*Node) UpdateMetadata(ctx context.Context, req *adminclientpb.UpdateMetadatRequest) (*adminclientpb.UpdateMetadataResponse, error) {
	// TODO: Add update metatdata handler function
	return &adminclientpb.UpdateMetadataResponse{Response: &commonpb.Response{Status: commonpb.ResponseStatus_SUCCESS}}, nil
}
