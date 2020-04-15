package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// LoadJSONData load the data into the given buffer
func LoadJSONData(path string, buffer interface{}) {
	// load the JSON byte data
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}

	// parse the JSON byte into structs
	if err := json.Unmarshal([]byte(data), buffer); err != nil {
		panic(err)
	}
}
