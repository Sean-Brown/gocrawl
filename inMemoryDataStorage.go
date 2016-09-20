package gocrawl

import "log"

type InMemoryDataStorage struct {
	/* Implement the DataStorage interface */
	DataStorage
}

func (storage *InMemoryDataStorage) Store(data string) {
	log.Println("Storing data: ", data)
}
