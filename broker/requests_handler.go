package broker

import (
	"AKFAK/proto/adminclientpb"
	"AKFAK/proto/clientpb"
	"AKFAK/proto/clustermetadatapb"
	"AKFAK/proto/commonpb"
	"AKFAK/proto/consumepb"
	"AKFAK/proto/consumermetadatapb"
	"AKFAK/proto/metadatapb"
	"AKFAK/proto/producepb"
	"AKFAK/proto/recordpb"
	"AKFAK/proto/zkmessagepb"
	"context"
	"errors"
	"io"
	"log"
)

// Produce used to receive message batch from Producer and forward other brokers in the cluster
func (n *Node) Produce(stream clientpb.ClientService_ProduceServer) error {
	// define the mapping for replicaConn & fileHandlers
	replicaConn := make(map[int]clientpb.ClientService_ProduceClient)
	fileHandlerMapping := make(map[int]*recordpb.FileRecord)

	// optimise the replica stream setup
	doneReplicaStreamSetup := false

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error while reading client stream: %v", err)
			cleanupProducerResource(replicaConn, fileHandlerMapping)
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

		// setup Produce RPC call to all the insync replicas if the current broker is the leader
		if !doneReplicaStreamSetup {
			for _, partState := range partitions {
				if int(partState.GetLeader()) == n.ID {
					// get all insync replicas for each partition of the specific topic
					insycBrks := partState.GetIsr()
					// create stream connection to each insync replica
					for _, brkID := range insycBrks {
						brkIDInt := int(brkID)
						if _, exist := replicaConn[brkIDInt]; !exist {
							// setup the stream connection
							stream, err := n.clientServiceClient[brkIDInt].Produce(context.Background())
							if err != nil {
								// TODO: Update controller about the fail broker
							}
							replicaConn[brkIDInt] = stream
						}
					}
				}
			}
			doneReplicaStreamSetup = true
		}

		for _, tpData := range topicData {
			for _, partState := range partitions {
				// find the partition state that match the request partition ID
				partID := int(partState.GetPartitionIndex())
				if partID == int(tpData.GetPartition()) {
					// current broker is the leader of the partition
					if int(partState.GetLeader()) == n.ID {
						// save to local
						log.Printf("Broker %v receive message for partition %v\n", n.ID, partID)
						n.writeRecordBatchToLocal(topicName, partID, fileHandlerMapping, tpData.GetRecordSet())

						// broadcast to all insync replica
						log.Printf("Broker %v broadcast message for partition %v to %v\n", n.ID, partID, partState.GetIsr())
						for _, brkID := range partState.GetIsr() {
							err := replicaConn[int(brkID)].Send(producepb.InitProduceRequest(topicName, partID, tpData.GetRecordSet().GetRecords()...))
							if err != nil {
								log.Printf("Broker %v unable to send message to replica %v\n", n.ID, brkID)
								// TODO: Update controller about the fail broker
								// TODO: Exclude the current broker from the ISR
							}

							_, err = replicaConn[int(brkID)].Recv()
						}
					} else {
						// insync replica broker, save to local
						log.Printf("Broker %v receive replica message for partition %v\n", n.ID, partID)
						n.writeRecordBatchToLocal(topicName, partID, fileHandlerMapping, tpData.GetRecordSet())
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
	cleanupProducerResource(replicaConn, fileHandlerMapping)

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
		log.Printf("Node %v received controller election request for broker %v\n", n.ID, req.GetBrokerID())
		n.initControllerRoutine()
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

	// update local cluster metadata cache
	n.ClusterMetadata.UpdateClusterMetadata(newClsInfo)

	// update client service peer connection
	n.updateClientPeerConnection()

	// controller updating peer
	if n.ID == int(n.ClusterMetadata.GetController().GetID()) {
		// update admin service peer connection
		n.updateAdminPeerConnection()

		// update peer cluster metadata
		for _, peer := range n.adminServiceClient {
			_, err := peer.UpdateMetadata(context.Background(), req)
			if err != nil {
				// TODO: update ZK about broker failure
			}
		}
	}

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

// Consume responds to pull request from consumer, sending record batch on topic-X partition-Y
func (n *Node) Consume(stream clientpb.ClientService_ConsumeServer) error {
	// TODO: find if consumer group id that is pulling messages is new (not in assignments), if so, update zookeeper.
	// TODO: retrieve message
	// TODO: update offset in consumer metadata & zk

	req, err := stream.Recv()
	if err != nil {
		return err
	}
	// Check if there is an assignment to this broker
	assignment, assignmentErr := n.checkAndGetAssignment(req)
	if assignmentErr != nil {
		return err
	}

	offset := int32(0)
	// When there is no assignment:
	// 1) Call ZK and init cache again to get latest offset for assignments
	// 2) Update MetadataAssignment and update the Broker field to the current Broker handling this consume request
	if assignment == nil {
		// Call ZK and init cache again to make sure most updated information in cache
		n.initConsumerMetadataCache()
		offset = n.ConsumerMetadata.GetOffset(assignment, req.GetGroupID())
		for _, group := range n.ConsumerMetadata.GetConsumerGroups() {
			if group.GetID() == req.GetGroupID() {
				for _, assignment := range group.GetAssignments() {
					assignedBroker := n.ID
					assignment.Broker = int32(assignedBroker)
				}
			}
		}
	} else {
		// not sure if this is required - depends on whether n.ReadRecordBatchFromLocal will keep track of offset by itself
		offset = n.ConsumerMetadata.GetOffset(assignment, req.GetGroupID())
	}


	// Retrieve and send batch record to consumer
	recordBatch, fileErr := n.ReadRecordBatchFromLocal(req.GetTopicName(), int(req.GetPartition()), int64(offset))
	if fileErr != nil {
		return fileErr
	}
	stream.Send(&consumepb.ConsumeResponse{
		TopicName:            req.GetTopicName(),
		Partition:            req.GetPartition(),
		RecordSet:            recordBatch,
	})
	// Update offset in ConsumerMetadata
	n.ConsumerMetadata.UpdateOffset(assignment, req.GetGroupID())

	// Send latest ConsumerMetadata to ZK
	newConsumerMetadataState := &consumermetadatapb.MetadataConsumerState{
		ConsumerGroups: n.ConsumerMetadata.GetConsumerGroups(),
	}
	_, zkErr := n.zkClient.UpdateConsumerMetadata(context.Background(), &zkmessagepb.UpdateConsumerMetadataRequest{
		NewState: newConsumerMetadataState,
	})
	if zkErr != nil {
		return errors.New("Could not update zookeeper")
	}
	return nil
}

// GetAssignment assigns replicas for partitions of a topic
func (n *Node) GetAssignment(ctx context.Context, req *consumepb.GetAssignmentRequest) (*consumepb.GetAssignmentResponse, error) {

	cg := req.GetGroupID()
	topicName := req.GetTopicName()

	// check if replicas already assigned to consumer group for that topic
	for _, grp := range n.ConsumerMetadata.GetConsumerGroups() {
		if grp.GetID() == cg {
			for _, assignment := range grp.GetAssignments() {
				if assignment.GetTopicName() == topicName {
					return nil, nil
				}
			}
		}
	}

	assignments := []*consumepb.MetadataAssignment{}

	// if no available partition for that topic
	partitions := n.ClusterMetadata.GetAvailablePartitionsByTopic(topicName)
	if partitions == nil {
		return nil, errors.New("Topic not available")
	}

	// loop through partitions and add an assignment for every partition
	for _, partState := range partitions {
		isrBrokers := []*clustermetadatapb.MetadataBroker{}
		brokerIDs := partState.GetIsr()
		for _, i := range brokerIDs {
			isrBrokers = append(isrBrokers, n.ClusterMetadata.GetNodesByID(int(i)))
		}
		replica := &consumepb.MetadataAssignment{
			TopicName:      topicName,
			PartitionIndex: int32(partState.GetPartitionIndex()),
			Broker:         int32(brokerIDs[0]),
			IsrBrokers:     isrBrokers,
		}
		assignments = append(assignments, replica)
	}

	// update consumer group metadata
	for _, grp := range n.ConsumerMetadata.GetConsumerGroups() {
		if grp.GetID() == cg {
			n.ConsumerMetadata.UpdateAssignments(int(cg), assignments)
		}
	}

	newConsumerMetadataState := &consumermetadatapb.MetadataConsumerState{
		ConsumerGroups: n.ConsumerMetadata.GetConsumerGroups(),
	}

	// update zookeeper
	_, err := n.zkClient.UpdateConsumerMetadata(
		context.Background(),
		&zkmessagepb.UpdateConsumerMetadataRequest{
			NewState: newConsumerMetadataState,
		})

	if err != nil {
		return nil, errors.New("Could not update zookeeper")
	}

	return &consumepb.GetAssignmentResponse{Assignments: assignments}, nil
}
