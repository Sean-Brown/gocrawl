package gocrawl

import (
	"flag"
	"os/signal"
	"os"
	"log"
)

var config = flag.String("config", "", "The path to the configuration file")

func main() {
	/* *** Setup *** */
	/* channel to funnel URLs through */
	urls := make(chan string, 1)
	urls <- *urlStart
	/* channel to funnel URL DOM data through */
	data := make(chan DataCollection)
	/* channel to signal workers to quit */
	quit := make(chan int, 1)
	/* channel to receive OS interrupts on */
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	/* *** Consume URLs and data *** */
	urlConsumer := NewURLConsumer(urls, data, quit, NewURLParsingRules())
	go urlConsumer.Consume()
	dataConsumer := NewDataConsumer(data, quit, nil, &InMemoryDataStorage{})
	go dataConsumer.Consume()

	/* *** Wait until the program receives an interrupt *** */
	interrupt := <-sig
	log.Println(interrupt)
	/* signal the threads to quit */
	quit <- 1
	/* wait for the threads to exit */
	urlConsumer.waitGroup.Wait()
	dataConsumer.waitGroup.Wait()
}
