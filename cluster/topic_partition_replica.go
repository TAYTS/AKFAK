package cluster

import (
	"fmt"
	"hash/adler32"
)

// TopicPartitionReplica is the topic name, partition number and the brokerId of the replica
type TopicPartitionReplica struct {
	hash      int
	brokerID  int
	partition int
	topic     string
}

///////////////////////////////////
// 		   Public Methods		 //
///////////////////////////////////

// InitTopicPartitionReplica create new TopicPartitionReplica and return its pointer type
func InitTopicPartitionReplica(topic string, partition int, brokerID int) *TopicPartitionReplica {
	return &TopicPartitionReplica{
		hash:      0,
		brokerID:  brokerID,
		topic:     topic,
		partition: partition,
	}
}

// GetTopic return the TopicPartitionReplica topic
func (tpPart *TopicPartitionReplica) GetTopic() string {
	return tpPart.topic
}

// GetPartition return the TopicPartitionReplica partition
func (tpPart *TopicPartitionReplica) GetPartition() int {
	return tpPart.partition
}

// GetBrokerID return the TopicPartitionReplica brokerID
func (tpPart *TopicPartitionReplica) GetBrokerID() int {
	return tpPart.brokerID
}

// HashCode return the TopicPartitionReplica hash
func (tpPart *TopicPartitionReplica) HashCode() int {
	if tpPart.hash != 0 {
		return tpPart.hash
	}

	const prime = 31
	result := 1
	result = prime*result + tpPart.partition
	result = prime*result + int(adler32.Checksum([]byte(tpPart.topic)))
	tpPart.hash = result

	return result
}

// String return the string representation of TopicPartitionReplica
func (tpPart *TopicPartitionReplica) String() string {
	return fmt.Sprintf("%v-%v-%v", tpPart.topic, tpPart.partition, tpPart.brokerID)
}
