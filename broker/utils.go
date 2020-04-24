package broker

import (
	"AKFAK/broker/partition"
	"AKFAK/proto/adminclientpb"
	"AKFAK/proto/adminpb"
	"AKFAK/proto/clientpb"
	"AKFAK/proto/heartbeatspb"
	"AKFAK/proto/recordpb"
	"AKFAK/proto/zkmessagepb"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
)

// cleanupProducerResource help to clean up the Producer resources
func cleanupProducerResource(replicaConn map[int]clientpb.ClientService_ProduceClient, fileHandlerMapping map[int]*recordpb.FileRecord) {
	for _, rCon := range replicaConn {
		err := rCon.CloseSend()
		if err != nil {
			log.Printf("Closing connection error: %v\n", err)
		}
	}

	for _, fileHandler := range fileHandlerMapping {
		err := fileHandler.CloseFile()
		if err != nil {
			log.Printf("Closing file error: %v\n", err)
		}
	}
}

// generateNewPartitionRequestData create a mapping of the brokerID and the create new partition RPC request
func (n *Node) generateNewPartitionRequestData(topicName string, numPartitions int, replicaFactor int) (map[int]*adminclientpb.AdminClientNewPartitionRequest, error) {
	// check if the topic exist
	topicExist := false
	if len(n.ClusterMetadata.GetPartitionsByTopic(topicName)) > 0 {
		topicExist = true
	}

	numBrokers := len(n.ClusterMetadata.GetLiveBrokers())

	// check if the current available broker can support replication
	if numBrokers < replicaFactor {
		return nil, errors.New("Number of brokers available is less than the replica factor")
	}

	if !topicExist {
		newTopicPartitionRequests := make(map[int]*adminclientpb.AdminClientNewPartitionRequest)

		for partID := 0; partID < numPartitions; partID++ {
			// distribute the partitions among the brokers
			brokerID := int(n.ClusterMetadata.GetLiveBrokers()[partID%numBrokers].GetID())
			leader := brokerID
			for replicaIdx := 0; replicaIdx < replicaFactor; replicaIdx++ {
				// distribute the replicas among the brokers
				replicaBrokerIdx := (brokerID + replicaIdx + partID) % numBrokers
				replicaBrokerID := int(n.ClusterMetadata.GetLiveBrokers()[replicaBrokerIdx].GetID())

				request, exist := newTopicPartitionRequests[replicaBrokerID]
				if exist {
					if replicaBrokerID == leader {
						request.PartitionID = append([]int32{int32(partID)}, request.PartitionID...)
					} else {
						request.PartitionID = append(request.PartitionID, int32(partID))
					}
				} else {
					request := &adminclientpb.AdminClientNewPartitionRequest{
						Topic:       topicName,
						PartitionID: []int32{int32(partID)},
						ReplicaID:   int32(replicaBrokerID),
					}
					newTopicPartitionRequests[replicaBrokerID] = request
				}
			}
		}

		return newTopicPartitionRequests, nil
	}
	return nil, errors.New("Topic already exist")
}

// createLocalPartitionFromReq take the create new partition RPC request and create the local partition directory and log file
func (n *Node) createLocalPartitionFromReq(req *adminclientpb.AdminClientNewPartitionRequest) error {
	topicName := req.GetTopic()
	partitionID := req.GetPartitionID()

	log.Printf("Node %v: Create partition %v\n", n.ID, partitionID)

	rootPath := n.config.LogDir
	for _, partID := range partitionID {
		err := partition.CreatePartitionDir(rootPath, topicName, int(partID))
		if err != nil {
			log.Println("create partition err", err)
			return err
		}
	}
	return nil
}

