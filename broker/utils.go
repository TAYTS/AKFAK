package broker

import (
	"AKFAK/broker/partition"
	"AKFAK/proto/adminclientpb"
	"AKFAK/proto/adminpb"
	"AKFAK/proto/clientpb"
	"AKFAK/proto/consumepb"
	"AKFAK/proto/recordpb"
	"errors"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

// writeRecordBatchToLocal is a helper function for Produce request handler to save the RecordBatch to the local log file
func (n *Node) writeRecordBatchToLocal(topicName string, partitionID int, fileHandlerMapping map[int]*recordpb.FileRecord, recordBatch *recordpb.RecordBatch) {
	fHandler, exist := fileHandlerMapping[partitionID]
	if exist {
		fHandler.WriteToFile(recordBatch)
	} else {
		filePath := fmt.Sprintf("%v/%v/%v", n.config.LogDir, partition.ConstructPartitionDirName(topicName, partitionID), partition.ContructPartitionLogName(topicName))
		fileRecordHandler, _ := recordpb.InitialiseFileRecordFromFile(filePath)
		fileHandlerMapping[partitionID] = fileRecordHandler
		fileRecordHandler.WriteToFile(recordBatch)
	}
}

// ReadRecordBatchFromLocal is a helper function for Consume request handler to read the Record from local log file
func (n *Node) ReadRecordBatchFromLocal(topicName string, partitionID int) (*recordpb.RecordBatch, error) {
	filePath := fmt.Sprintf("%v/%v/%v", n.config.LogDir, partition.ConstructPartitionDirName(topicName, partitionID), partition.ContructPartitionLogName(topicName))
	fileRecordHandler, err := recordpb.InitialiseFileRecordFromFile(filePath)
	if err != nil {
		return nil, err
	}
	recordBatch, _ := fileRecordHandler.ReadNextRecordBatch()
	return recordBatch, nil
}

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
func (n *Node) updateAdminPeerConnection() {
	// add new connection
	for _, brk := range n.ClusterMetadata.GetLiveBrokers() {
		peerID := int(brk.GetID())
		if _, exist := n.adminServiceClient[peerID]; !exist && peerID != n.ID {
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

func (n *Node) checkAndGetAssignment(req *consumepb.ConsumeRequest) (*consumepb.MetadataAssignment, error) {
	for _, group := range n.ConsumerMetadata.ConsumerGroups {
		// check for assignments in the consumer group id
		if group.GetID() == req.GetGroupID() {
			assignments := group.GetAssignments()
			// if assignment broker id matches with its own id, return true
			for _, assignment := range assignments {
				if int(assignment.GetBroker()) == n.ID {
					return assignment, nil
				}
			}
			return nil, nil
		}
	}
	return nil, errors.New("No matching consumer group id found in consumer metadata")
}
