package gocrawl

import (
	"fmt"
	"time"

	"github.com/Sean-Brown/gocrawl/config"
)

/*
GoCrawl - Main application that crawls the given domain, following a given set of rules
*/
type GoCrawl struct {
	/* The url consumer -- crawls and parses URLs */
	UrlConsumer *UrlConsumer
	/* The data consumer -- parses the downloaded DOM for data */
	dataConsumer *DataConsumer
	/* The data storage */
	dataStore config.DataStorage
}

/*
NewGoCrawl - Construct an empty GoCrawl instance
*/
func NewGoCrawl() GoCrawl {
	return GoCrawl{}
}

/*
GetDS - Get the GoCrawler's data store
*/
func (crawler *GoCrawl) GetDS() config.DataStorage {
	return crawler.dataStore
}

/*
Crawl - Begin crawling
*/
func (crawler *GoCrawl) Crawl(crawlConfig config.Config, quit chan int, done chan int) {
	/* create a channel to funnel URLs through and add the start url */
	urls := make(chan UrlData, 1)
	fmt.Println("Starting with ", crawlConfig.StartUrl)
	urls <- InitUrlData(crawlConfig.StartUrl, 0)
	/* channel to funnel URL DOM data through */
	data := make(chan DomQuery)

	/* Consume urls */
	quit2 := make(chan int)
	crawler.UrlConsumer = NewUrlConsumer(urls, data, quit2, crawlConfig.UrlParsingRules)
	go crawler.UrlConsumer.Consume()
	quit3 := make(chan int)
	crawler.dataStore = crawlConfig.DataStore
	crawler.dataConsumer = NewDataConsumer(data, quit3, crawlConfig.DataParsingRules, crawler.dataStore)
	go crawler.dataConsumer.Consume()

	/* Wait until signalled to quit or until there are no more URLs and data */
loop:
	for {
		select {
		case <-quit:
			fmt.Println("quit signal received")
			break loop
		default:
			// are there no more urls and data?
			noData := len(urls) == 0 && len(data) == 0
			noWorkers := crawler.dataConsumer.NumWorkers() == 0 && crawler.UrlConsumer.NumWorkers() == 0
			if noData && noWorkers {
				// no more urls or data, break out of the loop
				break loop
			}
			// sleep for one second before continuing
			time.Sleep(1 * time.Second)
		}
	}
	/* Tell the servers to quit */
	fmt.Println("Quit the threads")
	quit2 <- 1
	quit3 <- 1
	/* Wait for the consumer threads to quit */
	crawler.UrlConsumer.WaitGroup.Wait()
	crawler.dataConsumer.WaitGroup.Wait()
	fmt.Println("Done waiting, Goodbye!")
	/* indicate that the routine is done */
	done <- 0
}
