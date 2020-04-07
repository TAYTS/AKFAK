package utils

import (
	"AKFAK/proto/zkpb"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func GetBrokers(path string) []*zkpb.Broker {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Successfully opened broker_config.json")
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var brokers []*zkpb.Broker
	if err := json.Unmarshal([]byte(byteValue), &brokers); err!= nil {
		panic(err)
	}
	return brokers
}
