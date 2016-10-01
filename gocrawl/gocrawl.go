package gocrawl

import (
	"github.com/Sean-Brown/gocrawl/config"
	"time"
	"fmt"
)

type GoCrawl struct {
	/* The url consumer -- crawls and parses URLs */
	urlConsumer *URLConsumer
	/* The data consumer -- parses the downloaded DOM for data */
	dataConsumer *DataConsumer
	/* The data storage */
	dataStore DataStorage
	/* Timeout -- the time (in seconds) to wait with empty channels, if this time elapses,
	then the crawler will assume there are no more urls or data left and will thus exit*/
	timeout int
}

func NewGoCrawl() GoCrawl {
	return GoCrawl{timeout: 3}
}

func (crawler *GoCrawl) GetDS() DataStorage {
	return crawler.dataStore
}

func (crawler *GoCrawl) Crawl(crawlConfig config.Config, quit chan int, done chan int) {
	/* create a channel to funnel URLs through and add the start url */
	urls := make(chan URLData, 1)
	fmt.Println("Starting with ", crawlConfig.StartUrl)
	urls <- InitURLData(crawlConfig.StartUrl, 0)
	/* channel to funnel URL DOM data through */
	data := make(chan DataCollection)

	/* Consume urls */
	quit2 := make(chan int)
	crawler.urlConsumer = NewURLConsumer(urls, data, quit2, crawlConfig.UrlParsingRules);
	go crawler.urlConsumer.Consume()
	quit3 := make(chan int)
	crawler.dataStore = CreateInMemoryDataStore()
	crawler.dataConsumer = NewDataConsumer(data, quit3, crawlConfig.DataParsingRules, crawler.dataStore)
	go crawler.dataConsumer.Consume()

	/* Wait until signalled to quit or until there are no more URLs and data */
	elapsed := time.Now()
loop:
	for {
		select {
		case <-quit:
			fmt.Println("quit signal received")
			break loop
		default:
			// are there no more urls and data?
			noData := len(urls) == 0 && len(data) == 0
			timeElapsed := (int)(time.Now().Sub(elapsed).Seconds())
			timedOut := (timeElapsed > crawler.timeout)
			if noData && timedOut {
				// no more urls or data, break out of the loop
				break loop
			} else if (!noData) {
				// there is still data in the channels, reset the elapsed timer
				elapsed = time.Now()
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
	crawler.urlConsumer.WaitGroup.Wait()
	crawler.dataConsumer.WaitGroup.Wait()
	fmt.Println("Done waiting, Goodbye!")
	/* indicate that the routine is done */
	done <- 0
}
