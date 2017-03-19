package gocrawl

import (
	"github.com/Sean-Brown/gocrawl/config"
)

/*
InMemoryDataStorage - Store data in cache as a key-value association
*/
type InMemoryDataStorage struct {
	/* Implement the DataStorage interfaces */
	config.DataStorage
	/* Store data in a map of <url, data> */
	ds map[string]string
}

/*
CreateInMemoryDataStore - Construct a new instance
*/
func CreateInMemoryDataStore() *InMemoryDataStorage {
	return &InMemoryDataStorage{ds: make(map[string]string)}
}

/*
Get the data for the given URL
*/
func (storage *InMemoryDataStorage) Get(url string) string {
	return storage.ds[url]
}

/*
NumItems - the number of items in the data store
*/
func (storage *InMemoryDataStorage) NumItems() int {
	return len(storage.ds)
}

/*
Store - store data for the given URL
*/
func (storage *InMemoryDataStorage) Store(url string, data string) {
	storage.ds[url] = data
}
