package partition

import (
	"fmt"
	"log"
	"os"
)

// CreatePartitionDir is a util func to create the local directory to store the partition info
func CreatePartitionDir(rootPath string, topicName string, partition int) error {
	filePath := fmt.Sprintf("%v/%v", rootPath, ConstructPartitionDirName(topicName, partition))

	// Check if directory exist and remove if it exist
	_, err := os.Stat(filePath)
	if !os.IsNotExist(err) {
		os.RemoveAll(filePath)
	}

	// create the partition directory
	os.MkdirAll(filePath, os.ModePerm)

	// create the log file inside the directory
	fileName := fmt.Sprintf("%v/%v", filePath, ContructPartitionLogName(topicName))
	f, err := os.Create(fileName)
	if err != nil {
		log.Printf("Fail to create file: %v\n", err)
		return err
	}
	f.Close()

	return nil
}

// ConstructPartitionDirName return the directory name for a given topic and partition
func ConstructPartitionDirName(topic string, partition int) string {
	return fmt.Sprintf("%v-%v", topic, partition)
}

// ContructPartitionLogName return the log filename given the topic
func ContructPartitionLogName(topic string) string {
	return fmt.Sprintf("%v.log", topic)
}
