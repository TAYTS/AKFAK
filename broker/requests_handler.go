package broker

import (
	"AKFAK/broker/partition"
	"AKFAK/proto/adminclientpb"
	"AKFAK/proto/adminpb"
	"AKFAK/proto/clientpb"
	"AKFAK/proto/clustermetadatapb"
	"AKFAK/proto/commonpb"
	"AKFAK/proto/heartbeatspb"
	"AKFAK/proto/metadatapb"
	"AKFAK/proto/producepb"
	"AKFAK/proto/recordpb"
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
							// TODO: check if the replica connection exist if not create it as the cluster is just updated
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

	// update local cluster metadata cache
	n.ClusterMetadata.UpdateTopicState(newTopicStates)

	// send UpdateMetadata request to ZK
	n.updateZKClusterMetadata()

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
	replicaID := int32(-1)
	const MAX_MESSAGE_PER_PARTITION = 5
	doneSetup := false

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Broker %v unable to get sync request from %v\n", n.ID, replicaID)
		}

		// set the requesting replicaID
		replicaID = req.GetReplicaID()

		topicPartMapping := make(map[string][]*adminclientpb.SyncPartition)
		partDirnameFileHandlerMap := make(map[string]*recordpb.FileRecord)

		// initial setup
		if !doneSetup {
			for _, syncTp := range req.GetTopics() {
				topicName := syncTp.GetTopic()
				// setup topic to partition mapping
				topicPartMapping[topicName] = syncTp.GetPartitions()

				// setup partition directory name to file handler mapping
				for _, syncPart := range syncTp.GetPartitions() {
					partDirname := partition.ConstructPartitionDirName(topicName, int(syncPart.GetPartition()))
					logName := partition.ContructPartitionLogName(topicName)
					// create the file handler
					fileHandler, _ := recordpb.InitialiseFileRecordFromFilepath(fmt.Sprintf("%v/%v/%v", n.config.LogDir, partDirname, logName))
					// shift the file read offset to the log start offset
					fileHandler.ShiftReadOffset(syncPart.GetLogStartOffset())
					partDirnameFileHandlerMap[partDirname] = fileHandler
				}
			}
			doneSetup = true
		}

		// prepare data to send
		syncTpResp := []*adminclientpb.SyncTopicResponse{}
		syncPartResp := []*adminclientpb.SyncPartitionResponse{}
		for _, syncTp := range req.GetTopics() {
			topicName := syncTp.GetTopic()
			for _, syncPart := range syncTp.GetPartitions() {
				partID := syncPart.GetPartition()
				partDirname := partition.ConstructPartitionDirName(topicName, int(partID))

				// pull all the record batches
				rcrdBatches := make([]*recordpb.RecordBatch, 0, MAX_MESSAGE_PER_PARTITION)
				for i := 0; i < MAX_MESSAGE_PER_PARTITION; i++ {
					rcrdBatch, err := partDirnameFileHandlerMap[partDirname].ReadNextRecordBatch()
					if rcrdBatch != nil && err != nil {
						rcrdBatches = append(rcrdBatches, rcrdBatch)
					} else {
						break
					}
				}

				syncPartResp = append(syncPartResp, &adminclientpb.SyncPartitionResponse{
					Partition:      partID,
					LogStartOffset: partDirnameFileHandlerMap[partDirname].GetFileSize(),
					RecordSets:     rcrdBatches,
				})
			}
			syncTpResp = append(syncTpResp, &adminclientpb.SyncTopicResponse{
				Topic:              topicName,
				PartitionResponses: syncPartResp,
			})
		}

		forgotTps := req.GetForgottenTopics()
		if len(forgotTps) > 0 {
			for _, forgotTp := range forgotTps {
				topicName := forgotTp.GetTopic()
				for _, forgotPartID := range forgotTp.GetPartitions() {
					// update local cluster metadata
					n.ClusterMetadata.MoveBrkToISRByPartition(req.GetReplicaID(), topicName, forgotPartID)

					// update cluster metadata
					if n.ID == int(n.ClusterMetadata.GetController().GetID()) {
						n.updateZKClusterMetadata()
						n.updatePeerClusterMetadata()
					} else {
						n.updateCtrlClusterMetadata()
					}

					partDirname := partition.ConstructPartitionDirName(topicName, int(forgotPartID))
					if fileHandler, exist := partDirnameFileHandlerMap[partDirname]; exist {
						fileHandler.CloseFile()
						delete(partDirnameFileHandlerMap, partDirname)
					}
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
