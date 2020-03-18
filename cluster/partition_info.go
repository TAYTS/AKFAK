package cluster

import "fmt"

// PartitionInfo describe per-partition state in the MetadataResponse
type PartitionInfo struct {
	topic           string
	partition       int
	leader          *Node
	replicas        []*Node
	inSyncReplicas  []*Node
	offlineReplicas []*Node
}

///////////////////////////////////
// 		   Public Methods		 //
///////////////////////////////////

// InitPartitionInfo create new PartitionInfo and return its pointer type
func InitPartitionInfo(topic string, partition int, leader *Node, replicas []*Node, inSyncReplicas []*Node, offlineReplicas []*Node) *PartitionInfo {
	return &PartitionInfo{
		topic:           topic,
		partition:       partition,
		leader:          leader,
		replicas:        replicas,
		inSyncReplicas:  inSyncReplicas,
		offlineReplicas: offlineReplicas,
	}
}

// GetTopic return the PartitionInfo topic
func (partInfo *PartitionInfo) GetTopic() string {
	return partInfo.topic
}

// GetPartition return the PartitionInfo partition
func (partInfo *PartitionInfo) GetPartition() int {
	return partInfo.partition
}

// GetLeader return the PartitionInfo leader
func (partInfo *PartitionInfo) GetLeader() *Node {
	return partInfo.leader
}

// GetReplicas return the PartitionInfo replicas
func (partInfo *PartitionInfo) GetReplicas() []*Node {
	return partInfo.replicas
}

// GetInSyncReplicas return the PartitionInfo inSyncReplicas
func (partInfo *PartitionInfo) GetInSyncReplicas() []*Node {
	return partInfo.inSyncReplicas
}

// GetOfflineReplicas return the PartitionInfo offlineReplicas
func (partInfo *PartitionInfo) GetOfflineReplicas() []*Node {
	return partInfo.offlineReplicas
}

// String return the string representatin of PartitionInfo
func (partInfo *PartitionInfo) String() string {
	leader := "none"
	if partInfo.leader != nil {
		leader = partInfo.leader.GetIDString()
	}

	return fmt.Sprintf("Partition(topic = %v, partition = %v, leader = %v, replicas = %v, isr = %v, offlineReplicas = %v)",
		partInfo.topic,
		partInfo.partition,
		leader,
		formatNodeIDs(partInfo.replicas),
		formatNodeIDs(partInfo.inSyncReplicas),
		formatNodeIDs(partInfo.offlineReplicas))
}

func formatNodeIDs(nodes []*Node) string {
	outputStr := "["

	if nodes != nil {
		nodesLength := len(nodes)
		for i, node := range nodes {
			outputStr += node.GetIDString()

			if i < nodesLength-1 {
				outputStr += ","
			}
		}
	}
	outputStr += "]"
	return outputStr
}
