package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	messagepb "AKFAK/proto/messagepb"
	metadatapb "AKFAK/proto/metadatapb"
	recordpb "AKFAK/proto/recordpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProducerRecord will be used to instantiate a record
type ProducerRecord struct {
	topic     string // topic the record will be appended to
	timestamp int64  // timestamp of the record in ms since epoch. If null, assign current time in ms
	value     []byte // the record contents
}

// getCurrentTimeinMs return current unix time in ms
func getCurrentTimeinMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// converts records from struct ProducerRecords to RecordBatch
func producerRecordsToRecordBatch(pRecords []ProducerRecord) *recordpb.RecordBatch {

	recordBatch := recordpb.InitialiseEmptyRecordBatch()
	recordBatch.FirstTimestamp = getCurrentTimeinMs()

	// convert []ProducerRecord to []recordpb.Record
	for _, pRecord := range pRecords {
		// assign current time in ms if record.timestamp is nil (0)
		if pRecord.timestamp == 0 {
			pRecord.timestamp = getCurrentTimeinMs()
		}
		// append record to record batch
		recordBatch.AppendRecord(&recordpb.Record{
			Length:         1,
			Attributes:     0,
			TimestampDelta: int32(getCurrentTimeinMs() - pRecord.timestamp), // current time-record.timestamp
			OffsetDelta:    1,
			ValueLen:       int32(len(pRecord.value)),
			Value:          pRecord.value,
		})
	}

	// timestamp when last record is added
	recordBatch.MaxTimestamp = getCurrentTimeinMs()
	recordBatch.Magic = 2
	recordBatch.LastOffsetDelta = 1 // still unsure about this

	return recordBatch
}

func dialBrokers(brokersAddr map[int]string) map[int]messagepb.MessageService_MessageBatchClient {
	opts := grpc.WithInsecure()
	time.Sleep(time.Second)
	brokersConnections := make(map[int]messagepb.MessageService_MessageBatchClient)

	for i, addr := range brokersAddr {
		fmt.Printf("Dialing %v\n", addr)
		waitc := make(chan struct{})
		go func() {
			for {
				cSock, err := grpc.Dial(addr, opts)
				if err != nil {
					fmt.Printf("Error did not connect: %v\n", err)
					continue
				}
				cRPC := messagepb.NewMessageServiceClient(cSock)
				// Create stream by invoking the client
				stream, err := cRPC.MessageBatch(context.Background())
				if err != nil {
					fmt.Printf("Error while creating stream: %v\n", err)
					continue
				}
				brokersConnections[i] = stream
				break
			}
			close(waitc)
		}()
		<-waitc
	}

	return brokersConnections
}

// Send takes records and sends to partition decided by round-robin
func Send(metadataProd metadatapb.MetadataServiceClient, records []ProducerRecord) {

	// all records in a recordbatch should have the same topic
	topic := records[0].topic

	// get brokers for topic
	brokersAddr := getBrokersForTopic(metadataProd, topic)

	// create connection to all brokers for that topic
	brokersConnections := dialBrokers(brokersAddr)

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
		for { // TODO: Remove while true loop, this is for demo purposes
			for _, req := range requests {
				fmt.Printf("Sending message: %v\n", msg)
				brokersConnections[partition].Send(req)
				time.Sleep(1000 * time.Millisecond)
			}
		}
		// brokersConnections[partition].CloseSend()
	}()

	// receive messages from broker
	go func() {
		for {
			res, err := brokersConnections[partition].Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	//block until everything is done
	<-waitc
}

// wait for cluster metadata including partitions for the given topic
// and partition (if specified, 0 if no preference) to be available
func waitOnMetadata(producer metadatapb.MetadataServiceClient, partition int, topicName string, maxWaitMs time.Duration) *metadatapb.MetadataResponse {

	req := &metadatapb.MetadataRequest{
		TopicName: topicName,
	}

	ctx, cancel := context.WithTimeout(context.Background(), maxWaitMs)
	defer cancel()

	res, err := producer.WaitOnMetadata(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Metadata not received for topics %v after %d ms", topicName, maxWaitMs)
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

func getBrokersForTopic(producer metadatapb.MetadataServiceClient, topic string) map[int]string {
	fmt.Println("in get brokers for topic")

	brokers := make(map[int]string)

	// get metadata on topic
	metadata := waitOnMetadata(producer, 0, topic, 3000*time.Millisecond)

	// get addresses of other brokers having partitions of topic in request
	for i, broker := range metadata.GetBrokers() {
		brokers[i] = broker.GetHost() + strconv.Itoa(int(broker.GetPort()))
	}

	fmt.Println(brokers)
	return brokers
}

func main() {
	// Set up a connection to a broker
	fmt.Println("running producer.go...")
	broker0 := "0.0.0.0:50051"
	conn, err := grpc.Dial(broker0, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// create a metadataclient
	metadataProd := metadatapb.NewMetadataServiceClient(conn)

	records := []ProducerRecord{
		ProducerRecord{
			topic:     "topic1",
			timestamp: getCurrentTimeinMs(),
			value:     []byte("this is the message content"),
		},
	}

	fmt.Println("sending message")

	Send(metadataProd, records)
}
