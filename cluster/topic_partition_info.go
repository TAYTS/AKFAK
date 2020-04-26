package cluster

// TopicPartitionInfo containing leadership, replicas and ISR information for a topic partition.
type TopicPartitionInfo struct {
	partition int
	leader    *Node
	replicas  []*Node
	isr       []*Node
}

///////////////////////////////////
// 		   Public Methods		 //
///////////////////////////////////

// InitTopicPartitionInfo create new TopicPartitionInfo and return its pointer type
func InitTopicPartitionInfo(partition int, leader *Node, replicas []*Node, isr []*Node) *TopicPartitionInfo {
	return &TopicPartitionInfo{
		partition: partition,
		leader:    leader,
		replicas:  replicas,
		isr:       isr,
	}
}

// GetPartition return the TopicPartitionInfo partition
func (tpPart *TopicPartitionInfo) GetPartition() int {
	return tpPart.partition
}

// GetLeader return the TopicPartitionInfo leader
func (tpPart *TopicPartitionInfo) GetLeader() *Node {
	return tpPart.leader
}

// GetReplicas return the TopicPartitionInfo replicas
func (tpPart *TopicPartitionInfo) GetReplicas() []*Node {
	return tpPart.replicas
}

// GetISR return the TopicPartitionInfo isr
func (tpPart *TopicPartitionInfo) GetISR() []*Node {
	return tpPart.isr
}
