package llpm

import (
	"encoding/json"
	"os"
)

func LoadJson(filePath string, result interface{}) error {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(fileBytes, result)
}

func SaveJson(filePath string, data interface{}) error {
	destination, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer destination.Close()

	fileBytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	_, err = destination.Write(fileBytes)
	return err
}
