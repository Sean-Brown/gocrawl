package crawler

import (
	"sync"
	//"log"
)

type DataConsumer struct {
	/* Compose with the Consumer struct */
	Consumer
	/* the channel of data that the consumer will parse */
	data chan DataCollection
}
/* Make a new Data consumer */
func NewDataConsumer(data chan DataCollection, quit chan int) *DataConsumer {
	c := &DataConsumer{
		Consumer: Consumer {
			quit: quit,
			waitGroup: &sync.WaitGroup{},
		},
		data: data,
	}
	c.waitGroup.Add(1)
	return c
}

/* rules for parsing data from the DOM */
type DOMParsingRules struct {
	tags []string
}

/* Consumption Loop */
//func (consumer *DataConsumer) Consume() {
//	defer consumer.waitGroup.Done()
//	for {
//		select {
//		case <-consumer.quit:
//			log.Println("data consumer received the quit signal")
//			break
//		case data := <-consumer.data:
//			log.Println("data consumer received data for ", data.url)
//			consumer.consume(data)
//		}
//	}
//}

/* Consume the data */
func (consumer *DataConsumer) consume(data DataCollection) {
	/* TODO consume the data */
}