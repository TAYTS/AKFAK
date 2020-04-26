package consumer

import (
	"AKFAK/proto/metadatapb"
	"AKFAK/proto/recordpb"
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"AKFAK/proto/clientpb"
	"AKFAK/proto/consumepb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const METADATA_TIMEOUT = 500 * time.Millisecond

// Consumer is a Kafka consumer
type Consumer struct {
	ID int
	// assignments are given an idx (key)
	//brokerAddr    string
	topic         string
	brokerAddrMap map[int32]string
	metadata      *metadatapb.MetadataResponse
	brokerCon     clientpb.ClientService_ConsumeClient
	grpcConn      *grpc.ClientConn
	metadataMux   sync.RWMutex
	partitionIdx  int32 // partition index of the topic to consume
	readOffset    int64 // used to store the consumer current read offset
}

///////////////////////////////////
//         Public Methods        //
///////////////////////////////////

// InitConsumer will fetch the metadata for the topic requested and return the consumer instance and all the available partitions index
func InitConsumer(topic string, brokerAddr string) (*Consumer, []int32) {
	// Dial to broker to get metadata
	c := &Consumer{
		topic:      topic,
		readOffset: 0,
	}
	// get metadata and wait for 500ms
	err := c.waitOnMetadata(brokerAddr, METADATA_TIMEOUT)
	if err != nil {
		log.Fatalf("Unable to get Topic Metadata: %v\n", err)
	}

	// create and store the mapping of brokerID to broker connection
	// address for easy access later
	brkAddrMapping := make(map[int32]string)
	for _, brk := range c.metadata.GetBrokers() {
		brkAddrMapping[brk.GetNodeID()] = fmt.Sprintf("%v:%v", brk.GetHost(), brk.GetPort())
	}
	c.brokerAddrMap = brkAddrMapping

	// get all the available partition index
	partitions := c.metadata.GetTopic().GetPartitions()
	partitionIdx := make([]int32, 0, len(partitions))
	for _, part := range partitions {
		partitionIdx = append(partitionIdx, part.GetPartitionIndex())
	}

	return c, partitionIdx
}

// SetPartitionIdx is used set the partition index that the consumer need to pull from
func (c *Consumer) SetPartitionIdx(partIdx int) {
	partitionExist := false
	partIdxInt32 := int32(partIdx)

	for _, part := range c.metadata.GetTopic().GetPartitions() {
		if part.GetPartitionIndex() == partIdxInt32 {
			c.partitionIdx = partIdxInt32
			partitionExist = true
			break
		}
	}

	// fail the consumer if invalid partition index is given
	if !partitionExist {
		log.Fatalln("Invalid partition index")
	}
}

// Consume is used to start the messages pulling
func (c *Consumer) Consume() {
	// setup the stream connections to all the required brokers
	c.setupStreamToConsumeMsg()

	for {
		// setup the request
		req := &consumepb.ConsumeRequest{
			Partition: c.partitionIdx,
			TopicName: c.topic,
			Offset:    c.readOffset,
		}

		// send the consume request
		err := c.brokerCon.Send(req)
		if err != nil {
			// log.Println("err:", err)
			// reset the connection and try again
			c.resetBrokerConnection()
			continue
		}

		// get the consume response
		resp, err := c.brokerCon.Recv()
		if err == io.EOF {
			// reset the connection and try again
			c.resetBrokerConnection()
			continue
		} else {
			statusErr, ok := status.FromError(err)
			if ok {
				if statusErr.Message() == recordpb.ErrNoRecord.Error() {
					// if no record available sleep for 500 milliseconds
					time.Sleep(500 * time.Millisecond)
					// move to the next iteration
					continue
				} else {
					c.resetBrokerConnection()
				}
			} else {
				c.resetBrokerConnection()
			}
		}

		// print out all the record
		printRecordBatchMsg(resp.GetRecordSet())

		// update the consumer read offset
		c.updateReadOffset(resp.GetOffset())
	}
}

// CleanupResources used to cleanup the Consumer resources
func (c *Consumer) CleanupResources() {
	c.brokerCon.CloseSend()
	c.grpcConn.Close()
}

///////////////////////////////////
//         Private Methods       //
///////////////////////////////////

// waitOnMetadata used to fetch the metadata for the requested topic
func (c *Consumer) waitOnMetadata(brokerAddr string, maxWaitMs time.Duration) error {
	conn, err := grpc.Dial(brokerAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	// create gRPC client
	conClient := clientpb.NewClientServiceClient(conn)

	// send get metadata request
	req := &metadatapb.MetadataRequest{
		TopicName: c.topic,
	}

	// create context with timeout
	protoCtx, protoCancel := context.WithTimeout(context.Background(), maxWaitMs)
	defer protoCancel()

	// send request to get metadata
	res, err := conClient.WaitOnMetadata(protoCtx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Printf("Metadata not received for topic %v after %d ms\n", c.topic, maxWaitMs)
			} else {
				log.Printf("Unexpected error: %v\n", statusErr.Message())
			}
		} else {
			log.Printf("could not get metadata. error: %v\n", err)
		}
		return err
	}
	// save the metadata response to the consumer
	c.metadataMux.Lock()
	c.metadata = res
	c.metadataMux.Unlock()

	return nil
}

// setupStreamToConsumeMsg setup the gRPC stream for pulling message
// it will stop the consumer program if there is no broker available for the partition
func (c *Consumer) setupStreamToConsumeMsg() {

	// setup the consume stream connection
	for {
		leaderID := int32(-1)

		// get the leader ID for the partition
		for _, part := range c.metadata.GetTopic().GetPartitions() {
			if part.GetPartitionIndex() == c.partitionIdx {
				leaderID = part.GetLeaderID()
			}
		}

		// if no broker available for the current partition TERMINATE the consumer
		if leaderID == -1 {
			log.Fatalf("No brokers available for the current partition(%v)\n", c.partitionIdx)
		}

		// setup gRPC connection
		conn, err := grpc.Dial(c.brokerAddrMap[leaderID], grpc.WithInsecure())
		if err != nil {
			c.refreshMetadata()
			continue
		}
		// setup gRPC service
		conClient := clientpb.NewClientServiceClient(conn)

		// setup stream
		stream, err := conClient.Consume(context.Background())
		if err != nil {
			conn.Close()
			c.refreshMetadata()
			continue
		}

		// store the grpc and stream connection
		c.grpcConn = conn
		c.brokerCon = stream

		break
	}
}

// refreshMetadata fetch new metadata and replace the current metadata cache
// it will fail if there is no broker available to get the topic metadata
func (c *Consumer) refreshMetadata() {
	for _, brkAddr := range c.brokerAddrMap {
		// refresh the metadata
		err := c.waitOnMetadata(brkAddr, METADATA_TIMEOUT)
		if err != nil {
			continue
		}
		return
	}
	log.Fatalln("No brokers available")
}

// resetBrokerConnection is used to refresh the metadata, setup stream connection and check the partition available
func (c *Consumer) resetBrokerConnection() {
	// close the current grpc connection
	c.grpcConn.Close()

	c.refreshMetadata()

	c.setupStreamToConsumeMsg()
}

// printRecordBatchMsg used to print all the Record messages in the RecordBatch
func printRecordBatchMsg(rcdBatch *recordpb.RecordBatch) {
	for _, rcd := range rcdBatch.GetRecords() {
		fmt.Println(rcd.GetValue())
	}
}

// updateReadOffset used to update the consumer RecordBatch read offset
func (c *Consumer) updateReadOffset(offset int64) {
	c.readOffset = offset
}
