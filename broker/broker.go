package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	messagepb "AKFAK/proto/messagepb"
	metadatapb "AKFAK/proto/metadatapb"
	"google.golang.org/grpc"
)

type messageServer struct{}
type metadataServer struct{}

func (*messageServer) MessageBatch(stream messagepb.MessageService_MessageBatchServer) error {
	fmt.Println("MessageBatch function was invoked with a streaming request")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		records := req.GetRecords()
		// TODO: get fields from records and do stuff with it
		fmt.Println("Records: %v", records)

		result := "Received"
		sendErr := stream.Send(&messagepb.MessageBatchResponse{
			Result: result,
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", sendErr)
			return sendErr
		}
	}
}

func (*metadataServer) GetMetadata(ctx context.Context, req *metadatapb.MetadataRequest) (*metadatapb.MetadataResponse, error) {
	fmt.Printf("GetMetadata function was invoked with %v\n", req)

	//topicNames := req.GetTopicNames()

	// TODO: get brokers
	broker := &metadatapb.Broker{
		NodeID: 1,
		Host:   "127.0.0.1",
		Port:   1,
	}

	// TODO: get topics
	topic := &metadatapb.Topic{
		ErrorCode: 0,
		Name:      "name",
		Partitions: []*metadatapb.Partition{
			{
				ErrorCode:      0,
				PartitionIndex: 1,
				LeaderID:       1,
			},
			{
				ErrorCode:      0,
				PartitionIndex: 2,
				LeaderID:       1,
			},
		},
	}

	res := &metadatapb.MetadataResponse{
		Brokers: []*metadatapb.Broker{broker},
		Topics:  []*metadatapb.Topic{topic},
	}
	return res, nil
}

func main() {

	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		return
	}

	// message batch
	// gRPCmsgServer := grpc.NewServer()
	// messagepb.RegisterMessageServiceServer(gRPCmsgServer, &messageServer{})

	// if err := gRPCmsgServer.Serve(listen); err != nil {
	// 	log.Fatalf("Failed to serve: %v", err)
	// }

	// metadata
	gRPCmetadataServer := grpc.NewServer()
	metadatapb.RegisterMetadataServiceServer(gRPCmetadataServer, &metadataServer{})

	if err := gRPCmetadataServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
