package consumer

import (
	"context"
	"fmt"
	"time"

	"AKFAK/proto/clientpb"
	"AKFAK/proto/consumermetadatapb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConsumerGroup struct {
	id               int
	ConsumerGroupCon map[string]clientpb.ClientService_ConsumeClient
	metadata         *consumermetadatapb.MetadataConsumerState
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

	//get replicas
	cg.getReplicas(brokersAddr, 100*time.Millisecond)
	// set up connection to all brokers. key - address : val - stream
	cg.dialConsumerGroup(brokersAddr)

	return &cg
}

func (cg *ConsumerGroup) dialConsumerGroup(brokersAddr map[int]string) {
	opts := grpc.WithInsecure()
	consumerGroupConnection := make(map[string]clientpb.ClientService_ConsumeClient)

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
				stream, err := cRPC.Consume(context.Background())
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

func (cg *ConsumerGroup) getReplicas(brokersAddr map[int]string, maxWaitMs time.Duration) {
	// dial one of the broker to get the topic metadata
	for _, addr := range brokersAddr {
		// connect to one of the broker
		opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
		conn, err := grpc.Dial(addr, opts...)
		if err != nil {
			// move to the next broker if fail to connect
			continue
		}

		// create gRPC client
		prdClient := clientpb.NewClientServiceClient(conn)

		// send get metadata request
		req := &consumermetadatapb.MetadataConsumerState{
			ConsumerGroups: cg.id,
		}

		// create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), maxWaitMs)

		// send request to get metadata
		res, err := prdClient.GetReplicas(ctx, req)
		if err != nil {
			statusErr, ok := status.FromError(err)
			if ok {
				if statusErr.Code() == codes.DeadlineExceeded {
					fmt.Printf("Metadata not received for topic %v after %d ms", cg.topic, maxWaitMs)
				} else {
					fmt.Printf("Unexpected error: %v", statusErr)
				}
			} else {
				fmt.Printf("could not get metadata. error: %v", err)
			}

			// move to the next broker if fail to connect
			cancel()
			conn.Close()
			continue
		}

		// save the metadata response to the consumer
		cg.metadata = res

		// clean up resources
		cancel()
		conn.Close()

		// stop contacting other brokers
		break
	}
}
