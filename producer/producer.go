package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"io"

	messagepb "AKFAK/proto/messagepb"
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
	return 1
}

func send(producer messagepb.MessageServiceClient) {
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
	go func(){
		for _, msg := range msgs{
			fmt.Printf("Sending message: %v\n", msg)
			stream.Send(msg)
			time.Sleep(1000* time.Millisecond)
		}
		stream.CloseSend()
	}()

	// receive messages from broker
	go func () {
		for {
			res, err := stream.Recv()
			if err == io.EOF{
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


func main() {
	// Set up a connection to the chosen broker chosen via round-robin
	address := "0.0.0.0:50051"
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	prod := messagepb.NewMessageServiceClient(conn)

	send(prod)
}
