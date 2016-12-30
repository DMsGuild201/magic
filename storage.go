package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Storage interface {
	Save(interface{}) error
}

type jsonStorage struct {
	*os.File
}

func JsonStorage(name string) (*jsonStorage, error) {
	f, err := os.OpenFile(fmt.Sprintf("_output/%s.json", name), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &jsonStorage{
		File: f,
	}, nil
}

func (j jsonStorage) Save(data interface{}) error {
	defer j.Close()
	return json.NewEncoder(j).Encode(&data)
}
