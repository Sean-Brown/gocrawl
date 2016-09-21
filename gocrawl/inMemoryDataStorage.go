package gocrawl

import "log"

type InMemoryDataStorage struct {
	/* Implement the DataStorage interfaces */
	DataStorage
	/* Store data in a map of <url, data> */
	ds map[string]string
}

func (storage *InMemoryDataStorage) Store(url string, data string) {
	log.Println("Received data for url: ", url, ", data: ", data)
	storage.ds[url] = data
}
