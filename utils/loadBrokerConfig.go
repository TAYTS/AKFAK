package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Broker struct {
	Id int32 `json:"id"`
	Host string `json:"host"`
	Port int32 `json:"port"`
	IsCoordinator bool `json:"isCoordinator"`
}

func GetBrokers(path string) []Broker {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Successfully opened broker_config.json")
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var brokers []Broker
	if err := json.Unmarshal([]byte(byteValue), &brokers); err!= nil {
		panic(err)
	}
	return brokers
}
