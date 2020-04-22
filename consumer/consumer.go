package consumer

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"AKFAK/proto/clientpb"
	"AKFAK/proto/consumepb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ConsumerGroup holds consumers
type ConsumerGroup struct {
	id    			int
	topics 			[]string
	consumers   	[]*Consumer
	assignments 	[]*consumepb.MetadataAssignment
	// key - topic, value - consumer
	topicConsumer map[string][]*Consumer
	mux           sync.RWMutex
	// key - topic, value - next partition idx to read from
	topicPartPoint map[string]int
}

// Consumer is a member of a consumer group
type Consumer struct {
	id      int
	groupID int
	// assignments are given an idx (key)
	assignments map[int]*consumepb.MetadataAssignment
	brokerCon   map[int]clientpb.ClientService_ConsumeClient
	// key - brokerID, value - assignmentIDs
	brokerAssignmentMap map[int][]int
	brokersAddr         map[int]string
}

// InitConsumerGroup creates a consumergroup and sets up broker connections
func InitConsumerGroup(id int, _topics string, brokerAddr string) *ConsumerGroup {
	topics := strings.Split(_topics, ",")
	cg := ConsumerGroup{
		id: id,
		topics: topics,
	}
	fmt.Printf("Group %d initiated\n", id)
	const numConsumers = 2
	for i := 1; i <= numConsumers; i++ {
		cg.consumers = append(cg.consumers, initConsumer(i, id))
		fmt.Printf("Consumer %d created under Group %d\n", i, id)
	}
	fmt.Println("Done initialising group")
	fmt.Println(topics)

	err := cg.getAssignments(brokerAddr, topics, 100*time.Millisecond)
	if err != nil {
		panic(fmt.Sprintf("Consumer group %d unable to get assignments :%v\n", id, err))
	}

	// distribute assignments to consumers in a round-robin manner
	cg.distributeAssignments(numConsumers)

	// set the first partition idx to be read from = 0
	for _, topic := range topics {
		cg.topicPartPoint[topic] = 0
	}

	return &cg
}

// initConsumer creates a consumer
func initConsumer(id int, groupID int) *Consumer {
	return &Consumer{
		id:      id,
		groupID: groupID,
	}
}

// Consume tries to consume information on the topic
func (cg *ConsumerGroup) Consume(_topic string) error {
	for i, topic := range cg.topics {
		if topic == _topic {
			break
		} else if i == len(cg.topics)-1 {
			// already reached the end
			fmt.Println("Consumer group not set up to pull from this topic")
		}
	}
	// TODO: Add/disregard comments based on your own intuition
	// sorted by partitionIndex
	cg.createTopicConsumerMap(_topic)

	// key - topic, value - []*Consumer
	// can assume one assignment for the same topic for each consumer
	// and the consumerlist is sorted based on partitionId
	for _, consumer := range cg.topicConsumer[_topic] {
		consumer.assignments[0].GetPartitionIndex()
		// Connect to broker
		// Get the message
		//
	}

	// check which consumers have partitions of that topic
	// a variable like `topicConsumer` might be useful. In CG struct but not constructed yet.
	// check which consumer should call consume now
	consumers := cg.consumers
	for _, consumer := range consumers {
		assignments := consumer.assignments
		for _, assignment := range assignments {
			if assignment.GetTopicName() == _topic {
				consumerList := cg.topicConsumer[assignment.GetTopicName()]
				consumerList = append(consumerList, consumer)
			}
		}
	}
	return nil

	//cg.mux.Lock()
	// check which consumer should call consume now 
	/* say if consumer 1 has assignment T0-P1 and consumer 2 has assignment TO-P2
	// then make consumer 1 consume ->  consumer 2 consume -> consumer 1 consume so that the
	// order of consumption of messages is the same as the order of production of messages */
	// lock is here because the initial idea is that there might be multiple consumer
	// threads consuming and changing the `topicPartPoint` value, remove if not required

	//cg.mux.Unlock()

	// handle different cases of consume
	// 1) Normal consumption with no problem -> print msg

	// 2) Consume fails --> setup stream for next broker and call consume again
}

func (cg *ConsumerGroup) doSend(brokerID int) {
	for {
		select {

		}
	}
}

func (cg *ConsumerGroup) getAssignments(brokerAddr string, topics []string, maxWaitMs time.Duration) error {
	// connect to a broker
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	conn, err := grpc.Dial(brokerAddr, opts...)
	if err != nil {
		return err
	}

	// create gRPC client
	consumerClient := clientpb.NewClientServiceClient(conn)

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), maxWaitMs)

	for _, topic := range topics {
		// create a request for a topic
		req := &consumepb.GetAssignmentRequest{
			GroupID:   int32(cg.id),
			TopicName: topic,
		}
		// send request to get assignments
		res, err := consumerClient.GetAssignment(ctx, req)
		if err != nil {
			statusErr, ok := status.FromError(err)
			if ok {
				if statusErr.Code() == codes.DeadlineExceeded {
					fmt.Printf("assignment not received for topic %v after %d ms", topic, maxWaitMs)
				} else {
					fmt.Printf("unexpected error: %v", statusErr)
				}
			} else {
				fmt.Printf("could not get assignment. error: %v", err)
			}

			// close the connection
			cancel()
			conn.Close()
			return err
		}
		// save the assignment to the consumer group
		cg.assignments = append(cg.assignments, res.GetAssignments()...)
	}

	// clean up resources
	cancel()
	conn.Close()

	return nil
}

func (c *Consumer) getBrokerIdxAddrForAssignment(assignmentIdx int) (int, string) {
	brokerID := c.assignments[assignmentIdx].GetBroker()
	for i, isrbroker := range c.assignments[assignmentIdx].GetIsrBrokers() {
		if isrbroker.GetID() == brokerID {
			return i, fmt.Sprintf("%v:%v", isrbroker.GetHost(), isrbroker.GetPort())
		} else {
			panic("No matching broker ID in isrbrokers")
		}
	}
	return -1, ""
}

// distributeAssignments to consumers in a round-robin manner
func (cg *ConsumerGroup) distributeAssignments(numConsumers int) {
	for i, assignment := range cg.assignments {
		idx := i % numConsumers
		cg.consumers[idx].assignments[i/numConsumers] = assignment
	}
}

func (cg *ConsumerGroup) createTopicConsumerMap(topic string) {
	for _, consumer := range cg.topicConsumer[topic] {
		assignments := consumer.assignments
		for _, assignment := range assignments {
			if assignment.GetTopicName() == topic {
				cg.topicConsumer[topic] = append(cg.topicConsumer[assignment.GetTopicName()], consumer)
			}
		}
	}
	sort.Slice(cg.topicConsumer[topic], func(i, j int) bool {
		if cg.topicConsumer[topic][i].assignments[0].GetPartitionIndex() <= cg.topicConsumer[topic][j].assignments[0].GetPartitionIndex() {
			return cg.topicConsumer[topic][i].assignments[0].GetPartitionIndex() < cg.topicConsumer[topic][j].assignments[0].GetPartitionIndex()
		} else {
			return cg.topicConsumer[topic][i].assignments[0].GetPartitionIndex() > cg.topicConsumer[topic][j].assignments[0].GetPartitionIndex()
		}
	})
}

func (c *Consumer) createBrokerAssignmentMap() {
	for k, v := range c.assignments {
		c.brokerAssignmentMap[int(v.GetBroker())] = append(c.brokerAssignmentMap[int(v.GetBroker())], k)
	}
}