// updateAdminPeerConnection is used by controller to store all the peer gRPC connections for admin service
func (n *Node) updateAdminPeerConnection() []int {
	newPeer := []int{}
	// add new connection
	for _, brk := range n.ClusterMetadata.GetLiveBrokers() {
		peerID := int(brk.GetID())
		if _, exist := n.adminServiceClient[peerID]; !exist && peerID != n.ID {
			newPeer = append(newPeer, peerID)
			peerAddr := fmt.Sprintf("%v:%v", brk.GetHost(), brk.GetPort())
			clientCon, err := grpc.Dial(peerAddr, grpc.WithInsecure())
			if err != nil {
				log.Printf("Fail to connect to %v: %v\n", peerAddr, err)
				// TODO: Update the ZK about the fail node
				continue
			}
			adminServiceClient := adminpb.NewAdminServiceClient(clientCon)
			n.adminServiceClient[peerID] = adminServiceClient
		}
	}

	// remove dead broker
	for ID := range n.adminServiceClient {
		if n.ClusterMetadata.GetNodesByID(ID) == nil {
			delete(n.adminServiceClient, ID)
		}
	}

	return newPeer
}

// updateClientPeerConnection is used to update the mapping that store all the peer gRPC connections for client service
func (n *Node) updateClientPeerConnection() {
	// add new connection
	for _, brk := range n.ClusterMetadata.GetLiveBrokers() {
		peerID := int(brk.GetID())
		if _, exist := n.clientServiceClient[peerID]; !exist && peerID != n.ID {
			peerAddr := fmt.Sprintf("%v:%v", brk.GetHost(), brk.GetPort())
			clientCon, err := grpc.Dial(peerAddr, grpc.WithInsecure())
			if err != nil {
				log.Printf("Fail to connect to %v: %v\n", peerAddr, err)
				// TODO: Update the ZK about the fail node
				continue
			}
			clientServiceClient := clientpb.NewClientServiceClient(clientCon)
			n.clientServiceClient[peerID] = clientServiceClient
		}
	}

	// remove dead broker
	if len(n.ClusterMetadata.GetLiveBrokers()) != len(n.clientServiceClient) {
		for ID := range n.clientServiceClient {
			if n.ClusterMetadata.GetNodesByID(ID) == nil {
				delete(n.clientServiceClient, ID)
			}
		}
	}
}

// setupPeerHeartbeatsReceiver create a client stream to the peer broker
func (n *Node) setupPeerHeartbeatsReceiver(peerIDs []int) {
	for _, peerID := range peerIDs {
		if adminClient, exist := n.adminServiceClient[peerID]; exist {
			log.Printf("Controller setup heartbeat connection to Broker %v\n", peerID)
			stream, err := adminClient.Heartbeats(context.Background())
			if err != nil {
				n.handleBrokerFailure(int32(peerID))
				continue
			}

			// start new routine to handle the heartbeat of each peer
			go func() {
				for {
					_, err := stream.Recv()
					if err != nil {
						n.handleBrokerFailure(int32(peerID))
						break
					}
					log.Printf("Controller receive heartbeats response from  Broker %v\n", peerID)
				}
			}()
		}
	}
}

// setupZKHeartbeatsConnection used by controller to send heartbeat request to ZK
func (n *Node) setupZKHeartbeatsRequest() {
	stream, err := n.zkClient.Heartbeats(context.Background())
	if err != nil {
		// Should not happened in this implementation because we assume
		// ZK is always available
		log.Printf("Controller detect ZK has failed\n")
	}

	// send the heartbeats request to ZK every 1 second
	for {
		err := stream.Send(&heartbeatspb.HeartbeatsRequest{BrokerID: int32(n.ID)})
		if err != nil {
			// Should not happened in this implementation because we assume
			// ZK is always available
			log.Printf("Controller detect ZK has failed\n")
		}

		// wait for 1 second
		time.Sleep(time.Second)
	}
}

// handleClusterUpdateReq is update the cluster state based on the role of the broker
// - controller: update ZK then peer
// - broker: update controller
func (n *Node) handleClusterUpdateReq() {
	if n.ID == int(n.ClusterMetadata.GetController().GetID()) {
		n.updateZKClusterMetadata()
		n.updatePeerClusterMetadata()
	} else {
		n.updateCtrlClusterMetadata()
	}
}

// updateZKClusterMetadata is used by controller to update ZK about the current cluster state
func (n *Node) updateZKClusterMetadata() {
	log.Println("Controller update ZK about new cluster state")

	for {
		if n.zkClient != nil {
			n.zkClient.UpdateClusterMetadata(
				context.Background(),
				&zkmessagepb.UpdateClusterMetadataRequest{
					NewClusterInfo: n.ClusterMetadata.MetadataCluster,
				})
			break
		} else {
			time.Sleep(300 * time.Millisecond)
		}
	}
}

