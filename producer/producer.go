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
	grpcConn         map[int]*grpc.ClientConn
	brokerCon        map[int]clientpb.ClientService_ProduceClient
	metadata         *metadatapb.MetadataResponse
	rr               *RoundRobinPartitioner
	inflightRequests map[int]*inflightRequest
	timers           map[int]*time.Timer
	mux              sync.RWMutex // used to ensure only one routine can send/modify a request to a broker
	metadataMux      sync.RWMutex
}

type inflightRequest struct {
	msgCount int
	req      *producepb.ProduceRequest
}

const METADATA_TIMEOUT = 100 * time.Millisecond

var errBrkNotAvailable = errors.New("Broker not available")

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
		grpcConn:         make(map[int]*grpc.ClientConn),
		brokerCon:        make(map[int]clientpb.ClientService_ProduceClient),
	}

	// get metadata and wait for 500ms
	err := p.waitOnMetadata(brokerAddr, METADATA_TIMEOUT)
	if err != nil {
		panic(fmt.Sprintf("Unable to get Topic Metadata: %v\n", err))
	}

	// setup the stream connections to all the required brokers
	p.setupStreamToSendMsg()

	// check if there are partitions available to send
	p.failIfNoAvailablePartition()

	return p
}

// Send used to send ProduceRequest to the broker
func (p *Producer) Send(message string) {
	p.metadataMux.RLock()
	// get partition idx for topic
	partIdx := p.getNextPartition()

	p.metadataMux.RUnlock()

	// create new record
	newRcd := recordpb.InitialiseRecordWithMsg(message)

	// pass the record to the request
	var produceReq *inflightRequest
	p.mux.RLock()
	produceReq, exist := p.inflightRequests[partIdx]
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
		p.inflightRequests[partIdx] = produceReq
		// set 500ms timeout for the sending the new request
		p.timers[partIdx] = time.NewTimer(500 * time.Millisecond)
		p.mux.Unlock()
		go p.doSend(partIdx)
	}
}

// CleanupResources used to cleanup the Producer resources
func (p *Producer) CleanupResources() {
	for _, conn := range p.brokerCon {
		conn.CloseSend()
	}
}

///////////////////////////////////
// 		    Private Methods		     //
///////////////////////////////////

// wait for cluster metadata including partitions for the given topic
// and partition (if specified, 0 if no preference) to be available
func (p *Producer) waitOnMetadata(brokerAddr string, maxWaitMs time.Duration) error {
	// dial one of the broker to get the topic metadata
	conn, err := grpc.Dial(brokerAddr, grpc.WithInsecure())
	defer conn.Close()
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
	protoCtx, protoCancel := context.WithTimeout(context.Background(), maxWaitMs)
	defer protoCancel()

	// send request to get metadata
	res, err := prdClient.WaitOnMetadata(protoCtx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Printf("Metadata not received for topic %v after %d ms\n", p.topic, maxWaitMs)
			} else {
				log.Printf("Unexpected error: %v\n", statusErr.Message())
			}
		} else {
			log.Printf("could not get metadata. error: %v\n", err)
		}

		return err
	}

	// save the metadata response to the producer
	p.metadataMux.Lock()
	p.metadata = res
	p.metadataMux.Unlock()

	return nil
}

// setupStreamToSendMsg setup the gRPC stream for sending message batch to the leader broker and attach to the producer instance
func (p *Producer) setupStreamToSendMsg() {
	// create a mapping of brokerID to broker info for easy access later
	brkMapping := make(map[int]*metadatapb.Broker)
	for _, brk := range p.metadata.GetBrokers() {
		brkMapping[int(brk.GetNodeID())] = brk
	}

	for _, part := range p.metadata.GetTopic().GetPartitions() {
		// try to find a broker for the partiton
		for {
			leaderID := p.getLeaderIDByPartition(int(part.GetPartitionIndex()))
			// if no broker available for the current partition move on to the next one
			if leaderID == -1 {
				break
			}
			leader := brkMapping[leaderID]

			if _, exist := p.brokerCon[leaderID]; !exist {
				// setup gRPC connection
				conn, err := grpc.Dial(fmt.Sprintf("%v:%v", leader.GetHost(), leader.GetPort()), grpc.WithInsecure())
				if err != nil {
					p.refreshMetadata()
					continue
				}
				// setup gRPC service
				prdClient := clientpb.NewClientServiceClient(conn)

				// setup stream
				stream, err := prdClient.Produce(context.Background())
				if err != nil {
					conn.Close()
					p.refreshMetadata()
					continue
				}

				p.grpcConn[leaderID] = conn
				p.brokerCon[leaderID] = stream
			}

			break
		}
	}
}

