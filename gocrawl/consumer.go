package gocrawl

import "sync"

type iConsumer interface {
	Consume()
}

type Consumer struct {
	/* this struct implement the "Consumer" interfaces */
	iConsumer
	/* channel to tell the consumer when to quit */
	Quit chan int
	/* waitGroup to signal to the main service when the consumer thread has stopped */
	WaitGroup *sync.WaitGroup
}
