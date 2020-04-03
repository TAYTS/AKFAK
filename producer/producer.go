package producer

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"AKFAK/proto/clientpb"
	"AKFAK/proto/messagepb"
	"AKFAK/proto/metadatapb"
	"AKFAK/proto/recordpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Producer represents a Kafka producer
type Producer struct {
	ID            int
	brokerCon     map[string]clientpb.ClientService_MessageBatchClient
	clientService clientpb.ClientServiceClient
}

// Record will be used to instantiate a record
type Record struct {
	Topic     string // topic the record will be appended to
	Timestamp int64  // timestamp of the record in ms since epoch. If null, assign current time in ms
	Value     []byte // the record contents
}

// InitProducer creates a producer and sets up broker connections
func InitProducer(id int, brokersAddr map[int]string) *Producer {
	p := Producer{ID: id}
	// set up connection to all brokers. key - address : val - stream
	p.dialBrokers(brokersAddr)
	conn, err := grpc.Dial(brokersAddr[0], grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("could not connect to broker: %v", err)
	}
	p.clientService = clientpb.NewClientServiceClient(conn)

	return &p
}

// GetCurrentTimeinMs return current unix time in ms
func GetCurrentTimeinMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// Send takes records and sends to partition decided by round-robin
func (p *Producer) Send(records []Record) {

	// all records in a recordbatch should have the same topic
	topic := records[0].Topic

	// get brokers for topic
	topicBrokers := p.getBrokersForTopic(topic)

	// TODO: obtain cluster from broker
	// cluster := cluster.getCluster()

	// compute partition to send to
	// partitioner := RoundRobinPartitioner{}
	// partition := partitioner.getPartition(topic, cluster)
	partition := 0

	recordBatch := producerRecordsToRecordBatch(records)
	msg := &messagepb.MessageBatchRequest{
		Records: recordBatch,
	}
	requests := []*messagepb.MessageBatchRequest{msg}

	waitc := make(chan struct{})

	// send messages to broker
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", msg)
			p.brokerCon[topicBrokers[partition]].Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		// p.brokerCon[topicBrokers[partition]].CloseSend()
	}()

	// receive messages from broker
	go func() {
		for {
			res, err := p.brokerCon[topicBrokers[partition]].Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res)
		}
		close(waitc)
	}()

	//block until everything is done
	// <-waitc
}

// converts records from struct Record to RecordBatch
func producerRecordsToRecordBatch(pRecords []Record) *recordpb.RecordBatch {

	recordBatch := recordpb.InitialiseEmptyRecordBatch()
	recordBatch.Magic = 2 // default

	// convert []Record to []recordpb.Record
	for i, pRecord := range pRecords {

		// assign current time in ms if record.Timestamp is nil (0)
		if pRecord.Timestamp == 0 {
			pRecord.Timestamp = GetCurrentTimeinMs()
		}

		recordBatch.FirstTimestamp = pRecord.Timestamp

		// append record to record batch
		recordBatch.AppendRecord(&recordpb.Record{
			Length:         1,
			Attributes:     0,
			TimestampDelta: int32(GetCurrentTimeinMs() - recordBatch.GetFirstTimestamp()),
			OffsetDelta:    int32(i),
			ValueLen:       int32(len(pRecord.Value)),
			Value:          pRecord.Value,
		})

		if i == len(pRecords)-1 {
			// timestamp when last record is added
			recordBatch.MaxTimestamp = pRecord.Timestamp
			// offset of last message in RecordBatch
			recordBatch.LastOffsetDelta = int32(i)
		}
	}
	return recordBatch
}

func (p *Producer) dialBrokers(brokersAddr map[int]string) {
	opts := grpc.WithInsecure()
	time.Sleep(time.Second)
	brokersConnections := make(map[string]clientpb.ClientService_MessageBatchClient)

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
				brokersConnections[addr] = stream
				break
			}
			close(waitc)
		}()
		<-waitc
	}

	p.brokerCon = brokersConnections
}

// wait for cluster metadata including partitions for the given topic
// and partition (if specified, 0 if no preference) to be available
func (p *Producer) waitOnMetadata(partition int, topicName string, maxWaitMs time.Duration) *metadatapb.MetadataResponse {

	req := &metadatapb.MetadataRequest{
		TopicName: topicName,
	}

	ctx, cancel := context.WithTimeout(context.Background(), maxWaitMs)
	defer cancel()

	res, err := p.clientService.WaitOnMetadata(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Printf("Metadata not received for topics %v after %d ms", topicName, maxWaitMs)
			} else {
				fmt.Printf("Unexpected error: %v", statusErr)
			}
		} else {
			log.Fatalf("could not get metadata. error: %v", err)
		}
	}

	log.Print("received metadata: ", res)

	return res
}

func (p *Producer) getBrokersForTopic(topic string) map[int]string {

	brokers := make(map[int]string)

	// get metadata on topic
	metadata := p.waitOnMetadata(0, topic, 3000*time.Millisecond)

	// get addresses of other brokers having partitions of topic in request
	for i, broker := range metadata.GetBrokers() {
		brokers[i] = fmt.Sprintf("%s:%d", broker.GetHost(), broker.GetPort())
	}

	fmt.Println(brokers)
	return brokers
}
