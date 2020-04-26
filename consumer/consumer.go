package consumer

import (
	"AKFAK/proto/metadatapb"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"AKFAK/proto/clientpb"
	"AKFAK/proto/consumepb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const METADATA_TIMEOUT = 100 * time.Millisecond

// Consumer is a member of a consumer group
type Consumer struct {
	ID int
	// assignments are given an idx (key)
	//brokerAddr    string
	topic         string
	brokerAddrMap map[int]string
	metadata      *metadatapb.MetadataResponse
	brokerCon     map[int]clientpb.ClientService_ConsumeClient
	grpcConn      map[int]*grpc.ClientConn
	timers        map[int]*time.Timer
	mux           sync.RWMutex // used to ensure only one routine can send/modify a request to a broker
	metadataMux   sync.RWMutex
	partitionIdx  int
	offset        int // it will not remember if it switches to read a new topic and read back the old topic
}

// InitConsumerGroup creates a consumergroup and sets up broker connections
func InitConsumer(id int, topic string, partitionIdx int, brokerAddr string) *Consumer {

	// Dial to broker to get metadata
	c := &Consumer{
		ID: id,
		topic: topic,
		brokerCon: make(map[int]clientpb.ClientService_ConsumeClient),
		grpcConn:	make(map[int]*grpc.ClientConn),
		partitionIdx: partitionIdx,
	}
	// get metadata and wait for 500ms
	err := c.waitOnMetadata(brokerAddr, METADATA_TIMEOUT)
	if err != nil {
		panic(fmt.Sprintf("Unable to get Topic Metadata: %v\n", err))
	}

	// setup the stream connections to all the required brokers
	c.setupStreamToSendMsg()

	// check if there are partitions available to consume
	c.failIfNoAvailablePartition()

	return c
}

func (c *Consumer) waitOnMetadata(brokerAddr string, maxWaitMs time.Duration) error {
	// find a broker to get metadata

	conn, err := grpc.Dial(brokerAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	// create gRPC client
	conClient := clientpb.NewClientServiceClient(conn)

	// send get metadata request
	req := &metadatapb.MetadataRequest{
		TopicName: c.topic,
	}

	// create context with timeout
	protoCtx, protoCancel := context.WithTimeout(context.Background(), maxWaitMs)
	defer protoCancel()

	// send request to get metadata
	res, err := conClient.WaitOnMetadata(protoCtx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Printf("Metadata not received for topic %v after %d ms\n", c.topic, maxWaitMs)
			} else {
				log.Printf("Unexpected error: %v\n", statusErr.Message())
			}
		} else {
			log.Printf("could not get metadata. error: %v\n", err)
		}
		return err
	}
	// save the metadata response to the consumer
	c.metadataMux.Lock()
	c.metadata = res
	c.metadataMux.Unlock()

	return nil
}

func (c *Consumer) Consume() {
	c.metadataMux.RLock()
	// get partition idx for topic
	brkID := c.getLeaderIDByPartition(c.partitionIdx)
	c.metadataMux.RUnlock()

	c.mux.Lock()
	c.timers[brkID] = time.NewTimer(500 * time.Millisecond)
	c.mux.Unlock()

	go c.doConsume(brkID, c.partitionIdx)

}

func (c *Consumer)doConsume(brokerID int, partitionIdx int) {
	for {
		//c.mux.Lock()
		// send to broker for consumereq
		// get res
		// if res == nil
		// sleep

		// refresh new connection
		// change new connection
		// if leaderId != -1 // no brokers from that partition.
		select {
		case <-c.timers[brokerID].C:
			log.Println("Sending request to Broker", brokerID)
			if conn, exist := c.brokerCon[brokerID]; exist {
				req := consumepb.ConsumeRequest{
					ConsumerID:           int32(c.ID),
					Partition:            int32(partitionIdx),
					TopicName:            c.topic,
				}
				conn.Send(&req)
				res, err := conn.Recv()
				// get the response
				c.postDoConsumeHook(brokerID, res, err)
			} else {
				log.Fatalln("Broker not available")
			}
		}
	}
}

