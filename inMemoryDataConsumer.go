package crawler

import (
	"sync"
	"log"
)

type InMemoryDataConsumer struct {
	/* Inherit from Consumer */
	Consumer
	/* the channel of data that the consumer will parse */
	data chan DataCollection
}

/* Make a new Data consumer */
func NewDataConsumer(data chan DataCollection, quit chan int) *InMemoryDataConsumer {
	c := &InMemoryDataConsumer{
		Consumer: Consumer {
			quit: quit,
			waitGroup: &sync.WaitGroup{},
		},
		data: data,
	}
	c.waitGroup.Add(1)
	return c
}

/* Consumption Loop */
func (consumer *InMemoryDataConsumer) Consume() {
	defer consumer.waitGroup.Done()
	for {
		select {
		case <-consumer.quit:
			log.Println("data consumer received the quit signal")
			break
		case data := <-consumer.data:
			log.Println("data consumer received data for ", data.url)

		}
	}
}

/* Consume the data */
func (consumer *InMemoryDataConsumer) consume(data DataCollection) {
	/* TODO consume the data */
}