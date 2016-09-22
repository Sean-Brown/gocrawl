package gocrawl

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/Sean-Brown/gocrawl/config"
	"github.com/bobesa/go-domain-util/domainutil"
	"log"
	"strings"
	"sync"
)

/* The URL Consumer */
type URLConsumer struct {
	/* Compose with the Consumer struct */
	Consumer
	/* channel of urls (and their corresponding depth) that the url consumer consumes */
	urls chan URLData
	/* the channel of data that will be parsed by a data consumer */
	data chan DataCollection
	/* the parsing rules */
	rules config.URLParsingRules
	/* a store of urls already crawled */
	crawled map[string]bool
}

/* Make a new URL consumer */
func NewURLConsumer(urls chan URLData, data chan DataCollection, quit chan int, rules config.URLParsingRules) *URLConsumer {
	c := &URLConsumer{
		Consumer: Consumer{
			Quit:      quit,
			WaitGroup: &sync.WaitGroup{},
		},
		urls:  urls,
		data:  data,
		rules: rules,
	}
	c.WaitGroup.Add(1)
	return c
}

/* Consumption Loop */
func (consumer *URLConsumer) Consume() {
	defer consumer.WaitGroup.Done()
loop:
	for {
		select {
		case <-consumer.Quit:
			log.Println("url onsumer received the quit signal")
			break loop
		case urlData := <-consumer.urls:
			log.Println("url consumer consuming: ", urlData.URL)
			/* Download the DOM */
			doc, err := goquery.NewDocument(urlData.URL)
			if err != nil {
				log.Println(err)
			} else if urlData.Depth < consumer.rules.MaxDepth && !consumer.crawled[urlData.URL] {
				/* consume the document in a separate thread */
				go consumer.consume(doc, urlData.Depth+1)
			}
		}
	}
}

/* Consume the url */
func (consumer *URLConsumer) consume(doc *goquery.Document, depth int) {
	/* Parse and enqueue the links */
	consumer.parseLinks(doc, depth)

	/* enqueue the data */
	consumer.data <- InitDataCollection(doc.Url.String(), doc)
}

/* Parse and enqueue the links from the document */
func (consumer *URLConsumer) parseLinks(doc *goquery.Document, depth int) {
	domain := domainutil.Domain(doc.Url.String())
	doc.Find(a).Each(func(_ int, sel *goquery.Selection) {
		href, exists := sel.Attr(href)
		if exists {
			/* there is an href attribute, try adding it to the urls channel */
			if consumer.rules.SameDomain {
				/* check that the domains are equal */
				if strings.EqualFold(domain, domainutil.Domain(href)) {
					/* the domains are equal, enqueue the href */
					consumer.urls <- InitURLData(href, depth)
				}
			} else {
				/* enqueue the href without checking the domain */
				consumer.urls <- InitURLData(href, depth)
			}
		}
	})
}
