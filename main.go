package gocrawl

import (
	"time"
	"os/signal"
	"os"
	"log"
	"github.com/PuerkitoBio/goquery"
)

type DataCollection struct {
	/* the url fetched */
	url string
	/* the time the url was fetched */
	fetched time.Time
	/* the parsed data */
	dom *goquery.Document
}

func Init(url string, dom *goquery.Document) DataCollection {
	return DataCollection{url:url, fetched:time.Now().UTC(), dom:dom}
}

func main() {
	/* *** Setup *** */
	/* channel to funnel URLs through */
	urls := make(chan string)
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
