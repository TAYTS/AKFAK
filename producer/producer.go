package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"
	"strconv"

	messagepb "AKFAK/proto/messagepb"
	metadatapb "AKFAK/proto/metadatapb"
	recordpb "AKFAK/proto/recordpb"
	"google.golang.org/grpc"
)

// ProducerRecord will be used to instantiate a record
type ProducerRecord struct {
	topic     string
	partition int32
	timestamp float32
	key       byte
	value     byte
	headers   recordpb.Header
}

// round-robin returns an integer to tell the producer which partition to send to
func decidePartition() int {
	return 50051
}

// implementation modelled after doSend(...) in Kafka's Java implementation of KafkaProducer
// send sends message to different partitions in a round-robin fashion
func send(producer messagepb.MessageServiceClient) {

	// make sure metadata for topic is available
	// serialize fields in record
	// fill in fields in record
	// send to broker

	fmt.Println("Starting to send message...")

	// Create stream by invoking the client
	stream, err := producer.MessageBatch(context.Background())

	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	// dummy msg
	msgs := []*messagepb.MessageBatchRequest{
		{
			Records: &recordpb.RecordBatch{
				BaseOffset:           1,
				BatchLength:          []byte{byte('a')},
				PartitionLeaderEpoch: 1,
				Magic:                1,
				Crc:                  1,
				Attributes:           1,
				LastOffsetDelta:      1,
				FirstTimestamp:       1,
				MaxTimestamp:         1,
				ProducerId:           1,
				ProducerEpoch:        1,
				BaseSequence:         1,
				Records: []*recordpb.Record{
					&recordpb.Record{
						Length:         1,
						Attributes:     1,
						TimestampDelta: 1,
						OffsetDelta:    1,
						KeyLength:      1,
						Key:            []byte{byte('a')},
						ValueLen:       1,
						Value:          []byte{byte('a')},
						Headers: []*recordpb.Header{
							&recordpb.Header{
								HeaderKeyLength:   1,
								HeaderKey:         "a",
								HeaderValueLength: 1,
								Value:             []byte{byte('a')},
							},
						},
					},
				},
			},
		},
	}

	waitc := make(chan struct{})

	// send messages to broker
	go func() {
		for _, msg := range msgs {
			fmt.Printf("Sending message: %v\n", msg)
			stream.Send(msg)
			time.Sleep(1000 * time.Millisecond)
		}
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

// wait for cluster metadata including partitions for the given topic to be available
func getMetadata(producer metadatapb.MetadataServiceClient, topicNames []string) *metadatapb.MetadataResponse {

	req := &metadatapb.MetadataRequest{
		TopicNames: topicNames,
	}

	res, err := producer.GetMetadata(context.Background(), req)
	if err != nil {
		log.Fatalf("could not get metadata. error: %v", err)
	}

	log.Print("received metadata: ", res)

	return res

}

func main() {
	// Set up a connection to the chosen broker chosen via round-robin
	address := "0.0.0.0:" + strconv.Itoa(decidePartition())
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()


	// msgProd := messagepb.NewMessageServiceClient(conn)
	metadataProd := metadatapb.NewMetadataServiceClient(conn)
	
	// send(prod)
	getMetadata(metadataProd, []string{"topic1","topic2"})
}
