package broker

import (
	"AKFAK/proto/adminclientpb"
	"AKFAK/proto/clientpb"
	"AKFAK/proto/clustermetadatapb"
	"AKFAK/proto/commonpb"
	"AKFAK/proto/metadatapb"
	"AKFAK/proto/producepb"
	"AKFAK/proto/recordpb"
	"AKFAK/proto/zkmessagepb"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
)

// Produce used to receive message batch from Producer and forward other brokers in the cluster
func (n *Node) Produce(stream clientpb.ClientService_ProduceServer) error {
	// define the mapping for replicaConn & fileHandlers
	replicaConn := make(map[int]clientpb.ClientService_ProduceClient)
	fileHandlerMapping := make(map[int]*recordpb.FileRecord)

	for {
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

		// check if the topic exist or the partition for requested topic is available
		partitions := n.ClusterMetadata.GetAvailablePartitionsByTopic(topicName)
		if partitions == nil {
			// TODO: Test this part
			return errors.New("Topic Not Available")
		}

		// setup Produce RPC call to all the insync replicas
		for _, partState := range partitions {
			// get all insync replicas for each partition of the specific topic
			insycBrks := partState.GetIsr()
			// create stream connection to each insync replica
			for _, brkID := range insycBrks {
				// setup the stream connection
				stream, err := n.clientServiceClient[int(brkID)].Produce(context.Background())
				if err != nil {
					// TODO: Update controller about the fail broker
				}
				replicaConn[int(brkID)] = stream
			}
		}

		for _, tpData := range topicData {
			for _, partState := range partitions {
				// find the partition state that match the request partition ID
				partID := int(partState.GetPartitionIndex())
				if partID == int(tpData.GetPartition()) {
					// current broker is the leader of the partition
					if int(partState.GetLeader()) == n.ID {
						// save to local
						fmt.Printf("Broker %v receive message for partition %v\n", n.ID, partID)
						WriteRecordBatchToLocal(topicName, partID, fileHandlerMapping, tpData.GetRecordSet())

						// broadcast to all insync replica
						fmt.Printf("Broker %v broadcast message for partition %v\n", n.ID, partID)
						for _, brkID := range partState.GetIsr() {
							err := replicaConn[int(brkID)].Send(producepb.InitProduceRequest(topicName, partID, tpData.GetRecordSet().GetRecords()...))
							if err != nil {
								// TODO: Update controller about the fail broker
								// TODO: Exclude the current broker from the ISR
							}

							_, err = replicaConn[int(brkID)].Recv()
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

// WaitOnMetadata get the metadata about the kafka cluster related to the requested topic
func (n *Node) WaitOnMetadata(ctx context.Context, req *metadatapb.MetadataRequest) (*metadatapb.MetadataResponse, error) {
	// retrieve the requested topic name
	topic := req.GetTopicName()

	// get all the partitions of the requested topic from the cluster metadata cache
	partitions := n.ClusterMetadata.GetAvailablePartitionsByTopic(topic)

	// check if the partitions exist, if not return error
	if partitions == nil {
		return nil, errors.New("Topic Not Available")
	}

	// create partitions in the response
	resPartitions := make([]*metadatapb.Partition, 0, len(partitions))
	for _, partStates := range partitions {
		// get the partition info
		leader := partStates.GetLeader()
		partIdx := partStates.GetPartitionIndex()

		// update the partition of the response
		resPartitions = append(
			resPartitions,
			&metadatapb.Partition{
				PartitionIndex: partIdx,
				LeaderID:       leader,
				ReplicaNodes:   partStates.GetReplicas(),
				IsrNodes:       partStates.GetIsr(),
			})
	}

	// create brokers in the response
	resBrokers := make([]*metadatapb.Broker, 0, len(n.ClusterMetadata.GetBrokers()))
	for _, brk := range n.ClusterMetadata.GetBrokers() {
		resBrokers = append(resBrokers, &metadatapb.Broker{
			NodeID: brk.GetID(),
			Host:   brk.GetHost(),
			Port:   brk.GetPort(),
		})
	}

	metadataResp := &metadatapb.MetadataResponse{
		Brokers: resBrokers,
		Topic: &metadatapb.Topic{
			Name:       topic,
			Partitions: resPartitions,
		},
	}

	return metadataResp, nil
}

// ControllerElection used for the ZK to inform the broker to start the controller routine
func (n *Node) ControllerElection(ctx context.Context, req *adminclientpb.ControllerElectionRequest) (*adminclientpb.ControllerElectionResponse, error) {
	// get the selected brokerID to be the controller from ZK
	brokerID := int(req.GetBrokerID())

	// if the broker got selected start the controller routine
	if brokerID == n.ID {
		fmt.Printf("Node %v received controller election request for broker %v\n", n.ID, req.GetBrokerID())
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
	newPartitionReqMap, err := n.generateNewPartitionRequestData(topicName, numPartitions, replicaFactor)
	if err != nil {
		// topic existed
		return &adminclientpb.AdminClientNewTopicResponse{
			Response: &commonpb.Response{Status: commonpb.ResponseStatus_FAIL}}, err
	}

	// create MetadataTopicState to store the information about the
	// newly created topic and its partitions
	metaTopicState := &clustermetadatapb.MetadataTopicState{
		TopicName:       topicName,
		PartitionStates: make([]*clustermetadatapb.MetadataPartitionState, numPartitions),
	}

	// send request to each broker to create the partition
	for brokerID, req := range newPartitionReqMap {
		// local; create partition directly directly
		if brokerID == n.ID {
			err := n.createLocalPartitionFromReq(req)
			if err != nil {
				return &adminclientpb.AdminClientNewTopicResponse{
					Response: &commonpb.Response{Status: commonpb.ResponseStatus_FAIL}}, err
			}
		} else {
			_, err := n.adminServiceClient[brokerID].AdminClientNewPartition(context.Background(), req)
			if err != nil {
				// Terminate the partition creation
				return &adminclientpb.AdminClientNewTopicResponse{
					Response: &commonpb.Response{Status: commonpb.ResponseStatus_FAIL}}, err
			}
		}

		for _, partID := range req.GetPartitionID() {
			brokerIDint32 := int32(brokerID)

			partIDInt := int(partID)
			partitionState := metaTopicState.GetPartitionStates()[partIDInt]
			if partitionState == nil {
				replicas := make([]int32, 0, replicaFactor)
				replicas = append(replicas, brokerIDint32)

				metaTopicState.GetPartitionStates()[partIDInt] = &clustermetadatapb.MetadataPartitionState{
					TopicName:       topicName,
					PartitionIndex:  partID,
					Leader:          brokerIDint32,
					Isr:             make([]int32, 0, replicaFactor-1),
					Replicas:        replicas,
					OfflineReplicas: []int32{},
				}
			} else {
				partitionState.Isr = append(partitionState.Isr, brokerIDint32)
				partitionState.Replicas = append(partitionState.Replicas, brokerIDint32)
			}
		}
	}

	// create new MetadataCluster instead of update the local cache
	// as ZK might fail (ZK is the persistent store)
	newTopicStates := append(n.ClusterMetadata.GetTopicStates(), metaTopicState)
	newClusterState := &clustermetadatapb.MetadataCluster{
		LiveBrokers: n.ClusterMetadata.GetLiveBrokers(),
		Brokers:     n.ClusterMetadata.GetBrokers(),
		Controller:  n.ClusterMetadata.GetController(),
		TopicStates: newTopicStates,
	}

	// send UpdateMetadata request to ZK
	_, err = n.zkClient.UpdateClusterMetadata(
		context.Background(),
		&zkmessagepb.UpdateClusterMetadataRequest{
			NewClusterInfo: newClusterState,
		})
	if err != nil {
		return &adminclientpb.AdminClientNewTopicResponse{
			Response: &commonpb.Response{Status: commonpb.ResponseStatus_FAIL}}, errors.New("Fail to update ZK")
	}

	// update local cluster metadata cache
	n.ClusterMetadata.UpdateTopicState(newTopicStates)

	// send UpdateMetadata request to every live broker
	for _, brk := range n.ClusterMetadata.GetLiveBrokers() {
		brkID := int(brk.GetID())
		// update other broker
		if brkID != n.ID {
			_, err := n.adminServiceClient[brkID].UpdateMetadata(
				context.Background(),
				&adminclientpb.UpdateMetadataRequest{
					NewClusterInfo: newClusterState,
				},
			)
			if err != nil {
				// update ZK about the fail broker
			}
		}
	}

	// response
	return &adminclientpb.AdminClientNewTopicResponse{
		Response: &commonpb.Response{Status: commonpb.ResponseStatus_SUCCESS, Message: "Topic created successfully"}}, nil
}

// AdminClientNewPartition create new partition
func (n *Node) AdminClientNewPartition(ctx context.Context, req *adminclientpb.AdminClientNewPartitionRequest) (*adminclientpb.AdminClientNewPartitionResponse, error) {
	err := n.createLocalPartitionFromReq(req)
	if err != nil {
		return &adminclientpb.AdminClientNewPartitionResponse{Response: &commonpb.Response{Status: commonpb.ResponseStatus_FAIL}}, err
	}

	return &adminclientpb.AdminClientNewPartitionResponse{Response: &commonpb.Response{Status: commonpb.ResponseStatus_SUCCESS}}, nil
}

// LeaderAndIsr update the state of the local replica
func (*Node) LeaderAndIsr(ctx context.Context, req *adminclientpb.LeaderAndIsrRequest) (*adminclientpb.LeaderAndIsrResponse, error) {
	// TODO: Add update the leader and isr handler function
	return &adminclientpb.LeaderAndIsrResponse{Response: &commonpb.Response{Status: commonpb.ResponseStatus_SUCCESS}}, nil
}

// UpdateMetadata update the Metadata state of the broker
func (n *Node) UpdateMetadata(ctx context.Context, req *adminclientpb.UpdateMetadataRequest) (*adminclientpb.UpdateMetadataResponse, error) {
	newClsInfo := req.GetNewClusterInfo()

	n.ClusterMetadata.UpdateClusterMetadata(newClsInfo)

	return &adminclientpb.UpdateMetadataResponse{Response: &commonpb.Response{Status: commonpb.ResponseStatus_SUCCESS}}, nil
}

// GetController gets infomation about the controller
func (n *Node) GetController(ctx context.Context, req *adminclientpb.GetControllerRequest) (*adminclientpb.GetControllerResponse, error) {
	controller := n.ClusterMetadata.GetController()

	return &adminclientpb.GetControllerResponse{
		ControllerID: controller.GetID(),
		Host:         controller.GetHost(),
		Port:         controller.GetPort(),
	}, nil
}
