package main

import (
	"AKFAK/proto/clientpb"
	"AKFAK/proto/messagepb"
	"AKFAK/proto/recordpb"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:5001", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := clientpb.NewClientServiceClient(cc)

	stream, err := c.MessageBatch(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v\n", err)
	}

	reqs := []*messagepb.MessageBatchRequest{
		&messagepb.MessageBatchRequest{
			Records: &recordpb.RecordBatch{
				Records: []*recordpb.Record{
					&recordpb.Record{
						Length:   10,
						ValueLen: 10,
						Value:    []byte("testing1-1"),
					},
					&recordpb.Record{
						Length:   10,
						ValueLen: 10,
						Value:    []byte("testing1-2"),
					},
					&recordpb.Record{
						Length:   10,
						ValueLen: 10,
						Value:    []byte("testing1-3"),
					},
					&recordpb.Record{
						Length:   10,
						ValueLen: 10,
						Value:    []byte("testing1-4"),
					},
				},
			},
		},
		&messagepb.MessageBatchRequest{
			Records: &recordpb.RecordBatch{
				Records: []*recordpb.Record{
					&recordpb.Record{
						Length:   10,
						ValueLen: 10,
						Value:    []byte("testing2-1"),
					},
					&recordpb.Record{
						Length:   10,
						ValueLen: 10,
						Value:    []byte("testing2-2"),
					},
					&recordpb.Record{
						Length:   10,
						ValueLen: 10,
						Value:    []byte("testing2-3"),
					},
					&recordpb.Record{
						Length:   10,
						ValueLen: 10,
						Value:    []byte("testing2-4"),
					},
				},
			},
		},
		&messagepb.MessageBatchRequest{
			Records: &recordpb.RecordBatch{
				Records: []*recordpb.Record{
					&recordpb.Record{
						Length:   10,
						ValueLen: 10,
						Value:    []byte("testing3-1"),
					},
					&recordpb.Record{
						Length:   10,
						ValueLen: 10,
						Value:    []byte("testing3-2"),
					},
					&recordpb.Record{
						Length:   10,
						ValueLen: 10,
						Value:    []byte("testing3-3"),
					},
					&recordpb.Record{
						Length:   10,
						ValueLen: 10,
						Value:    []byte("testing3-4"),
					},
				},
			},
		},
	}

	waitc := make(chan struct{})
	go func() {
		for _, req := range reqs {
			stream.Send(req)
			time.Sleep(time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v\n", err)
			}
			fmt.Printf("Received: %v\n", res.GetResponse().GetStatus())
		}
		close(waitc)
	}()

	<-waitc
}
