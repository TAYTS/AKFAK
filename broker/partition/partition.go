package partition

import (
	"fmt"
	"log"
	"os"
)

// Partition is used to describe per-partition state in the MetadataResponse
type Partition struct {
	topic           string
	partition       int
	leader          int
	replicas        []int
	inSyncReplicas  []int
	offlineReplicas []int
}

// InitPartition used to initialise the default value for the Partition type
func InitPartition(topicName string, partitionIdx int, leader int, replicas []int, inSyncReplicas []int, offlineReplicas []int) *Partition {
	return &Partition{
		topic:           topicName,
		partition:       partitionIdx,
		leader:          leader,
		replicas:        replicas,
		inSyncReplicas:  inSyncReplicas,
		offlineReplicas: offlineReplicas,
	}
}

// GetTopic return the topic name of the partition
func (p *Partition) GetTopic() string {
	return p.topic
}

// GetPartitionID return the partition ID of the partition
func (p *Partition) GetPartitionID() int {
	return p.partition
}

// GetLeader return the leader of the partition
func (p *Partition) GetLeader() int {
	return p.leader
}

// GetReplicas return the replicas list of the partition
func (p *Partition) GetReplicas() []int {
	return p.replicas
}

// GetInSyncReplicas return the inSyncReplicas list of the partition
func (p *Partition) GetInSyncReplicas() []int {
	return p.inSyncReplicas
}

// GetOfflineReplicas return the offlineReplicas list of the partition
func (p *Partition) GetOfflineReplicas() []int {
	return p.offlineReplicas
}

// CreatePartitionDir is a util func to create the local directory to store the partition info
func CreatePartitionDir(rootPath string, topicName string, partition int) error {
	filePath := fmt.Sprintf("%v/%v-%v", rootPath, topicName, partition)

	// Check if directory exist and remove if it exist
	_, err := os.Stat(filePath)
	if !os.IsNotExist(err) {
		os.RemoveAll(filePath)
	}

	// create the partition directory
	os.MkdirAll(filePath, os.ModePerm)

	// create the log file inside the directory
	fileName := fmt.Sprintf("%v.log", topicName)
	f, err := os.Create(fileName)
	if err != nil {
		log.Printf("Fail to create file: %v\n", err)
		return err
	}
	f.Close()

	return nil
}
