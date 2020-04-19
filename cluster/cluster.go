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

// UpdateClusterMetadata replace the current MetadataCluster local cache
func (cls *Cluster) UpdateClusterMetadata(newClsMeta *clustermetadatapb.MetadataCluster) {
	cls.mux.Lock()
	cls.MetadataCluster = newClsMeta
	cls.mux.Unlock()
	cls.populateCluster()
}

func (cls *Cluster) populateCluster() {
	cls.mux.Lock()
	// update nodesByID mapping
	nodesByID := make(map[int]*clustermetadatapb.MetadataBroker)
	for _, node := range cls.MetadataCluster.GetLiveBrokers() {
		nodesByID[int(node.GetID())] = node
	}
	cls.nodesByID = nodesByID
	cls.mux.Unlock()

	// update partitionsByTopic and availablePartitionsByTopic mapping
	cls.UpdateTopicState(nil)
}