func (c *Consumer) postDoConsumeHook(brokerID int, res *consumepb.ConsumeResponse, err error) {
	if err != nil {
		log.Printf("Detect Broker %v failure, retry to send message to other broker", brokerID)
		newBrkID := brokerID
		if err == errors.New("Broker not available") {
			newBrkID = -1
		}
		// try until all the partitions are exhausted

	} else {
		records := res.RecordSet.GetRecords()
		for _, record := range records {
			bytes := record.GetValue()
			msg := string(bytes)
			log.Printf("Consumer %v has received message: %v\n", c.ID, msg)
		}
	}
}

// CleanupResources used to cleanup the Producer resources
func (c *Consumer) CleanupResources() {
	for _, conn := range c.brokerCon {
		conn.CloseSend()
	}
}

// setupStreamToSendMsg setup the gRPC stream for sending message batch to the leader broker and attach to the producer instance
func (c *Consumer) setupStreamToSendMsg() {
	// create a mapping of brokerID to broker info for easy access later
	brkMapping := make(map[int]*metadatapb.Broker)
	for _, brk := range c.metadata.GetBrokers() {
		brkMapping[int(brk.GetNodeID())] = brk
	}

	for _, part := range c.metadata.GetTopic().GetPartitions() {
		// try to find a broker for the partiton
		for {
			leaderID := c.getLeaderIDByPartition(int(part.GetPartitionIndex()))
			// if no broker available for the current partition move on to the next one
			if leaderID == -1 {
				break
			}
			leader := brkMapping[leaderID]

			if _, exist := c.brokerCon[leaderID]; !exist {
				// setup gRPC connection
				conn, err := grpc.Dial(fmt.Sprintf("%v:%v", leader.GetHost(), leader.GetPort()), grpc.WithInsecure())
				if err != nil {
					c.refreshMetadata()
					continue
				}
				// setup gRPC service
				conClient := clientpb.NewClientServiceClient(conn)

				// setup stream
				stream, err := conClient.Consume(context.Background())
				if err != nil {
					conn.Close()
					c.refreshMetadata()
					continue
				}

				c.grpcConn[leaderID] = conn
				c.brokerCon[leaderID] = stream
			}
			break
		}
	}
}

// getLeaderIDByPartition find the leader ID based on the given partition ID
// return -1 if no broker available for the given partition
func (c *Consumer) getLeaderIDByPartition(partitionIdx int) int {
	for _, part := range c.metadata.GetTopic().GetPartitions() {
		if int(part.GetPartitionIndex()) == partitionIdx {
			return int(part.GetLeaderID())
		}
	}
	return -1
}

// refreshMetadata fetch new metadata and validate it
// Will fail the Producer if the new metadata does not have any partition available
func (c *Consumer) refreshMetadata() {
	for _, brk := range c.metadata.GetBrokers() {
		// refresh the metadata
		err := c.waitOnMetadata(fmt.Sprintf("%v:%v", brk.GetHost(), brk.GetPort()), METADATA_TIMEOUT)
		if err != nil {
			continue
		}
		return
	}
	log.Fatalln("No brokers available")
}


// failIfNoAvailablePartition used to terminate the Producer if there is no partition available
func (c *Consumer) failIfNoAvailablePartition() {
	if len(c.getAvailablePartitionID()) == 0 {
		log.Fatalln("No partition available")
	}
}

// getAvailablePartitionID get all the available partition ID based on the alive brokers
func (c *Consumer) getAvailablePartitionID() []*metadatapb.Partition {
	availablePartitionIDs := []*metadatapb.Partition{}

	for _, part := range c.metadata.GetTopic().GetPartitions() {
		leaderID := int(part.GetLeaderID())
		if leaderID == -1 {
			continue
		} else {
			availablePartitionIDs = append(availablePartitionIDs, part)
		}
	}
	return availablePartitionIDs
}
