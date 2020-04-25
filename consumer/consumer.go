package consumer

import (
	"context"
	"errors"
	"fmt"
	"log"
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
}

// Consumer is a member of a consumer group
type Consumer struct {
	id      int
	groupID int
	// assignments are given an idx (key)
	assignments map[int]*consumepb.MetadataAssignment
	brokerCon   map[int]clientpb.ClientService_ConsumeClient
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
	fmt.Println("Topics:",topics)

	err := cg.getAssignments(brokerAddr, topics, 100*time.Millisecond)
	if err != nil {
		panic(fmt.Sprintf("Consumer group %d unable to get assignments :%v\n", id, err))
	}

	// distribute assignments to consumers in a round-robin manner
	cg.distributeAssignments(numConsumers)

	for _, c := range cg.consumers {
		fmt.Printf("Consumer %v has the following assignments %v\n", c.id, c.assignments)
		c.createBrokersAddrMap()
		fmt.Sprintf("Consumer %v has created a broker-address map", c.id)
	}


	for _, c := range cg.consumers {
		fmt.Printf("BrokerAddrMap for Consumer %v, BrokerAddrMap: %v\n", c, c.brokersAddr)
		err := c.setupStreamToConsumeMsg()
		if err != nil {
			panic(fmt.Sprintf("Unable to set up stream to broker: %v\n", err))
		}
		fmt.Sprintf("Consumer %v has set up a broker-connection map\n", c.id)
	}

	cg.createTopicConsumerMap(_topics)
	fmt.Sprintf("Consumer Group %v has created a topic-consumers map", cg.id)

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
	fmt.Println("Consuming....")
	// TODO: Consider the event when there is a fault
	for i, topic := range cg.topics {
		if topic == _topic {
			// topic-consumer takes in a topic as key, and returns an ordered array of consumer, with the
			// ascending partitionIndex
			consumers := cg.topicConsumer[_topic]
			go doConsume(consumers, _topic)
			break
		} else if i == len(cg.topics)-1 {
			// already reached the end
			fmt.Println("Consumer group not set up to pull from this topic")
			return errors.New("Consumer group not set up to pull from this topic")
		}
	}
	return nil
}

func doConsume(sortedConsumers[]*Consumer, topic string) {
	fmt.Println("Running doConsume...")
	fmt.Println("Searching for topic", topic)
	fmt.Println("sortedCusumers", sortedConsumers)
	for _, c := range sortedConsumers {
		fmt.Printf("Consumer %v has assignment %v\n", c.id, c.assignments)
		for _, assignment := range c.assignments {
			// search for the correct assignment
			fmt.Printf("\nAssignment %v by Consumer %v\n", assignment, c.id)
			if assignment.GetTopicName() == topic {
				fmt.Println("Found topic name in assignment")
				var connectedBrokenID int
				stream := c.brokerCon[int(assignment.GetBroker())]
				err := stream.Send(&consumepb.ConsumeRequest{
					GroupID:              int32(c.groupID),
					ConsumerID:           int32(c.id),
					Partition:            assignment.GetPartitionIndex(),
					TopicName:            topic,
				})
				if err != nil {
					fmt.Println("Fail to connect. Retrying...")
					for i := 0; i < len(assignment.GetIsrBrokers()); i ++ {
						chosenBrokenID := assignment.GetIsrBrokers()[i].ID
						if chosenBrokenID == assignment.GetBroker(){
							// main broker is down based on above attempt
							continue
						}
						stream := c.brokerCon[int(chosenBrokenID)]
						err2 := stream.Send(&consumepb.ConsumeRequest{
							GroupID:              int32(c.groupID),
							ConsumerID:           int32(c.id),
							Partition:            assignment.GetPartitionIndex(),
							TopicName:            topic,
						})
						// When there is an error and it is the last broker it tries to connect
						if err2 != nil && i == len(assignment.GetIsrBrokers()) - 1 {
							log.Fatalf("ISRs exhausted. No more brokers to connect to. Error msg: %v\n", err)
						} else if err2 == nil {
							connectedBrokenID = int(chosenBrokenID)
							assignment.Broker = chosenBrokenID
							fmt.Println("Failed but recovered")
							fmt.Println("Connecting to Broker:", connectedBrokenID)
							break
						}
					}
				} else {
					connectedBrokenID = int(assignment.GetBroker())
					fmt.Println("Connecting to Broker:", int(assignment.GetBroker()))
				}
				fmt.Println("What is ")
				fmt.Println("BrokerCon Map:", c.brokerCon)
				fmt.Println("Connecting to BrokerID:", connectedBrokenID)
				res, err3 := stream.Recv()
				//fmt.Println("Response:", res)
				if err3 != nil {
					log.Fatalf("Error connecting")
				}
				responseHandler(int(assignment.GetBroker()), c.id, res, err3)
				//go func(connectedBrokenID int) {
				//	res, err3 := c.brokerCon[connectedBrokenID].Recv()
				//	fmt.Println("Response:", res)
				//	responseHandler(int(assignment.GetBroker()), c.id, res, err3)
				//}(connectedBrokenID)
			}
		}
	}
}

