package main

import (
	"AKFAK/proto/clientpb"
	"AKFAK/proto/consumepb"
	"context"
	"fmt"
	"log"
	"io"
	"google.golang.org/grpc"
)

func main() {
	// test getassignment
	id := int32(1)
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("0.0.0.0:5000", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := clientpb.NewClientServiceClient(cc)
	
	/*
	req := &consumepb.GetAssignmentRequest{
		GroupID:	id,
		TopicName:	"topic1",
	}

	res, err := c.GetAssignment(context.Background(), req)
	if err != nil {
		log.Fatalf("Error whil calling controller RPC: %v\n", err)
	}

	fmt.Println("Result:", res)*/

	// test consume request
	
	// ====== //
	// Add this in front of the n.ReadRecordBatchFromLocal() part if testing without producer
	// fileHandlerMapping := make(map[int]*recordpb.FileRecord)
	// newRcd := recordpb.InitialiseRecordWithMsg("some random msg")
	// reqqq := producepb.InitProduceRequest(req.GetTopicName(), int(req.GetPartition()), newRcd)
	// recordBatchTest := reqqq.GetTopicData()[0].GetRecordSet()
	// // byteee, _ := proto.Marshal(recordBatchTest)
	// // log.Println("This should be the data that will be written:", byteee)
	// // test := &recordpb.RecordBatch{}
	// // proto.Unmarshal(byteee, test)
	// // log.Println(test)
	// n.writeRecordBatchToLocal(req.GetTopicName(), int(req.GetPartition()), fileHandlerMapping, recordBatchTest)
	//=====//


	// setup stream
	stream, err := c.Consume(context.Background())
	if err != nil {
		log.Fatalf("could not consume: %v", err)
	}
	stream.Send(&consumepb.ConsumeRequest{
		GroupID:	id,
		ConsumerID:	1,
		Partition:	1,
		TopicName:	"topic1",
	})
	res2, err2 := stream.Recv()
	
	if err2 != nil{
		if err2 == io.EOF {
			fmt.Println("Nothing to consume")
		} else {
			log.Fatalf("Error: %v\n", err2)
		}
	}
	fmt.Println("Response:", res2)
	
}