package main

import (
	"flag"
	"github.com/Sean-Brown/gocrawl/config"
	"github.com/Sean-Brown/gocrawl/gocrawl"
	"log"
	"os"
	"os/signal"
)

var configPath = flag.String("conf", "", "The path to the configuration file")

func main() {
	/* *** Setup *** */
	/* read the configuration file */
	conf := config.ReadConfig(*configPath)
	/* create the channel to funnel URLs through and add the start url */
	urls := make(chan gocrawl.URLData, 1)
	urls <- gocrawl.InitURLData(conf.StartUrl, 0)
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
