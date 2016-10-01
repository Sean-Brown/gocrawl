package gocrawl

import (
	"testing"
)

func TestInMemoryDataStorage_Store(t *testing.T) {
	store := CreateInMemoryDataStore()
	store.Store("test.com", "hello world")
	if store.ds["test.com"] != "hello world" {
		t.Fatal("Failed to store data")
	}
}
