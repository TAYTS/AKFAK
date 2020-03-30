package broker

import (
	"AKFAK/proto/adminpb"
	"errors"
)

func (n *Node) newPartitionRequestData(topicName string, numPartitions int, replicaFactor int) (map[int]*adminpb.AdminClientNewPartitionRequest, error) {
	// check if the topic exist
	topicExist := false // TODO: Verify if the topic exist

	// TODO: Get from ZK
	numBrokers := len(n.peerCon) + 1

	if !topicExist {
		newTopicPartitionRequests := make(map[int]*adminpb.AdminClientNewPartitionRequest)

		for partID := 0; partID < numPartitions; partID++ {
			// TODO: Update this after implement ZK Metadata
			// distribute the partitions among the brokers
			brokerID := partID % numBrokers
			for replicaIdx := 0; replicaIdx < replicaFactor; replicaIdx++ {
				// distribute the replicas among the brokers
				replicaBrokerID := (brokerID + replicaIdx + partID) % numBrokers

				request, exist := newTopicPartitionRequests[replicaBrokerID]
				if exist {
					request.PartitionID = append(request.PartitionID, int32(partID))
				} else {
					request := &adminpb.AdminClientNewPartitionRequest{
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
