package gocrawl

import (
	"flag"
	"log"
	"os"
	"os/signal"
)

var configPath = flag.String("config", "", "The path to the configuration file")

func main() {
	/* *** Setup *** */
	/* read the configuration file */
	config := ReadConfig(*configPath)
	/* create the channel to funnel URLs through and add the start url */
	urls := make(chan URLData, 1)
	urls <- InitURLData(config.StartUrl, 0)
	/* channel to funnel URL DOM data through */
	data := make(chan DataCollection)
	/* channel to signal workers to quit */
	quit := make(chan int, 1)
	/* channel to receive OS interrupts on */
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	/* *** Consume URLs and data *** */
	urlConsumer := NewURLConsumer(urls, data, quit, config.UrlParsingRules)
	go urlConsumer.Consume()
	dataConsumer := NewDataConsumer(data, quit, config.DataParsingRules, &InMemoryDataStorage{})
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
