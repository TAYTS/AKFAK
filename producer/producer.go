package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"
	
	messagepb "AKFAK/proto/messagepb"
	metadatapb "AKFAK/proto/metadatapb"
	recordpb "AKFAK/proto/recordpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

// ProducerRecord will be used to instantiate a record
type ProducerRecord struct {
	topic     string	// topic the record will be appended to
	partition int		// partition to which the record should be sent
	timestamp int64		// timestamp of the record in ms since epoch. If null, assign current time in ms
	key       []byte		// the key that will be included in the record
	value     []byte		// the record contents
}


// getCurrentTimeinMs return current unix time in ms
func getCurrentTimeinMs() int64{
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func producerRecordsToRecordBatch(pRecords []ProducerRecord) *recordpb.RecordBatch {

	recordBatch := recordpb.InitialiseEmptyRecordBatch() 
	recordBatch.FirstTimestamp = getCurrentTimeinMs()

	// convert []ProducerRecord to []recordpb.Record
	for _, pRecord := range pRecords {
		// assign current time in ms if record.timestamp is nil (0)
		if pRecord.timestamp == 0 {
			timestamp := getCurrentTimeinMs()
		} else {
			timestamp := pRecord.timestamp
		}
		// append record to record batch
		recordBatch.AppendRecord(&recordpb.Record{
						Length:         1,
						Attributes:     0,
						TimestampDelta: int32(getCurrentTimeinMs()-pRecord.timestamp), // current time-record.timestamp
						OffsetDelta:    1,
						KeyLength:      int32(len(pRecord.key)),
						Key:            pRecord.key,
						ValueLen:       int32(len(pRecord.value)),
						Value:          pRecord.value,
					})
	}
	
	// timestamp when last record is added
	recordBatch.MaxTimestamp = getCurrentTimeinMs()
	recordBatch.Magic = 2
	recordBatch.LastOffsetDelta = 1 // still unsure about this
	recordBatch.ProducerId = 1 // still unsure about this
	
	return recordBatch
}

// Send sends message to different partitions in a round-robin fashion
func Send(msgProducer messagepb.MessageServiceClient, metadataProd metadatapb.MetadataServiceClient, records []ProducerRecord) {

	// cluster := cluster.getCluster() // TODO: get cluster from broker
	partitioner := RoundRobinPartitioner{}
	partition := partitioner.getPartition(record.topic, cluster) // TODO: realised nowhere to indicate partition in record
		
	// Create stream by invoking the client
	stream, err := msgProducer.MessageBatch(context.Background())
	
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	recordBatch := producerRecordsToRecordBatch(records)
	msg := &messagepb.MessageBatchRequest{
		Records:	recordBatch,
	}

	waitc := make(chan struct{})

	// send messages to broker
	go func() {
		fmt.Printf("Sending message: %v\n", msg)
		stream.Send(msg)
		time.Sleep(1000 * time.Millisecond)
		
		stream.CloseSend()
	}()

	// receive messages from broker
	go func() {
		for {
			res, err := stream.Recv()
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
func waitOnMetadata(producer metadatapb.MetadataServiceClient, partition int, topicNames []string, maxWaitMs time.Duration) *metadatapb.MetadataResponse {

	req := &metadatapb.MetadataRequest{
		TopicNames: topicNames,
	}

	ctx, cancel := context.WithTimeout(context.Background(), maxWaitMs)
	defer cancel()

	res, err := producer.WaitOnMetadata(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Metadata not received for topics %v after %d ms", topicNames, maxWaitMs)
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

func main() {
	// Set up a connection to the chosen broker chosen via round-robin
	address := "0.0.0.0:50051"
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// msgProd := messagepb.NewMessageServiceClient(conn)
	metadataProd := metadatapb.NewMetadataServiceClient(conn)

	// === test functions ===
	// Send(prod)
	waitOnMetadata(metadataProd, 0, []string{"topic1", "topic2"}, 3000*time.Millisecond)
}
