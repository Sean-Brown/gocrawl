package crawler

import (
	"sync"
	"log"
)

type DataConsumer struct {
	/* the channel of data that the consumer will parse */
	data chan DataCollection
	/* channel to tell the consumer when to quit */
	quit chan int
	/* waitGroup to signal to the main service when the consumer thread has stopped */
	waitGroup *sync.WaitGroup
}

/* Make a new Data consumer */
func NewDataConsumer(data chan DataCollection, quit chan int) *DataConsumer {
	c := &DataConsumer{
		data: data,
		quit: quit,
		waitGroup: &sync.WaitGroup{},
	}
	c.waitGroup.Add(1)
	return c
}

/* Consumption Loop */
func (consumer *DataConsumer) Consume() {
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
func (consumer *DataConsumer) consume(data DataCollection) {
	/* TODO consume the data */
}