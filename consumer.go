package gocrawl

import "sync"

type iConsumer interface {
	Consume()
}

type Consumer struct {
	/* implement the "Consumer" interface */
	iConsumer
	/* channel to tell the consumer when to quit */
	quit chan int
	/* waitGroup to signal to the main service when the consumer thread has stopped */
	waitGroup *sync.WaitGroup
}