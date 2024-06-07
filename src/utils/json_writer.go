package utils

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// Creates directory if it does not exist
func createDirectoryIfNotExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func JsonWriter(contractInformation interface{}, filePath string) {
	err := createDirectoryIfNotExists(filepath.Dir(filePath))
	if err != nil {
		log.Fatal("Error creating directory:", err)
		return
	}

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
