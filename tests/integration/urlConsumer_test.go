package integration

import (
	"testing"
	"github.com/Sean-Brown/gocrawl/gocrawl"
	"github.com/Sean-Brown/gocrawl/gctesting"
	"os"
	"os/signal"
	"log"
	"sync"
)

func TestConsumesAllURLS(t *testing.T) {
	/* start the server */
	quit := make(chan int)
	wait := sync.WaitGroup{}
	ports := make(chan int, 4)
	go gctesting.Serve(&wait, quit, ports)
	/* create the channel to funnel URLs through and add the start url */
	urls := make(chan gocrawl.URLData, 1)
	urls <- gocrawl.InitURLData("http://hosta", 0)
	/* channel to funnel URL DOM data through */
	data := make(chan gocrawl.DataCollection)
	/* channel to signal workers to quit */
	quit := make(chan int, 1)
	/* channel to receive OS interrupts on */
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	/* *** Consume URLs and data *** */
	urlConsumer := gocrawl.NewURLConsumer(urls, data, quit, conf.UrlParsingRules)
	go urlConsumer.Consume()
	dataConsumer := gocrawl.NewDataConsumer(data, quit, conf.DataParsingRules, &gocrawl.InMemoryDataStorage{})
	go dataConsumer.Consume()

	/* *** Wait until the program receives an interrupt *** */
	interrupt := <-sig
	log.Println(interrupt)
	/* signal the threads to quit */
	quit <- 1
	/* wait for the threads to exit */
	urlConsumer.WaitGroup.Wait()
	dataConsumer.WaitGroup.Wait()

}