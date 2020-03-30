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

/*
********************************* NOTE ***************************************
Code below is only for local testing of send/receive functions for producer.go
******************************************************************************
*/

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
		fmt.Println("Records:", records)

		result := "Broker Received"
		sendErr := stream.Send(&messagepb.MessageBatchResponse{
			Result: result,
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", sendErr)
			return sendErr
		}
	}
}

func (*metadataServer) WaitOnMetadata(ctx context.Context, req *metadatapb.MetadataRequest) (*metadatapb.MetadataResponse, error) {
	fmt.Printf("WaitOnMetadata function was invoked with %v\n", req)

	broker1 := &metadatapb.Broker{
		NodeID: 1,
		Host:   "0.0.0.0",
		Port:   50052,
	}

	broker2 := &metadatapb.Broker{
		NodeID: 1,
		Host:   "0.0.0.0",
		Port:   50053,
	}

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
		Brokers: []*metadatapb.Broker{broker1, broker2},
		Topics:  topic,
	}
	return res, nil
}

func main() {
	// set up servers listening for messageBatch request
	go func() {
		listen1, err1 := net.Listen("tcp", "0.0.0.0:50052")
		if err1 != nil {
			log.Fatalf("Failed to listen: %v", err1)
			return
		}
		gRPCmsgServer1 := grpc.NewServer()
		messagepb.RegisterMessageServiceServer(gRPCmsgServer1, &messageServer{})

		if err := gRPCmsgServer1.Serve(listen1); err1 != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	go func() {
		listen2, err2 := net.Listen("tcp", "0.0.0.0:50053")
		if err2 != nil {
			log.Fatalf("Failed to listen: %v", err2)
			return
		}
		gRPCmsgServer2 := grpc.NewServer()

		messagepb.RegisterMessageServiceServer(gRPCmsgServer2, &messageServer{})
		if err := gRPCmsgServer2.Serve(listen2); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		return
	}

	// server for metadata
	gRPCmetadataServer := grpc.NewServer()
	metadatapb.RegisterMetadataServiceServer(gRPCmetadataServer, &metadataServer{})

	if err := gRPCmetadataServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
