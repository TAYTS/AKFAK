package broker

import (
	"AKFAK/broker/partition"
	"AKFAK/proto/adminclientpb"
	"AKFAK/proto/adminpb"
	"AKFAK/proto/clientpb"
	"AKFAK/proto/clustermetadatapb"
	"AKFAK/proto/commonpb"
	"AKFAK/proto/consumepb"
	"AKFAK/proto/consumermetadatapb"
	"AKFAK/proto/heartbeatspb"
	"AKFAK/proto/metadatapb"
	"AKFAK/proto/producepb"
	"AKFAK/proto/recordpb"
	"AKFAK/proto/zkmessagepb"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"
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
			cleanupProducerResource(replicaConn, fileHandlerMapping)
			return err
		}
		topicName := req.GetTopicName()
		topicData := req.GetTopicData()

		// check if the topic exist or the partition for requested topic is available
		partitions := n.ClusterMetadata.GetAvailablePartitionsByTopic(topicName)
		if partitions == nil {
			return errors.New("Topic Not Available")
		}

		for _, tpData := range topicData {
			for _, partState := range partitions {
				// find the partition state that match the request partition ID
				partID := int(partState.GetPartitionIndex())
				if partID == int(tpData.GetPartition()) {
					// update file handler mapping
					if _, exist := fileHandlerMapping[partID]; !exist {
						filePath := fmt.Sprintf("%v/%v/%v", n.config.LogDir, partition.ConstructPartitionDirName(topicName, partID), partition.ContructPartitionLogName(topicName))
						fileRecordHandler, _ := recordpb.InitialiseFileRecordFromFilepath(filePath)
						fileHandlerMapping[partID] = fileRecordHandler
					}

					// save to local
					newRcdBatch, _ := fileHandlerMapping[partID].WriteToFileByBaseOffset(tpData.GetRecordSet())

					// update the request RecordSet with the BaseOffset
					tpData.RecordSet = newRcdBatch

					// current broker is the leader of the partition
					if int(partState.GetLeader()) == n.ID {
						log.Printf("Broker %v receive message for partition %v\n", n.ID, partID)

						// broadcast to all insync replica
						log.Printf("Broker %v broadcast message for partition %v to %v\n", n.ID, partID, partState.GetIsr())

						for _, brkID := range partState.GetIsr() {
							brkIDInt := int(brkID)
							if _, exist := replicaConn[brkIDInt]; !exist {
								// setup the stream connection
								stream, err := n.clientServiceClient[brkIDInt].Produce(context.Background())
								if err != nil {
									delete(replicaConn, brkIDInt)
									continue
								} else {
									replicaConn[brkIDInt] = stream
								}
							}

							// broadcast message
							err := replicaConn[brkIDInt].Send(req)
							if err != nil {
								// this is used for the case where the stream connection was broken then the broker come back alive
								// before sending the message
								newStream, err := n.clientServiceClient[brkIDInt].Produce(context.Background())
								if err != nil {
									delete(replicaConn, brkIDInt)
									log.Printf("Broker %v unable to send message to replica %v: %v\n", n.ID, brkID, err)
									continue
								} else {
									replicaConn[brkIDInt] = newStream
									newStream.Send(producepb.InitProduceRequest(topicName, partID, tpData.GetRecordSet().GetRecords()...))
								}
							}

							// if the broker does not reply mean it has failed
							_, err = replicaConn[brkIDInt].Recv()
							if err != nil {
								delete(replicaConn, brkIDInt)
							}
						}
					} else {
						// this else case solely for logging purpose
						// insync replica broker, save to local
						log.Printf("Broker %v receive replica message for partition %v\n", n.ID, partID)
					}
				}
			}
		}

		sendErr := stream.Send(&producepb.ProduceResponse{Response: &commonpb.Response{Status: commonpb.ResponseStatus_SUCCESS, Message: "Received"}})
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

