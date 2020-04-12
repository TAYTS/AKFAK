package consumer

type Consumer struct {
	topic     string
	partition int
	broker    int
}

func InitConsumer(topicName string, partitionIdx int, brokerIdx int) *Consumer {
	return &Consumer{
		topic:     topicName,
		partition: partitionIdx,
		broker:    brokerIdx,
	}
}

// Get requests
func (p *Consumer) GetTopic() string {
	return p.topic
}

func (p *Consumer) GetPartitionID() int {
	return p.partition
}

func (p *Consumer) GetBroker() int {
	return p.broker
}
