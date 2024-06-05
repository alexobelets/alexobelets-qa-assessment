package utils

import (
	"encoding/json"
	"log"
	"os"
)

func JsonWriter(contractInformation interface{}, filePath string) {
	jsonData, err := json.MarshalIndent(contractInformation, "", "  ")
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		log.Fatal("Error writing JSON file:", err)
		return
	}
	log.Println("Contract data is written to JSON file successfully!")
}
