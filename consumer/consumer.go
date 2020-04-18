package consumer

import (
	"context"
	"fmt"
	"log"

	"AKFAK/proto/clientpb"

	"google.golang.org/grpc"
)

type ConsumerGroup struct {
	id               int
	ConsumerGroupCon map[string]clientpb.ClientService_MessageBatchClient
	clientService    clientpb.ClientServiceClient
}

type Consumer struct {
	id               int
	groupID          int
	assignedReplicas []AssignedReplica
}

type AssignedReplica struct {
	// multiple topic
	topic     string
	partition int
	broker    int
}

// InitConsumerGroup creates a consumergroup and sets up broker connections
func InitGroupConsumer(id int, brokersAddr map[int]string) *ConsumerGroup {
	fmt.Printf("Group %d initiated\n", id)
	NUM_CONSUMER := 3
	for i := 1; i <= NUM_CONSUMER; i++ {
		InitConsumer(i, id)
		fmt.Printf("Consumer %d created under %d\n", i, id)
	}
	fmt.Printf("==========Done initializing groups\n")

	cg := ConsumerGroup{id: id}
	// set up connection to all brokers. key - address : val - stream
	cg.dialConsumerGroup(brokersAddr)
	conn, err := grpc.Dial(brokersAddr[0], grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("could not connect to broker: %v", err)
	}
	cg.clientService = clientpb.NewClientServiceClient(conn)

	return &cg
}

func (cg *ConsumerGroup) dialConsumerGroup(brokersAddr map[int]string) {
	opts := grpc.WithInsecure()
	consumerGroupConnection := make(map[string]clientpb.ClientService_MessageBatchClient)

	for _, addr := range brokersAddr {
		fmt.Printf("Dialing %v\n", addr)
		waitc := make(chan struct{})
		go func() {
			for {
				cSock, err := grpc.Dial(addr, opts)
				if err != nil {
					fmt.Printf("Error did not connect: %v\n", err)
					continue
				}
				cRPC := clientpb.NewClientServiceClient(cSock)
				// Create stream by invoking the client
				stream, err := cRPC.MessageBatch(context.Background())
				if err != nil {
					fmt.Printf("Error while creating stream: %v\n", err)
					continue
				}
				consumerGroupConnection[addr] = stream
				break
			}
			close(waitc)
		}()
		<-waitc
	}

	cg.ConsumerGroupCon = consumerGroupConnection
}

func InitConsumer(id int, groupID int) *Consumer {
	return &Consumer{
		id:      id,
		groupID: groupID,
	}
}

//consumergroup call Getreplica and Replica response change the partition, broker
func (c *Consumer) Assignment(topicName string, partitionIdx int, brokerIdx int) {

	Assign := AssignedReplica{
		topic:     topicName,
		partition: partitionIdx,
		broker:    brokerIdx,
	}

	c.assignedReplicas = append(c.assignedReplicas, Assign)

	fmt.Println(stringPQ(c.assignedReplicas))
}

func stringPQ(pq []AssignedReplica) string {
	if len(pq) == 0 {
		return "Array EMPTY"
	}
	var ret string = "Array has"
	for _, msg := range pq {
		ret += fmt.Sprintf("[" + msg.topic + "]")
	}
	return ret
}
