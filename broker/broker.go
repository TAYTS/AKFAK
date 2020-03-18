package main

import (
	"net"
	"log"
	"fmt"
	"io"

	"google.golang.org/grpc"
	messagepb "AKFAK/proto/messagepb"
)

type server struct {}

func (*server) MessageBatch(stream messagepb.MessageService_MessageBatchServer) error {
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
			Result:	result,
		})
		if sendErr != nil{
			log.Fatalf("Error while sending data to client: %v", sendErr)
			return sendErr
		}		
	}
}

func main(){

	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil{
		log.Fatalf("Failed to listen: %v", err)
		return
	}

	s := grpc.NewServer()
	messagepb.RegisterMessageServiceServer(s, &server{})

	// reflection.Register(s)

	if err := s.Serve(listen); err != nil{
		log.Fatalf("Failed to serve: %v",err)
	}
}