// updateControllerMetadata used by broker to update controller about the new cluster metadata
func (n *Node) updateCtrlClusterMetadata() {
	log.Printf("Broker %v update controller about new cluster state\n", n.ID)

	for len(n.ClusterMetadata.GetLiveBrokers()) > 1 {
		ctrlAddr := fmt.Sprintf("%v:%v", n.ClusterMetadata.GetController().GetHost(), n.ClusterMetadata.GetController().GetPort())

		ctrlCon, err := grpc.Dial(ctrlAddr, grpc.WithInsecure())
		defer ctrlCon.Close()
		if err != nil {
			// broker ignore controller failure(none of their business)
			// try update controller until success
			log.Printf("Broker %v send update req to controller about new cluster state: %v", n.ID, err)
			time.Sleep(300 * time.Millisecond)
			continue
		}

		adminServiceClient := adminpb.NewAdminServiceClient(ctrlCon)
		// broker ignore controller failure
		_, err = adminServiceClient.UpdateMetadata(context.Background(), &adminclientpb.UpdateMetadataRequest{
			NewClusterInfo: n.ClusterMetadata.MetadataCluster,
		})
		log.Printf("Broker %v send update req to controller about new cluster state\n", n.ID)

		if err != nil {
			// broker ignore controller failure(none of their business)
			// try update controller until success
			log.Printf("Broker %v send update req to controller about new cluster state: %v", n.ID, err)
			time.Sleep(300 * time.Millisecond)
			continue
		}
		break
	}
	log.Printf("Broker %v done update controller about new cluster state\n", n.ID)
}

// updatePeerClusterMetadata used by controller to update peer about the new cluster metadata
func (n *Node) updatePeerClusterMetadata() {
	log.Println("Controller update peers about the new Cluster state")
	req := &adminclientpb.UpdateMetadataRequest{
		NewClusterInfo: n.ClusterMetadata.MetadataCluster,
	}

	ctrlID := n.ClusterMetadata.GetController().GetID()

	// update all live brokers
	for _, brk := range n.ClusterMetadata.GetLiveBrokers() {
		brkID := brk.GetID()
		if brkID != ctrlID {
			_, err := n.adminServiceClient[int(brkID)].UpdateMetadata(context.Background(), req)
			if err != nil {
				n.handleBrokerFailure(brkID)
			}
		}
	}
}

// handleBrokerFailure is used by controller to handle broker failure
func (n *Node) handleBrokerFailure(brkID int32) {
	log.Printf("Detect Broker %v failure, updating the Cluster state\n", brkID)

	// remove the broker from the ISR and elect new leader if required
	n.ClusterMetadata.MoveBrkToOfflineAndElectLeader(brkID)

	// clean up admin service peer connections
	n.updateAdminPeerConnection()

	// clean up client service peer connections
	n.updateClientPeerConnection()

	// update cluster state
	n.handleClusterUpdateReq()
}

