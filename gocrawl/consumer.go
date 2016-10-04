package gocrawl

import "sync"

type iConsumer interface {
	Consume()
	IncWorkers()
	DecWorkers()
	NumWorkers() int
}

type Consumer struct {
	/* this struct implement the "Consumer" interfaces */
	iConsumer
	/* channel to tell the consumer when to quit */
	Quit chan int
	/* waitGroup to signal to the main service when the consumer thread has stopped */
	WaitGroup *sync.WaitGroup
	/* the number of threads working */
	numWorkers int
	/* mutex for interacting with numWorkers */
	workerMux sync.RWMutex
}

func (consumer *Consumer) IncWorkers() {
	consumer.workerMux.Lock()
	consumer.numWorkers++
	consumer.workerMux.Unlock()
}
func (consumer *Consumer) DecWorkers() {
	consumer.workerMux.Lock()
	consumer.numWorkers--
	consumer.workerMux.Unlock()
}
func (consumer *Consumer) NumWorkers() int {
	var numWorkers int
	consumer.workerMux.RLock()
	numWorkers = consumer.numWorkers
	consumer.workerMux.RUnlock()
	return numWorkers
}