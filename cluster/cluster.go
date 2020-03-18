package cluster

import "fmt"

// Cluster is an immutable representation of a subset of nodes, topics, and partitions in the Kafka cluster
type Cluster struct {
	nodes                      []*Node
	isBootstrapConfigured      bool
	unauthorizedTopics         map[string]string
	invalidTopics              map[string]string
	internalTopics             map[string]string
	controller                 *Node
	partitionsByTopicPartition map[*TopicPartition]*PartitionInfo
	partitionsByTopic          map[string][]*PartitionInfo
	availablePartitionByTopic  map[string][]*PartitionInfo
	partitionsByNode           map[int][]*PartitionInfo
	nodesByID                  map[int]*Node
	clusterResource            *ClusterResource
}

// InitCluster create new Cluster instance and return its pointer type
func InitCluster(
	clusterID string,
	isBootstrapConfigured bool,
	nodes []*Node,
	partitions []*PartitionInfo,
	unauthorizedTopics map[string]string,
	invalidTopics map[string]string,
	internalTopics map[string]string,
	controller *Node) *Cluster {
	cluster := &Cluster{}

	cluster.isBootstrapConfigured = isBootstrapConfigured
	cluster.clusterResource = &ClusterResource{clusterID}

	// Index the nodes for quick lookup
	tmpNodesByID := make(map[int]*Node)
	tmpPartitionsByNode := make(map[int][]*PartitionInfo)
	for _, node := range nodes {
		tmpNodesByID[node.GetID()] = node
		// Populate the map here to make it easy to add the partitions per node efficiently when iterating over the partitions
		tmpPartitionsByNode[node.GetID()] = []*PartitionInfo{}
	}
	cluster.nodesByID = tmpNodesByID

	// index the partition infos by topic, topic+partition, and node
	tmpPartitionsByTopicPartition := make(map[*TopicPartition]*PartitionInfo)
	tmpPartitionsByTopic := make(map[string][]*PartitionInfo)
	for _, partitionInfo := range partitions {
		tmpPartitionsByTopicPartition[InitTopicPartition(partitionInfo.GetTopic(), partitionInfo.GetPartition())] = partitionInfo

		partArray, exist := tmpPartitionsByTopic[partitionInfo.GetTopic()]
		if !exist {
			tmpPartitionsByTopic[partitionInfo.GetTopic()] = []*PartitionInfo{partitionInfo}
		} else {
			tmpPartitionsByTopic[partitionInfo.GetTopic()] = append(partArray, partitionInfo)
		}

		// The leader may not be known
		if partitionInfo.GetLeader() == nil || partitionInfo.GetLeader().IsEmpty() {
			continue
		}

		// If it is known, its node information should be available
		partitionsForNode := tmpPartitionsByNode[partitionInfo.GetLeader().GetID()]
		tmpPartitionsByNode[partitionInfo.GetLeader().GetID()] = append(partitionsForNode, partitionInfo)
	}

	tmpAvailablePartitionsByTopic := make(map[string][]*PartitionInfo)
	for topic, partitions := range tmpPartitionsByTopic {
		// Optimise for the common case where all partitions are available
		foundUnavailablePartition := false
		for _, partition := range partitions {
			if partition.GetLeader() == nil {
				foundUnavailablePartition = true
				break
			}
		}

		availablePartitionsForTopic := partitions
		if foundUnavailablePartition {
			availablePartitionsForTopic = make([]*PartitionInfo, 0, len(partitions))

			for _, partition := range partitions {
				if partition.GetLeader() != nil {
					availablePartitionsForTopic = append(availablePartitionsForTopic, partition)
				}
			}
		}
		tmpAvailablePartitionsByTopic[topic] = availablePartitionsForTopic
	}

	cluster.partitionsByTopicPartition = tmpPartitionsByTopicPartition
	cluster.partitionsByTopic = tmpPartitionsByTopic
	cluster.availablePartitionByTopic = tmpAvailablePartitionsByTopic
	cluster.partitionsByNode = tmpPartitionsByNode

	cluster.unauthorizedTopics = unauthorizedTopics
	cluster.invalidTopics = invalidTopics
	cluster.internalTopics = internalTopics
	cluster.controller = controller

	return cluster
}

// Empty create an empty Cluster instance with no nodes and no topic-partitions
func Empty() *Cluster {
	return InitCluster("", false, nil, nil, nil, nil, nil, nil)
}

// Bootstrap create an "bootstrap" Cluster instance using the given host/ports
func Bootstrap(addresses map[string]int) *Cluster {
	nodes := make([]*Node, 0, len(addresses))
	nodeID := -1
	for host, port := range addresses {
		nodes = append(nodes, InitNode(nodeID, host, port))
		nodeID--
	}
	return InitCluster("", true, nodes, nil, nil, nil, nil, nil)
}

