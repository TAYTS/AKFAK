package cluster

import (
	"AKFAK/proto/clustermetadatapb"
	"sync"
)

// Cluster is a wrapper to clusterMetadatapb.MetadataCluster that provide thread safe for managing the cluster metadata
type Cluster struct {
	*clustermetadatapb.MetadataCluster
	partitionsByTopic          map[string][]*clustermetadatapb.MetadataPartitionState
	availablePartitionsByTopic map[string][]*clustermetadatapb.MetadataPartitionState
	nodesByID                  map[int]*clustermetadatapb.MetadataBroker
	mux                        sync.RWMutex
}

// InitCluster create a thread safe wrapper for handling the cluster metadata
func InitCluster(clusterMetadata *clustermetadatapb.MetadataCluster) *Cluster {
	cls := &Cluster{MetadataCluster: clusterMetadata}
	cls.populateCluster()
	return cls
}

// GetPartitionsByTopic return the partitionsByTopic
func (cls *Cluster) GetPartitionsByTopic(topicName string) []*clustermetadatapb.MetadataPartitionState {
	cls.mux.RLock()
	partition, exist := cls.partitionsByTopic[topicName]
	if !exist {
		cls.mux.RUnlock()
		return nil
	}
	cls.mux.RUnlock()
	return partition
}

// GetAvailablePartitionsByTopic return the availablePartitionsByTopic
func (cls *Cluster) GetAvailablePartitionsByTopic(topicName string) []*clustermetadatapb.MetadataPartitionState {
	cls.mux.RLock()
	partition, exist := cls.availablePartitionsByTopic[topicName]
	if !exist {
		cls.mux.RUnlock()
		return nil
	}
	cls.mux.RUnlock()
	return partition
}

// GetNodesByID return the nodesByID
func (cls *Cluster) GetNodesByID(nodeID int) *clustermetadatapb.MetadataBroker {
	cls.mux.RLock()
	node, exist := cls.nodesByID[nodeID]
	if !exist {
		cls.mux.RUnlock()
		return nil
	}
	cls.mux.RUnlock()
	return node
}

// UpdateTopicState replace the current topicState
func (cls *Cluster) UpdateTopicState(newTopicStates []*clustermetadatapb.MetadataTopicState) {
	cls.mux.Lock()
	if newTopicStates != nil {
		cls.TopicStates = newTopicStates
	}
	cls.mux.Unlock()

	// refrehs the partitionsByTopic & availablePartitionsByTopic mapping
	cls.refreshPartitionTopicMapping()
}

// UpdateClusterMetadata replace the current MetadataCluster local cache
func (cls *Cluster) UpdateClusterMetadata(newClsMeta *clustermetadatapb.MetadataCluster) {
	cls.mux.Lock()
	cls.MetadataCluster = newClsMeta
	cls.mux.Unlock()
	cls.populateCluster()
}

// populateCluster refresh the internal mapping for easy access
func (cls *Cluster) populateCluster() {
	cls.mux.Lock()
	// update nodesByID mapping
	nodesByID := make(map[int]*clustermetadatapb.MetadataBroker)
	for _, node := range cls.MetadataCluster.GetLiveBrokers() {
		nodesByID[int(node.GetID())] = node
	}
	cls.nodesByID = nodesByID
	cls.mux.Unlock()

	// refresh partitionsByTopic and availablePartitionsByTopic mapping
	cls.refreshPartitionTopicMapping()
}

// refreshPartitionTopicMapping refresh the partitionsByTopic & availablePartitionsByTopic mapping
func (cls *Cluster) refreshPartitionTopicMapping() {
	cls.mux.Lock()
	partitionsByTopic := make(map[string][]*clustermetadatapb.MetadataPartitionState)
	availablePartitionsByTopic := make(map[string][]*clustermetadatapb.MetadataPartitionState)
	for _, topic := range cls.MetadataCluster.GetTopicStates() {
		partitionsByTopic[topic.GetTopicName()] = topic.GetPartitionStates()

		// create a tmp array to store all the available partition (partition with leader != -1)
		tmpAvailablePartition := make([]*clustermetadatapb.MetadataPartitionState, 0, len(topic.GetPartitionStates()))

		for _, partition := range topic.GetPartitionStates() {
			// partition that is stil available (leader != -1)
			if int(partition.GetLeader()) != -1 {
				tmpAvailablePartition = append(tmpAvailablePartition, partition)
			}
		}
		availablePartitionsByTopic[topic.GetTopicName()] = tmpAvailablePartition
	}
	cls.partitionsByTopic = partitionsByTopic
	cls.availablePartitionsByTopic = availablePartitionsByTopic
	cls.mux.Unlock()
}

// moveBrkToISRByPartition move the specified broker into the ISR of the specified topic and partition index
func (cls *Cluster) moveBrkToISRByPartition(brkID int32, topicName string, partitionIdx int32) {
	cls.mux.Lock()
	// get the partitions of the topic
	partitions := cls.partitionsByTopic[topicName]

	// update the partition
	for _, partState := range partitions {
		if partState.GetPartitionIndex() == partitionIdx {
			for idx, replicaID := range partState.GetOfflineReplicas() {
				if replicaID == brkID {
					// add the broker ID to the ISR
					partState.Isr = append(partState.GetIsr(), brkID)

					// remove the broker ID from offline replicas
					offlineReplicaCopy := partState.GetOfflineReplicas()
					partState.OfflineReplicas = append(offlineReplicaCopy[:idx], offlineReplicaCopy[idx+1:]...)

					break
				}
			}
		}
	}
	cls.mux.Unlock()
}

// moveBrkToOfflineAndElectLeader move the specified broker to offline replicas
// for all the topics and partitions and elect new leader if required
func (cls *Cluster) moveBrkToOfflineAndElectLeader(brkID int32) {
	needRefereshPartMapping := false
	cls.mux.Lock()
	for _, tpState := range cls.GetTopicStates() {
		for _, partState := range tpState.GetPartitionStates() {
			// update ISR and Offline replicas
			for idx, isrID := range partState.GetIsr() {
				if isrID == brkID {
					// add broker ID to the offline replicas
					partState.OfflineReplicas = append(partState.GetOfflineReplicas(), brkID)

					// remove the broker ID from ISR
					isrCopy := partState.GetIsr()
					partState.Isr = append(isrCopy[:idx], isrCopy[idx+1:]...)

					break
				}
			}

			// if the broker is the leader then choose the next broker in the ISR as leader
			if partState.GetLeader() == brkID {
				needRefereshPartMapping = true
				// if there is broker in the ISR select the first one as the leader
				if len(partState.GetIsr()) > 0 {
					partState.Leader = partState.GetIsr()[0]
				} else {
					// if no broker in ISR set the leader to -1
					partState.Leader = -1
				}
			}
		}
	}
	cls.mux.Unlock()

	if needRefereshPartMapping {
		cls.refreshPartitionTopicMapping()
	}
}

// checkBrokerInSync check if the broker is a insync replica for all the partition
func (cls *Cluster) checkBrokerInSync(brkID int32) bool {
	cls.mux.RLock()
	for _, tpState := range cls.GetTopicStates() {
		for _, partState := range tpState.GetPartitionStates() {
			for _, offReplica := range partState.GetOfflineReplicas() {
				// if the broker is in offlineReplica mean it is not insync
				if brkID == offReplica {
					cls.mux.RUnlock()
					return false
				}
			}
		}
	}

	cls.mux.RUnlock()
	// if the broker is not in any of the offline replicas mean it is insync
	return true
}
