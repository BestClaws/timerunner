package jsondisk

import (
	"encoding/json"
	"os"
)

func Load[T any](path string) (*T, error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var obj T
	err = json.Unmarshal(fileData, &obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

func Save[T any](data T, path string) error {

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}
