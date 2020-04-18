package producer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"AKFAK/proto/clientpb"
	"AKFAK/proto/metadatapb"
	"AKFAK/proto/producepb"
	"AKFAK/proto/recordpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Producer represents a Kafka producer
type Producer struct {
	ID               int
	topic            string
	brokerCon        map[int]clientpb.ClientService_ProduceClient
	metadata         *metadatapb.MetadataResponse
	rr               *RoundRobinPartitioner
	inflightRequests map[int]*inflightRequest
	timers           map[int]*time.Timer
	mux              sync.RWMutex // used to ensure only one routine can send/modify a request to a broker
}

type inflightRequest struct {
	msgCount int
	req      *producepb.ProduceRequest
}

///////////////////////////////////
// 		      Public Methods		   //
///////////////////////////////////

// InitProducer creates a producer and sets up broker connections
func InitProducer(id int, topic string, brokerAddr string) *Producer {
	// initialise the Producer instance
	p := &Producer{
		ID:               id,
		topic:            topic,
		rr:               InitRoundRobin(),
		inflightRequests: make(map[int]*inflightRequest),
		timers:           make(map[int]*time.Timer),
	}

	// get metadata and wait for 500ms
	err := p.waitOnMetadata(brokerAddr, 100*time.Millisecond)
	if err != nil {
		panic(fmt.Sprintf("Unable to get Topic Metadata: %v\n", err))
	}

	// setup the stream connections to all the required brokers
	err = p.setupStreamToSendMsg()
	if err != nil {
		panic(fmt.Sprintf("Unable to send message to broker: %v\n", err))
	}

	return p
}

// Send used to send ProduceRequest to the broker
func (p *Producer) Send(message string) {
	// get partition idx for topic
	partIdx := p.getNextPartition()

	// get broker ID for a partition
	brkID := p.getBrkIDByPartition(partIdx)

	// try to access another partition if unable to find broker ID for given partition
	for brkID == -1 {
		partIdx = p.getNextPartition()
		brkID = p.getBrkIDByPartition(partIdx)
	}

	// create new record
	newRcd := recordpb.InitialiseRecordWithMsg(message)

	// pass the record to the request
	var produceReq *inflightRequest
	p.mux.RLock()
	produceReq, exist := p.inflightRequests[brkID]
	if exist {
		// if there is inflight request, append the message to it
		produceReq.req.AddRecord(partIdx, newRcd)
		produceReq.msgCount++
		p.mux.RUnlock()
	} else {
		p.mux.RUnlock()
		// else create new ProduceRequest
		produceReq = &inflightRequest{
			msgCount: 1,
			req:      producepb.InitProduceRequest(p.topic, partIdx, newRcd),
		}
		// attach new request to the Producer instance
		p.mux.Lock()
		p.inflightRequests[brkID] = produceReq
		// set 500ms timeout for the sending the new request
		p.timers[brkID] = time.NewTimer(500 * time.Millisecond)
		p.mux.Unlock()
		go p.doSend(brkID)
	}
}

// CleanupResources used to cleanup the Producer resources
func (p *Producer) CleanupResources() {
	for _, conn := range p.brokerCon {
		conn.CloseSend()
	}
}

///////////////////////////////////
// 		   Private Methods		 //
///////////////////////////////////

// wait for cluster metadata including partitions for the given topic
// and partition (if specified, 0 if no preference) to be available
func (p *Producer) waitOnMetadata(brokerAddr string, maxWaitMs time.Duration) error {
	// dial one of the broker to get the topic metadata
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	conn, err := grpc.Dial(brokerAddr, opts...)
	if err != nil {
		return err
	}

	// create gRPC client
	prdClient := clientpb.NewClientServiceClient(conn)

	// send get metadata request
	req := &metadatapb.MetadataRequest{
		TopicName: p.topic,
	}

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), maxWaitMs)

	// send request to get metadata
	res, err := prdClient.WaitOnMetadata(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Printf("Metadata not received for topic %v after %d ms", p.topic, maxWaitMs)
			} else {
				fmt.Printf("Unexpected error: %v", statusErr)
			}
		} else {
			fmt.Printf("could not get metadata. error: %v", err)
		}

		// close the connection
		cancel()
		conn.Close()
		return err
	}

	// save the metadata response to the producer
	p.metadata = res

	// clean up resources
	cancel()
	conn.Close()

	return nil
}

