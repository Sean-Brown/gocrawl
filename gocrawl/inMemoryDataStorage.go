package gocrawl

import (
	"github.com/Sean-Brown/gocrawl/config"
)

type InMemoryDataStorage struct {
	/* Implement the DataStorage interfaces */
	config.DataStorage
	/* Store data in a map of <url, data> */
	ds map[string]string
}

func CreateInMemoryDataStore() *InMemoryDataStorage {
	return &InMemoryDataStorage{ds: make(map[string]string)}
}

func (storage *InMemoryDataStorage) Get(url string) string {
	return storage.ds[url]
}

func (storage *InMemoryDataStorage) NumItems() int {
	return len(storage.ds)
}

func (storage *InMemoryDataStorage) Store(url string, data string) {
	storage.ds[url] = data
}