// WithPartition return a copy of current cluster combined with "partitions"
func (cluster *Cluster) WithPartition(partitions map[*TopicPartition]*PartitionInfo) *Cluster {
	combinedPartitions := make([]*PartitionInfo, 0, len(cluster.partitionsByTopicPartition))

	for _, partInfo := range cluster.partitionsByTopicPartition {
		combinedPartitions = append(combinedPartitions, partInfo)
	}

	newUnauthorizedTopics := map[string]string{}
	for k, v := range cluster.unauthorizedTopics {
		newUnauthorizedTopics[k] = v
	}

	newInvalidTopics := map[string]string{}
	for k, v := range cluster.invalidTopics {
		newInvalidTopics[k] = v
	}

	newInternalTopics := map[string]string{}
	for k, v := range cluster.internalTopics {
		newInternalTopics[k] = v
	}

	return InitCluster(cluster.clusterResource.GetClusterID(), false, cluster.nodes, combinedPartitions, newUnauthorizedTopics, newInvalidTopics, newInternalTopics, cluster.controller)
}

// GetNodes return the known set of nodes
func (cluster *Cluster) GetNodes() []*Node {
	return cluster.nodes
}

// GetNodeByID return the known set of nodes
func (cluster *Cluster) GetNodeByID(ID int) *Node {
	if node, exist := cluster.nodesByID[ID]; exist {
		return node
	}
	return nil
}

// GetNodeByIDIfOnline return the node by node id if the replica for the given partition is online
func (cluster *Cluster) GetNodeByIDIfOnline(partition *TopicPartition, ID int) *Node {
	node := cluster.GetNodeByID(ID)

	nodeExistInPartition := false
	partInfo := cluster.GetPartition(partition)
	if partInfo != nil {
		for _, repNode := range partInfo.GetOfflineReplicas() {
			if *node == *repNode {
				nodeExistInPartition = true
				break
			}
		}
	}

	if node != nil && !nodeExistInPartition {
		return node
	}
	return nil
}

// GetLeaderFor return the current leader for the given topic-partition
// The node that is the leader for this topic-partition, or null if there is currently no leader
func (cluster *Cluster) GetLeaderFor(topicPartition *TopicPartition) *Node {
	info, exist := cluster.partitionsByTopicPartition[topicPartition]
	if exist {
		return info.GetLeader()
	}

	return nil
}

// GetPartition return the metadata for the specified partition, return nil if not found
func (cluster *Cluster) GetPartition(topicPartition *TopicPartition) *PartitionInfo {
	if partitionInfo, exist := cluster.partitionsByTopicPartition[topicPartition]; exist {
		return partitionInfo
	}
	return nil
}

// GetPartitionsForTopic return the list of partitions for this topic
func (cluster *Cluster) GetPartitionsForTopic(topic string) []*PartitionInfo {
	partition, exist := cluster.partitionsByTopic[topic]
	if exist {
		return partition
	}

	return []*PartitionInfo{}
}

// GetPartitionCountForTopic return the number of partitions for the given topic
// Return the number of partitions or 0 if there is no corresponding metadata
func (cluster *Cluster) GetPartitionCountForTopic(topic string) int {
	partitions, exits := cluster.partitionsByTopic[topic]
	if exits {
		return len(partitions)
	}
	return 0
}

// GetAvailablePartitionsForTopic return the list of available partitions for this topic
// Return the list of partitions
func (cluster *Cluster) GetAvailablePartitionsForTopic(topic string) []*PartitionInfo {
	partitions, exits := cluster.availablePartitionByTopic[topic]
	if exits {
		return partitions
	}
	return []*PartitionInfo{}
}

// GetPartitionsForNode return the list of partitions whose leader is this node
// Return the list of partitions
func (cluster *Cluster) GetPartitionsForNode(nodeID int) []*PartitionInfo {
	partitions, exits := cluster.partitionsByNode[nodeID]
	if exits {
		return partitions
	}
	return []*PartitionInfo{}
}

// GetTopics return a map of all topics
func (cluster *Cluster) GetTopics() map[string]string {
	topics := map[string]string{}
	for topic := range cluster.partitionsByTopic {
		topics[topic] = topic
	}
	return topics
}

// GetUnauthorizedTopics return a map of all unauthorized topics
func (cluster *Cluster) GetUnauthorizedTopics() map[string]string {
	return cluster.unauthorizedTopics
}

// GetInvalidTopics return a map of all invalid topics
func (cluster *Cluster) GetInvalidTopics() map[string]string {
	return cluster.invalidTopics
}

// GetInternalTopics return a map of all internal topics
func (cluster *Cluster) GetInternalTopics() map[string]string {
	return cluster.internalTopics
}

// GetIsBootstrapConfigured return isBootstrapConfigured
func (cluster *Cluster) GetIsBootstrapConfigured() bool {
	return cluster.isBootstrapConfigured
}

// GetClusterResource return clusterResource
func (cluster *Cluster) GetClusterResource() *ClusterResource {
	return cluster.clusterResource
}

// GetController return controller
func (cluster *Cluster) GetController() *Node {
	return cluster.controller
}

// GetController return string representation of Cluster
func (cluster *Cluster) String() string {
	partitions := make([]*PartitionInfo, 0, len(cluster.partitionsByTopicPartition))

	for _, partition := range cluster.partitionsByTopicPartition {
		partitions = append(partitions, partition)
	}

	return fmt.Sprintf("Cluster(id = %v, nodes = %v, partitions = %v, controller = %v)",
		cluster.clusterResource.GetClusterID(),
		cluster.nodes,
		partitions,
		cluster.controller)

}
