package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// LoadJSONData load the data into the given buffer
func LoadJSONData(path string, buffer interface{}) error {
	// load the JSON byte data
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Unable to read the JSON file: %v\n", err)
		return err
	}

	// parse the JSON byte into structs
	if err := json.Unmarshal([]byte(data), buffer); err != nil {
		log.Printf("Failed to parse the JSON file: %v\n", err)
		return err
	}
	return nil
}

//FlushJSONData write the JSON data into file
func FlushJSONData(path string, data interface{}) error {
	// parse the data into JSON byte
	bytes, marshallErr := json.MarshalIndent(data, "", " ")
	if marshallErr != nil {
		log.Println("Unable to convert data into JSON:", marshallErr)
		return marshallErr
	}

	// flush the JSON byte intp file
	writeErr := ioutil.WriteFile(path, bytes, 0644)
	if writeErr != nil {
		log.Println("Unable to save JSON to file:", writeErr)
		return writeErr
	}

	return nil
}
