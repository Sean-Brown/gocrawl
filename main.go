package main

import (
	"flag"
	"fmt"
	"github.com/Sean-Brown/gocrawl/config"
	"github.com/Sean-Brown/gocrawl/gocrawl"
	"os"
	"os/signal"
)

var configPath = flag.String("conf", "", "The path to the configuration file")

func main() {
	/* *** Setup *** */
	/* read the configuration file */
	conf := config.ReadConfig(*configPath)
	/* create the channel to funnel URLs through and add the start url */
	urls := make(chan gocrawl.UrlData, 1)
	urls <- gocrawl.InitUrlData(conf.StartUrl, 0)
	/* channel to funnel URL DOM data through */
	data := make(chan gocrawl.DomQuery)
	/* channel to signal workers to quit */
	quit := make(chan int, 1)
	/* channel to receive OS interrupts on */
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	/* *** Consume URLs and data *** */
	UrlConsumer := gocrawl.NewUrlConsumer(urls, data, quit, conf.UrlParsingRules)
	go UrlConsumer.Consume()
	dataConsumer := gocrawl.NewDataConsumer(data, quit, conf.DataParsingRules, &gocrawl.InMemoryDataStorage{})
	go dataConsumer.Consume()

	/* *** Wait until the program receives an interrupt *** */
	interrupt := <-sig
	fmt.Println(interrupt)
	/* signal the threads to quit */
	quit <- 1
	/* wait for the threads to exit */
	UrlConsumer.WaitGroup.Wait()
	dataConsumer.WaitGroup.Wait()
}
