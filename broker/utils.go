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
	"io"
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

			for replicaIdx := 0; replicaIdx < replicaFactor; replicaIdx++ {
				// distribute the replicas among the brokers
				replicaBrokerIdx := (brokerID + replicaIdx + partID) % numBrokers
				replicaBrokerID := int(n.ClusterMetadata.GetLiveBrokers()[replicaBrokerIdx].GetID())

				request, exist := newTopicPartitionRequests[replicaBrokerID]
				if exist {
					request.PartitionID = append(request.PartitionID, int32(partID))
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
	if len(n.ClusterMetadata.GetLiveBrokers()) != len(n.adminServiceClient) {
		for ID := range n.adminServiceClient {
			if n.ClusterMetadata.GetNodesByID(ID) == nil {
				delete(n.adminServiceClient, ID)
			}
		}
	}

	return newPeer
}

// updateClientPeerConnection is used to store all the peer gRPC connections for client service
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
				log.Printf("Controller detect Broker %v has failed\n", peerID)
				// TODO: update ZK about fail broker
			}

			// start new routine to handle the heartbeat of each peer
			go func() {
				for {
					_, err := stream.Recv()
					if err != nil {
						log.Printf("Controller detect Broker %v has failed\n", peerID)
						// TODO: update ZK about fail broker
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

// updateZKClusterMetadata is used by controller to update ZK about the current cluster state
func (n *Node) updateZKClusterMetadata() {
	log.Println("Controller update ZK about new cluster state")

	n.zkClient.UpdateClusterMetadata(
		context.Background(),
		&zkmessagepb.UpdateClusterMetadataRequest{
			NewClusterInfo: n.ClusterMetadata.MetadataCluster,
		})
}

// updateControllerMetadata used by broker to update controller about the new cluster metadata
func (n *Node) updateCtrlClusterMetadata() {
	log.Printf("Broker %v update controller about new cluster state\n", n.ID)

	for {
		ctrlAddr := fmt.Sprintf("%v:%v", n.ClusterMetadata.GetController().GetHost(), n.ClusterMetadata.GetController().GetHost())

		ctrlCon, err := grpc.Dial(ctrlAddr, grpc.WithInsecure())
		if err != nil {
			// broker ignore controller failure(none of their business)
			// try update controller until success
			ctrlCon.Close()
			time.Sleep(300 * time.Millisecond)
			continue
		}

		adminServiceClient := adminpb.NewAdminServiceClient(ctrlCon)

		// broker ignore controller failure
		_, err = adminServiceClient.UpdateMetadata(context.Background(), &adminclientpb.UpdateMetadataRequest{
			NewClusterInfo: n.ClusterMetadata.MetadataCluster,
		})
		if err != nil {
			// broker ignore controller failure(none of their business)
			// try update controller until success
			ctrlCon.Close()
			time.Sleep(300 * time.Millisecond)
			continue
		}
		break
	}
}

// updatePeerClusterMetadata used by controller to update peer about the new cluster metadata
func (n *Node) updatePeerClusterMetadata() {
	req := &adminclientpb.UpdateMetadataRequest{
		NewClusterInfo: n.ClusterMetadata.MetadataCluster,
	}
	for _, peer := range n.adminServiceClient {
		_, err := peer.UpdateMetadata(context.Background(), req)
		if err != nil {
			// TODO: update ZK about broker failure
		}
	}
}

// handleBrokerFailure is used by controller to handle broker failure
func (n *Node) handleBrokerFailure(brkID int32) {
	// remove the broker from the ISR and elect new leader if required
	n.ClusterMetadata.MoveBrkToOfflineAndElectLeader(brkID)

	// update ZK about the new cluster state
	n.updateZKClusterMetadata()
}

// syncLocalPartition is used by broker on startup to sync the local partition with other replicas
func (n *Node) syncLocalPartition() {
	log.Printf("Broker %v starts to sync the local partition\n", n.ID)
	// get all partitions need to sync
	offlineTopics := n.ClusterMetadata.GetBrkOfflineTopics(int32(n.ID))

	// if there is no offline topic return
	if offlineTopics == nil {
		return
	}

	// prepare the sync requests
	leaderSyncReqMap := make(map[int32]*adminclientpb.SyncMessagesRequest)
	partFileHandleMap := make(map[string]*recordpb.FileRecord)
	for _, tpState := range offlineTopics {
		topicName := tpState.GetTopicName()

		// create mapping for keeping track the leader of each partition
		leaderPartsMap := make(map[int32][]*adminclientpb.SyncPartition)
		for _, partState := range tpState.GetPartitionStates() {
			// TODO: leaderID = -1
			leaderID := partState.GetLeader()
			partIdx := partState.GetPartitionIndex()

			// create the file handler
			partDirname := partition.ConstructPartitionDirName(topicName, int(partIdx))
			logName := partition.ContructPartitionLogName(topicName)
			fileHandler, _ := recordpb.InitialiseFileRecordFromFilepath(fmt.Sprintf("%v/%v/%v", n.config.LogDir, partDirname, logName))
			partFileHandleMap[partDirname] = fileHandler

			// create sync partition
			syncPart := &adminclientpb.SyncPartition{
				Partition:      partState.GetPartitionIndex(),
				LogStartOffset: fileHandler.GetFileSize(),
			}

			// add sync partition to the respective leader
			if parts, exist := leaderPartsMap[leaderID]; exist {
				parts = append(parts, syncPart)
			} else {
				leaderPartsMap[leaderID] = []*adminclientpb.SyncPartition{syncPart}
			}
		}

		// add/update sync req for the respective leader
		for leader, syncParts := range leaderPartsMap {
			syncTopic := &adminclientpb.SyncTopic{
				Topic:      topicName,
				Partitions: syncParts,
			}
			if syncReq, exist := leaderSyncReqMap[leader]; exist {
				syncReq.Topics = append(syncReq.GetTopics(), syncTopic)
			} else {
				leaderSyncReqMap[leader] = &adminclientpb.SyncMessagesRequest{
					ReplicaID: int32(n.ID),
					Topics:    []*adminclientpb.SyncTopic{syncTopic},
				}
			}
		}
	}

	// wait group for waiting all leader syncing is complete
	var syncWG sync.WaitGroup

	// sending sync request to each leader
	for leader, syncReq := range leaderSyncReqMap {
		brk := n.ClusterMetadata.GetNodesByID(int(leader))
		peerAddr := fmt.Sprintf("%v:%v", brk.GetHost(), brk.GetPort())
		clientCon, err := grpc.Dial(peerAddr, grpc.WithInsecure())
		defer clientCon.Close()
		if err != nil {
			log.Printf("Fail to connect to %v: %v\n", peerAddr, err)
			// TODO: Update the ZK about the fail node
			continue
		}
		adminServiceClient := adminpb.NewAdminServiceClient(clientCon)
		stream, err := adminServiceClient.SyncMessages(context.Background())
		if err != nil {
			// TODO: Handle broker failure
			continue
		}

		// wait for one more leader
		syncWG.Add(1)

		go func(syncReq *adminclientpb.SyncMessagesRequest) {
			// buffer used to keep track all partitions end offset
			partEndOffsetMap := make(map[string]int64)
			for {
				// send the sync request
				stream.Send(syncReq)

				// clear the forgotten topics
				syncReq.ForgottenTopics = []*adminclientpb.SyncForgottenTopicsData{}

				// get response from leader
				resp, err := stream.Recv()
				if err == io.EOF {
					// done syncing
					syncWG.Done()
					break
				}
				syncResp := resp.GetResponses()

				// writting the record batches to local log files
				for _, tp := range syncResp {
					topicName := tp.GetTopic()
					for _, part := range tp.GetPartitionResponses() {
						partIdx := part.GetPartition()
						partDirname := partition.ConstructPartitionDirName(topicName, int(partIdx))
						if offset, exist := partEndOffsetMap[partDirname]; exist {
							offset = part.GetLogStartOffset()
						} else {
							partEndOffsetMap[partDirname] = offset
						}

						// write record batches and get the updated record end offset
						var newEndOffset int64
						for _, record := range part.GetRecordSets() {
							partFileHandleMap[partDirname].WriteToFileByBaseOffset(record)
							newEndOffset = partFileHandleMap[partDirname].GetFileSize()
						}

						// update the req forgetten topic
						if newEndOffset == partEndOffsetMap[partDirname] {
							forgotTpExist := false
							for _, forgetTp := range syncReq.GetForgottenTopics() {
								if forgetTp.GetTopic() == topicName {
									forgotTpExist = true
									forgetTp.Partitions = append(forgetTp.GetPartitions(), partIdx)
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
		}(syncReq)
	}

	// wait for all the syncing with leaders to be done
	syncWG.Wait()

	// clean up resource
	for _, fileHandler := range partFileHandleMap {
		fileHandler.CloseFile()
	}
}
