package main

import (
    "fmt"
    "log"
    "os"
    "time"
    "consumer/common"
	"consumer/consumergroup"
	"consumer/config"
)

// const (
//     zookeeperConn = "10.4.1.29:2181"
//     cgroup = "zgroup"
//     topic = "senz"
// )

func main() {
    // setup sarama log to stdout
    sarama.Logger = log.New(os.Stdout, "", log.Ltime)

    // init consumer
    cg, err := initConsumer()
    if err != nil {
        fmt.Println("Error consumer goup: ", err.Error())
        os.Exit(1)
    }
    defer cg.Close()

    // run consumer
    consume(cg)
}

func initConsumer()(*consumergroup.ConsumerGroup, error) {
    // consumer config
    config := consumergroup.NewConfig()
    config.Offsets.Initial = sarama.OffsetOldest
    config.Offsets.ProcessingTimeout = 10 * time.Second

    // join to consumer group
    cg, err := consumergroup.JoinConsumerGroup(cgroup, []string{topic}, []string{zookeeperConn}, config)
    if err != nil {
        return nil, err
    }

    return cg, err
}

func consume(cg *consumergroup.ConsumerGroup) {
    for {
        select {
        case msg := <-cg.Messages():
            // messages coming through chanel
            // only take messages from subscribed topic
	    if msg.Topic != topic {
                continue
            }

            fmt.Println("Topic: ", msg.Topic)
            fmt.Println("Value: ", string(msg.Value))

            // commit to zookeeper that message is read
            // this prevent read message multiple times after restart
            err := cg.CommitUpto(msg)
            if err != nil {
                fmt.Println("Error commit zookeeper: ", err.Error())
            }
        }
    }
}