func (p *Producer) doSend(partIdx int) {
	for {
		p.mux.Lock()
		select {
		case <-p.timers[partIdx].C:
			brkID := p.getLeaderIDByPartition(partIdx)
			log.Println("Sending request to Broker", brkID)
			if conn, exist := p.brokerCon[brkID]; exist {
				conn.Send(p.inflightRequests[partIdx].req)
				// get the response
				_, err := p.brokerCon[brkID].Recv()
				p.postDoSendHook(partIdx, brkID, err)
			} else {
				p.postDoSendHook(partIdx, brkID, errBrkNotAvailable)
			}

			// remove the request
			delete(p.inflightRequests, partIdx)
			p.mux.Unlock()
			return
		default:
			// Send the request if the req has more than 15 messages
			if p.inflightRequests[partIdx].msgCount > 15 {
				brkID := p.getLeaderIDByPartition(partIdx)
				log.Println("Sending request to Broker", brkID)
				p.brokerCon[brkID].Send(p.inflightRequests[partIdx].req)

				// get the response
				_, err := p.brokerCon[brkID].Recv()
				p.postDoSendHook(partIdx, brkID, err)

				// remove the request
				delete(p.inflightRequests, partIdx)
				p.mux.Unlock()
				return
			}
			p.mux.Unlock()
		}
	}
}

// postDoSendHook is used after the doSend operation
// Used to handle broker failure and print send log
func (p *Producer) postDoSendHook(partIdx int, brokerID int, err error) {
	if err != nil {
		log.Printf("Detect Broker %v failure, retry to send message to other broker", brokerID)
		reqData := p.inflightRequests[partIdx].req
		newBrkID := brokerID

		if err == errBrkNotAvailable {
			newBrkID = -1
		}

		// try until the producer can send the message
		for {
			// clean up
			if newBrkID != -1 {
				p.brokerCon[newBrkID].CloseSend()
				p.grpcConn[newBrkID].Close()
				delete(p.brokerCon, newBrkID)
				delete(p.grpcConn, newBrkID)
			}

			// reset broker connection
			p.resetBrokerConnection()

			// get leader ID for a partition
			newBrkID = p.getLeaderIDByPartition(partIdx)

			// if the current partition does not have any brokers left
			// change the partition
			if newBrkID == -1 {
				partIdx = p.getNextPartition()
				newBrkID = p.getLeaderIDByPartition(partIdx)
				for _, topic := range reqData.GetTopicData() {
					topic.Partition = int32(partIdx)
				}
			}

			// send request
			err := p.brokerCon[newBrkID].Send(reqData)
			if err != nil {
				continue
			}

			// check response
			_, err = p.brokerCon[newBrkID].Recv()
			if err != nil {
				continue
			}

			break
		}
		log.Printf("Successfully send the request to Broker %v\n", newBrkID)
	} else {
		log.Printf("Successfully send the request to Broker %v\n", brokerID)
	}
}

// getAvailablePartitionID get all the available partition ID based on the alive brokers
func (p *Producer) getAvailablePartitionID() []*metadatapb.Partition {
	availablePartitionIDs := []*metadatapb.Partition{}

	for _, part := range p.metadata.GetTopic().GetPartitions() {
		leaderID := int(part.GetLeaderID())
		if leaderID == -1 {
			continue
		} else {
			availablePartitionIDs = append(availablePartitionIDs, part)
		}
	}
	return availablePartitionIDs
}

// getNextBroker find the next broker to send the message using Round Robin method and return the broker ID
func (p *Producer) getNextPartition() int {
	// get the next partition using Round Robin
	return p.rr.getPartition(p.topic, p.getAvailablePartitionID())
}

// getLeaderIDByPartition find the leader ID based on the given partition ID
// return -1 if no broker available for the given partition
func (p *Producer) getLeaderIDByPartition(partitionIdx int) int {
	for _, part := range p.metadata.GetTopic().GetPartitions() {
		if int(part.GetPartitionIndex()) == partitionIdx {
			return int(part.GetLeaderID())
		}
	}
	return -1
}

// refreshMetadata fetch new metadata and validate it
// Will fail the Producer if the new metadata does not have any partition available
func (p *Producer) refreshMetadata() {
	for _, brk := range p.metadata.GetBrokers() {
		// refresh the metadata
		err := p.waitOnMetadata(fmt.Sprintf("%v:%v", brk.GetHost(), brk.GetPort()), METADATA_TIMEOUT)
		if err != nil {
			continue
		}
		return
	}
	log.Fatalln("No brokers available")
}

// failIfNoAvailablePartition used to terminate the Producer if there is no partition available
func (p *Producer) failIfNoAvailablePartition() {
	if len(p.getAvailablePartitionID()) == 0 {
		log.Fatalln("No partition available")
	}
}

// resetBrokerConnection is used to refresh the metadata, setup stream connection and check the partition available
func (p *Producer) resetBrokerConnection() {
	p.refreshMetadata()

	p.setupStreamToSendMsg()

	p.failIfNoAvailablePartition()
}
