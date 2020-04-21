package consumer

import (
	"context"
	"fmt"
	"strings"
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
	consumers   	[]Consumer
	assignments 	[]*consumepb.MetadataAssignment
	// key - topic, value - consumer
	topicConsumer map[string]*Consumer
	mux     		sync.RWMutex
	// key - topic, value - next partition idx to read from
	topicPartPoint 	map[string]int
}

// Consumer is a member of a consumer group
type Consumer struct {
	id      	int
	groupID 	int
	// assignments are given an idx (key)
	assignments map[int]*consumepb.MetadataAssignment
	brokerCon   map[int]clientpb.ClientService_ConsumeClient
	// key - brokerID, value - assignmentID
	brokerAssignmentMap map[int]int
	brokersAddr         map[int]string
}

// InitConsumer creates a consumer
func InitConsumer(id int, groupID int) *Consumer {
	return &Consumer{
		id:      id,
		groupID: groupID,
	}
}

// InitConsumerGroup creates a consumergroup and sets up broker connections
func InitConsumerGroup(id int, _topics string, brokerAddr string) *ConsumerGroup {
	topics := strings.Split(_topics, ";")
	cg := ConsumerGroup{id: id}
	fmt.Printf("Group %d initiated\n", id)
	const numConsumers = 2
	for i := 1; i <= numConsumers; i++ {
		cg.consumers = append(cg.consumers, InitConsumer(i, id))
		fmt.Printf("Consumer %d created under %d\n", i, id)
	}
	fmt.Printf("Done initialising group\n")

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


// Consume tries to consume information on the topic
func (cg *ConsumerGroup) Consume(topic string) {
	// TODO: Add/disregard comments based on your own intuition

	// check which consumers have partitions of that topic
	// a variable like `topicConsumer` might be useful. In CG struct but not constructed yet.

	cg.mux.Lock()
	// check which consumer should call consume now 
	/* say if consumer 1 has assignment T0-P1 and consumer 2 has assignment TO-P2
	// then make consumer 1 consume ->  consumer 2 consume -> consumer 1 consume so that the 
	// order of consumption of messages is the same as the order of production of messages */
	// lock is here because the initial idea is that there might be multiple consumer
	// threads consuming and changing the `topicPartPoint` value, remove if not required

	cg.mux.Unlock()

	// handle different cases of consume
	// 1) Normal consumption with no problem -> print msg
	
	// 2) Consume fails --> setup stream for next broker and call consume again


}

func (cg *ConsumerGroup) getAssignments(brokerAddr string, topics []string, maxWaitMs time.Duration) error {
	// connect to a broker
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	conn, err := grpc.Dial(addr, opts...)
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
			GroupID:   cg.id,
			TopicName: topic,
		}
		// send request to get assignments
		res, err := consumerClient.GetAssignment(ctx, req)
		if err != nil {
			statusErr, ok := status.FromError(err)
			if ok {
				if statusErr.Code() == codes.DeadlineExceeded {
					fmt.Printf("assignment not received for topic %v after %d ms", cg.topic, maxWaitMs)
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
		cg.assignments = append(cg.assignments, res)
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
}


// distributeAssignments to consumers in a round-robin manner
func (cg *ConsumerGroup) distributeAssignments(numConsumers int) {
	for i, assignment := range cg.assignments {
		idx := i % numConsumers
		cg.consumers[idx].assignments[i/numConsumers] = assignment
	}
}

func (c *Consumer) createBrokerAssignmentMap() {
	for k, v := range c.assignments {
		c.brokerAssignmentMap[v.GetBroker()] = append(c.brokerAssignmentMap[v.GetBroker()], k)
	}
}