func responseHandler(brokerID int, consumerID int, res *consumepb.ConsumeResponse, err error) {
	fmt.Println("responseHandling...")
	if err != nil {
		log.Printf("Error when receiving messages from Broker %v\n", err)
	} else {
		recordSet := res.GetRecordSet()
		fmt.Println("RecordSet:", recordSet)
		records := recordSet.GetRecords()
		fmt.Println("Records:", records)
		for _, r := range records {
			bytes := r.GetValue()
			msg := string(bytes)
			fmt.Sprintf("Consumer %v has received message: %v\n", consumerID, msg)
		}
		log.Printf("Successfully receive the message from Broker %v\n", brokerID)
	}
}

func (cg *ConsumerGroup) getAssignments(brokerAddr string, topics []string, maxWaitMs time.Duration) error {
	fmt.Println("Connecting to", brokerAddr)
	// connect to a broker
	opts := []grpc.DialOption{grpc.WithInsecure()}
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
		// fmt.Println("assignments:", cg.assignments)
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
		cg.consumers[idx].assignments = make(map[int]*consumepb.MetadataAssignment)
		cg.consumers[idx].assignments[i/numConsumers] = assignment
	}
}

func (cg *ConsumerGroup) createTopicConsumerMap(topic string) {
	fmt.Println("Input topic:", topic)
	cg.topicConsumer = make(map[string][]*Consumer)
	for _, consumer := range cg.consumers {
		assignments := consumer.assignments
		for _, assignment := range assignments {
			if assignment.GetTopicName() == topic {
				cg.topicConsumer[topic] = append(cg.topicConsumer[assignment.GetTopicName()], consumer)
			}
		}
	}
	fmt.Println("topicConsumer Map before sort", cg.topicConsumer)
	sort.Slice(cg.topicConsumer[topic], func(i, j int) bool {
		if cg.topicConsumer[topic][i].assignments[0].GetPartitionIndex() <= cg.topicConsumer[topic][j].assignments[0].GetPartitionIndex() {
			return cg.topicConsumer[topic][i].assignments[0].GetPartitionIndex() < cg.topicConsumer[topic][j].assignments[0].GetPartitionIndex()
		} else {
			return cg.topicConsumer[topic][i].assignments[0].GetPartitionIndex() > cg.topicConsumer[topic][j].assignments[0].GetPartitionIndex()
		}
	})
	fmt.Println("topicConsumer Map after sort", cg.topicConsumer)

}

func (c *Consumer) createBrokersAddrMap() {
	c.brokersAddr = make(map[int]string)
	for _, assignment := range c.assignments {
		fmt.Printf("Assigment for Consumer %v is %v\n", c.id, c.assignments)
		// CANNOT FIND ISR Brokers
		for _, isr := range assignment.GetIsrBrokers() {
			c.brokersAddr[int(isr.GetID())] = fmt.Sprintf("%v:%v", isr.GetHost(), isr.GetPort())
		}
	}
	fmt.Println(c.brokersAddr)
}

func (c *Consumer) setupStreamToConsumeMsg() error {
	brokerCount := len(c.brokersAddr)
	brokersConnections := make(map[int]clientpb.ClientService_ConsumeClient)
	for brokerId, brokerAddr := range c.brokersAddr {
		opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
		conn, err := grpc.Dial(brokerAddr, opts...)
		if err != nil {
			brokerCount --
			continue
		}
		consumerClient := clientpb.NewClientServiceClient(conn)

		stream, err := consumerClient.Consume(context.Background())
		if err != nil {
			brokerCount --
			continue
		}
		brokersConnections[brokerId] = stream
	}
	c.brokerCon = brokersConnections
	if brokerCount == 0 {
		return errors.New("No brokers available")
	}
	return nil
}
