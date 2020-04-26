package cluster

import (
	"fmt"
	"hash/adler32"
)

// TopicPartition is a topic name and partition number
type TopicPartition struct {
	hash      int
	partition int
	topic     string
}

///////////////////////////////////
// 		   Public Methods		 //
///////////////////////////////////

// InitTopicPartition create new TopicPartition and return its pointer type
func InitTopicPartition(topic string, partition int) *TopicPartition {
	return &TopicPartition{
		partition: partition,
		topic:     topic,
	}
}

// GetPartition return the TopicPartition partition
func (tpPart *TopicPartition) GetPartition() int {
	return tpPart.partition
}

// GetTopic return the TopicPartition topic
func (tpPart *TopicPartition) GetTopic() string {
	return tpPart.topic
}

// HashCode return the TopicPartition hash
func (tpPart *TopicPartition) HashCode() int {
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

// String return the string representation of TopicPartition
func (tpPart *TopicPartition) String() string {
	return fmt.Sprintf("%v-%v", tpPart.topic, tpPart.partition)
}