// setupStreamToSendMsg setup the gRPC stream for sending message batch to the leader broker and attach to the producer instance
func (p *Producer) setupStreamToSendMsg() error {
	brkCount := len(p.metadata.GetTopic().GetPartitions())
	brokersConnections := make(map[int]clientpb.ClientService_ProduceClient)

	// create a mapping of brokerID to broker info for easy access later
	brkMapping := make(map[int32]*metadatapb.Broker)
	for _, brk := range p.metadata.GetBrokers() {
		brkMapping[brk.GetNodeID()] = brk
	}

	for _, part := range p.metadata.GetTopic().GetPartitions() {
		leaderID := part.GetLeaderID()
		leader := brkMapping[leaderID]

		// setup gRPC connection
		opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
		conn, err := grpc.Dial(fmt.Sprintf("%v:%v", leader.GetHost(), leader.GetPort()), opts...)
		if err != nil {
			brkCount--
			// try the next broker
			continue
		}

		// setup gRPC service
		prdClient := clientpb.NewClientServiceClient(conn)

		stream, err := prdClient.Produce(context.Background())
		if err != nil {
			brkCount--
			// try the next broker
			continue
		}
		brokersConnections[int(leaderID)] = stream
	}

	// return error if no brokers available
	if brkCount == 0 {
		return errors.New("No brokers available")
	}

	// attach the broker connections and locks mapping to the Producer instance
	p.brokerCon = brokersConnections

	return nil
}

func (p *Producer) doSend(brokerID int) {
	for {
		p.mux.Lock()
		select {
		case <-p.timers[brokerID].C:
			fmt.Printf("Sending request to Broker %v\n", brokerID)
			p.brokerCon[brokerID].Send(p.inflightRequests[brokerID].req)

			// get the response
			res, err := p.brokerCon[brokerID].Recv()
			responseHandler(brokerID, res, err)

			// remove the request
			delete(p.inflightRequests, brokerID)
			p.mux.Unlock()
			return
		default:
			// Send the request if the req has more than 15 messages
			if p.inflightRequests[brokerID].msgCount > 15 {
				p.brokerCon[brokerID].Send(p.inflightRequests[brokerID].req)

				// get the response
				res, err := p.brokerCon[brokerID].Recv()
				responseHandler(brokerID, res, err)

				// remove the request
				delete(p.inflightRequests, brokerID)
				p.mux.Unlock()
				return
			}
			p.mux.Unlock()
		}
	}
}

func responseHandler(brokerID int, res *producepb.ProduceResponse, err error) {
	if err != nil {
		// TODO: Retry sending the request? Remove broker from the connection? Ignore fail request?
		log.Printf("Error when sending messages to Broker %v\n", err)
	} else {
		log.Printf("Successfully send the request to Broker %v\n", brokerID)
	}
}

// getAvailablePartition get all the available partition based on the alive brokers
func (p *Producer) getAvailablePartition() []*metadatapb.Partition {
	availablePartition := []*metadatapb.Partition{}

	for _, part := range p.metadata.GetTopic().GetPartitions() {
		leaderID := int(part.GetLeaderID())
		if _, exist := p.brokerCon[leaderID]; exist {
			availablePartition = append(availablePartition, part)
		}
	}
	return availablePartition
}

// getNextBroker find the next broker to send the message using Round Robin method and return the broker ID
func (p *Producer) getNextPartition() int {
	// get the next partition using Round Robin
	return p.rr.getPartition(p.topic, p.getAvailablePartition())

}

func (p *Producer) getBrkIDByPartition(partitionIdx int) int {
	for _, part := range p.metadata.GetTopic().GetPartitions() {
		if int(part.GetPartitionIndex()) == partitionIdx {
			return int(part.GetLeaderID())
		}
	}
	fmt.Println("could not find broker id corresponding to partition", partitionIdx)
	return -1
}