// syncLocalPartition is used by broker on startup to sync the local partition with other replicas
func (n *Node) syncLocalPartition() {
	// get all partitions need to sync
	offlineTopics := n.ClusterMetadata.GetBrkOfflineTopics(int32(n.ID))

	// if there is no offline topic return
	if offlineTopics == nil {
		return
	}

	log.Printf("Broker %v start syncing the partition\n", n.ID)

	// prepare the sync requests
	partFileHandleMap := make(map[string]*recordpb.FileRecord)
	partEndOffsetMap := make(map[string]int64)
	allSyncTp := []*adminclientpb.SyncTopic{}
	for _, tpState := range offlineTopics {
		topicName := tpState.GetTopicName()

		// sync partition buffer for the current topic
		tpSyncParts := []*adminclientpb.SyncPartition{}
		for _, partState := range tpState.GetPartitionStates() {
			leaderID := partState.GetLeader()
			partIdx := partState.GetPartitionIndex()

			if leaderID == -1 {
				// if there is no leader available at the first place elect itself
				n.ClusterMetadata.MoveBrkToOnlineByPartition(int32(n.ID), topicName, partIdx)
				n.handleClusterUpdateReq()
				continue
			} else if int(leaderID) == n.ID {
				// if the current broker is the leader not need to update
				continue
			}

			// create the file handler
			partDirname := partition.ConstructPartitionDirName(topicName, int(partIdx))
			logName := partition.ContructPartitionLogName(topicName)
			fileHandler, _ := recordpb.InitialiseFileRecordFromFilepath(fmt.Sprintf("%v/%v/%v", n.config.LogDir, partDirname, logName))
			lastOffset := fileHandler.GetLastEndOffset()
			partFileHandleMap[partDirname] = fileHandler
			partEndOffsetMap[partDirname] = lastOffset

			// create sync partition
			syncPart := &adminclientpb.SyncPartition{
				Partition:      partState.GetPartitionIndex(),
				LogStartOffset: lastOffset,
			}

			// add to the sync partition buffer for the current topic
			tpSyncParts = append(tpSyncParts, syncPart)
		}

		allSyncTp = append(allSyncTp, &adminclientpb.SyncTopic{
			Topic:      topicName,
			Partitions: tpSyncParts,
		})
	}

	if len(allSyncTp) == 0 {
		return
	}

	syncReq := &adminclientpb.SyncMessagesRequest{
		ReplicaID: int32(n.ID),
		Topics:    allSyncTp,
	}

	var brkWG sync.WaitGroup

	// send the sync request to all live brokers, since only the leader brokers will reply
	for _, brk := range n.ClusterMetadata.GetLiveBrokers() {
		peerAddr := fmt.Sprintf("%v:%v", brk.GetHost(), brk.GetPort())
		clientCon, err := grpc.Dial(peerAddr, grpc.WithInsecure())
		defer clientCon.Close()
		if err != nil {
			continue
		}

		adminServiceClient := adminpb.NewAdminServiceClient(clientCon)
		stream, err := adminServiceClient.SyncMessages(context.Background())
		if err != nil {
			continue
		}
		brkWG.Add(1)

		go func() {
			for {
				stream.Send(syncReq)
				resp, err := stream.Recv()
				if err != nil {
					brkWG.Done()
					break
				}
				data := resp.GetResponses()
				for _, tpResp := range data {
					topicName := tpResp.GetTopic()
					for _, partResp := range tpResp.GetPartitionResponses() {
						partIdx := partResp.GetPartition()
						endOffset := partResp.GetLogStartOffset()
						rcrdBatches := partResp.GetRecordSets()

						partDirname := partition.ConstructPartitionDirName(topicName, int(partIdx))
						partEndOffsetMap[partDirname] = endOffset
						fileHandler := partFileHandleMap[partDirname]
						for _, rcrdBatch := range rcrdBatches {
							fileHandler.WriteToFileByBaseOffset(rcrdBatch)
						}

						// update the req forgetten topic
						newEndOffset := fileHandler.GetLastEndOffset()
						if newEndOffset >= partEndOffsetMap[partDirname] {
							forgotTpExist := false
							for _, forgetTp := range syncReq.GetForgottenTopics() {
								if forgetTp.GetTopic() == topicName {
									forgotTpExist = true
									forgotPartExist := false
									for _, partID := range forgetTp.GetPartitions() {
										if partID == partIdx {
											forgotPartExist = true
										}
									}
									if !forgotPartExist {
										forgetTp.Partitions = append(forgetTp.GetPartitions(), partIdx)
									}
								}
							}
							if !forgotTpExist {
								syncReq.ForgottenTopics = append(syncReq.GetForgottenTopics(), &adminclientpb.SyncForgottenTopicsData{
									Topic:      topicName,
									Partitions: []int32{partIdx},
								})
							}
						}
					}
				}
			}
		}()
	}
	brkWG.Wait()

	log.Printf("Broker %v done sending sync req\n", n.ID)

	// clean up resource
	for _, fileHandler := range partFileHandleMap {
		fileHandler.CloseFile()
	}
}