// ControllerElection used by the ZK to inform the broker to start the controller routine
func (n *Node) ControllerElection(ctx context.Context, req *adminclientpb.ControllerElectionRequest) (*adminclientpb.ControllerElectionResponse, error) {
	newClsInfo := req.GetNewClusterInfo()

	// get the selected brokerID to be the controller from ZK
	brokerID := int(newClsInfo.GetController().GetID())

	// if the broker got selected start the controller routine
	if brokerID == n.ID {
		log.Printf("Broker %v is the new controller\n", n.ID)

		// update local cluster metadata cache
		n.ClusterMetadata.UpdateClusterMetadata(newClsInfo)

		// start the controller routine
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

		for idx, partID := range req.GetPartitionID() {
			brokerIDint32 := int32(brokerID)

			partIDInt := int(partID)
			partitionState := metaTopicState.GetPartitionStates()[partIDInt]
			if partitionState == nil {
				replicas := make([]int32, 0, replicaFactor)
				replicas = append(replicas, brokerIDint32)
				isr := make([]int32, 0, replicaFactor-1)

				leader := int32(-1)
				if idx == 0 {
					leader = brokerIDint32
				} else {
					isr = append(isr, brokerIDint32)
				}

				metaTopicState.GetPartitionStates()[partIDInt] = &clustermetadatapb.MetadataPartitionState{
					TopicName:       topicName,
					PartitionIndex:  partID,
					Leader:          leader,
					Isr:             isr,
					Replicas:        replicas,
					OfflineReplicas: []int32{},
				}
			} else {
				if idx == 0 {
					partitionState.Leader = brokerIDint32
				} else {
					partitionState.Isr = append(partitionState.Isr, brokerIDint32)
				}
				partitionState.Replicas = append(partitionState.Replicas, brokerIDint32)
			}
		}
	}

	newTopicStates := append(n.ClusterMetadata.GetTopicStates(), metaTopicState)

	// update local cluster metadata cache
	n.ClusterMetadata.UpdateTopicState(newTopicStates)

	// update ZK and peer about the new cluster state
	n.handleClusterUpdateReq()

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
		newPeers := n.updateAdminPeerConnection()

		// setup hearbeats receiver with all the new peers
		n.setupPeerHeartbeatsReceiver(newPeers)

		// update peer cluster metadata
		n.updatePeerClusterMetadata()
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

// Heartbeats is used by the broker to send heartbeat to the controller every 1 second
func (n *Node) Heartbeats(stream adminpb.AdminService_HeartbeatsServer) error {
	for {
		err := stream.Send(&heartbeatspb.HeartbeatsResponse{BrokerID: int32(n.ID)})
		if err != nil {
			log.Printf("Broker %v detect controller fail\n", n.ID)
			return err
		}

		// wait for 1 second
		time.Sleep(time.Second)
	}
}

// SyncMessages is used to sync the record batch of the broker whose records are outdated
func (n *Node) SyncMessages(stream adminpb.AdminService_SyncMessagesServer) error {
	const MAX_MESSAGE_PER_PARTITION = 5

	partDirnameFileHandlerMap := make(map[string]*recordpb.FileRecord)
	partDirCompleteMap := make(map[string]bool)
	for {
		req, err := stream.Recv()
		if err != nil {
			return nil
		}

		syncTpResp := []*adminclientpb.SyncTopicResponse{}
		for _, syncTp := range req.GetTopics() {
			topicName := syncTp.GetTopic()

			syncPartResp := []*adminclientpb.SyncPartitionResponse{}
			for _, syncPart := range syncTp.GetPartitions() {
				partIdx := syncPart.GetPartition()

				// if the current broker is not the leader for this partition skip
				if !n.ClusterMetadata.CheckBrokerIsLeader(int32(n.ID), topicName, partIdx) {
					continue
				}

				partDirname := partition.ConstructPartitionDirName(topicName, int(partIdx))
				if _, exist := partDirCompleteMap[partDirname]; !exist {
					partDirCompleteMap[partDirname] = false
				}

				if !partDirCompleteMap[partDirname] {
					if _, exist := partDirnameFileHandlerMap[partDirname]; !exist {
						logName := partition.ContructPartitionLogName(topicName)
						// create the file handler
						fileHandler, _ := recordpb.InitialiseFileRecordFromFilepath(fmt.Sprintf("%v/%v/%v", n.config.LogDir, partDirname, logName))

						if fileHandler.GetLastEndOffset() == syncPart.GetLogStartOffset() {
							partDirCompleteMap[partDirname] = true
						}
						// shift the file read offset to the log start offset
						fileHandler.ShiftReadOffset(syncPart.GetLogStartOffset())
						partDirnameFileHandlerMap[partDirname] = fileHandler
					}

					// pull all the record batches
					rcrdBatches := make([]*recordpb.RecordBatch, 0, MAX_MESSAGE_PER_PARTITION)
					for i := 0; i < MAX_MESSAGE_PER_PARTITION; i++ {
						rcrdBatch, err := partDirnameFileHandlerMap[partDirname].ReadNextRecordBatch()
						if rcrdBatch != nil && err == nil {
							rcrdBatches = append(rcrdBatches, rcrdBatch)
						} else {
							break
						}
					}

					syncPartResp = append(syncPartResp, &adminclientpb.SyncPartitionResponse{
						Partition:      partIdx,
						LogStartOffset: partDirnameFileHandlerMap[partDirname].GetLastEndOffset(),
						RecordSets:     rcrdBatches,
					})
				}

			}
			if len(syncPartResp) > 0 {
				syncTpResp = append(syncTpResp, &adminclientpb.SyncTopicResponse{
					Topic:              topicName,
					PartitionResponses: syncPartResp,
				})
			}
		}

		if len(partDirnameFileHandlerMap) == 0 {
			return fmt.Errorf("Broker %v does not have or not the leader for any partitions requested", n.ID)
		}

		forgotTps := req.GetForgottenTopics()
		for _, forgotTp := range forgotTps {
			topicName := forgotTp.GetTopic()
			for _, forgotPartID := range forgotTp.GetPartitions() {
				partDirname := partition.ConstructPartitionDirName(topicName, int(forgotPartID))
				if fileHandler, exist := partDirnameFileHandlerMap[partDirname]; exist {
					// update local cluster metadata
					log.Printf("Broker %v update the ISR\n", n.ID)
					n.ClusterMetadata.MoveBrkToOnlineByPartition(req.GetReplicaID(), topicName, forgotPartID)

					// update cluster metadata
					n.handleClusterUpdateReq()

					// update partDirCompleteMap
					partDirCompleteMap[partDirname] = true

					fileHandler.CloseFile()
					delete(partDirnameFileHandlerMap, partDirname)
				}
			}
		}

		// sending response
		resp := &adminclientpb.SyncMessagesResponse{
			Responses: syncTpResp,
		}
		stream.Send(resp)
	}
}

// Consume responds to pull request from consumer, sending record batch on topic-X partition-Y
func (n *Node) Consume(stream clientpb.ClientService_ConsumeServer) error {
	req, err := stream.Recv()
	if err != nil {
		log.Printf("Error while reading client stream: %v", err)
		return err
	}
	// Check if consumer group has an existing assignment to this broker
	assignment, assignmentErr := n.checkAndGetAssignment(req)
	if assignmentErr != nil {
		return assignmentErr
	}
	offset := int32(0)
	// When there is no assignment:
	if assignment == nil {
		// Call ZK and init cache again to make sure most updated information in cache
		n.initConsumerMetadataCache()
		offset = n.ConsumerMetadata.GetOffset(assignment, req.GetGroupID())
		// update the metadata in node's own metadatacache
		for _, group := range n.ConsumerMetadata.GetConsumerGroups() {
			if group.GetID() == req.GetGroupID() {
				for _, assignment := range group.GetAssignments() {
					assignment.Broker = int32(n.ID)
				}
			}
		}
	} else {
		offset = n.ConsumerMetadata.GetOffset(assignment, req.GetGroupID())
	}
	log.Println("Offset:", offset)

	// Retrieve and send batch record to consumer
	recordBatch, fileErr := n.ReadRecordBatchFromLocal(req.GetTopicName(), int(req.GetPartition()), int64(offset))
	if fileErr != nil {
		log.Println("error when reading from file:", fileErr)
		return fileErr
	}
	stream.Send(&consumepb.ConsumeResponse{
		TopicName: req.GetTopicName(),
		Partition: req.GetPartition(),
		RecordSet: recordBatch,
	})
	// Update offset in ConsumerMetadata
	n.ConsumerMetadata.UpdateOffset(assignment, req.GetGroupID())

	// Send latest ConsumerMetadata to ZK
	newConsumerMetadataState := &consumermetadatapb.MetadataConsumerState{
		ConsumerGroups: n.ConsumerMetadata.GetConsumerGroups(),
	}

	// update zookeeper on new assignment
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
	assignments := []*consumepb.MetadataAssignment{}

	// check if replicas already assigned to consumer group for that topic
	for _, grp := range n.ConsumerMetadata.GetConsumerGroups() {
		if grp.GetID() == cg {
			for i, assignment := range grp.GetAssignments() {
				if assignment.GetTopicName() == topicName {
					assignments = append(assignments, assignment)
				}
				if i == (len(grp.GetAssignments())-1) && len(assignments) != 0 {
					return &consumepb.GetAssignmentResponse{Assignments: assignments}, nil
					// decide not to throw error so that consumer can get assignments again
					// return &consumepb.GetAssignmentResponse{Assignments: assignments}, errors.New("You already have an assignment for this topic!")
				}
			}
		}
	}

	// if no available partition for that topic
	partitions := n.ClusterMetadata.GetAvailablePartitionsByTopic(topicName)
	if partitions == nil {
		return nil, errors.New("Topic not available")
	}

	// loop through partitions and add an assignment for every partition
	for _, partState := range partitions {
		isrBrokers := []*clustermetadatapb.MetadataBroker{}
		brokerIDs := partState.GetReplicas()
		if len(brokerIDs) == 0 {
			return nil, errors.New("No broker is storing data on this topic")
		}
		for _, i := range brokerIDs {
			isrBrokers = append(isrBrokers, n.ClusterMetadata.GetNodesByID(int(i)))
		}
		replica := &consumepb.MetadataAssignment{
			TopicName:      topicName,
			PartitionIndex: int32(partState.GetPartitionIndex()),
			//assign the first broker to be the one the consumer should contact
			Broker:     int32(brokerIDs[0]),
			IsrBrokers: isrBrokers,
		}
		assignments = append(assignments, replica)
	}

	// if no consumergroups available in metadata
	if len(n.ConsumerMetadata.GetConsumerGroups()) == 0 {
		// add this consumer group
		n.ConsumerMetadata.AddAssignments(int(cg), assignments)
	}

	// update consumer group metadata
	for i, grp := range n.ConsumerMetadata.GetConsumerGroups() {
		if grp.GetID() == cg {
			n.ConsumerMetadata.UpdateAssignments(int(cg), assignments)
			break
		} else if i == len(n.ConsumerMetadata.GetConsumerGroups())-1 {
			// consumergroup not in metadata yet
			n.ConsumerMetadata.AddAssignments(int(cg), assignments)
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
