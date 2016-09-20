package crawler

import (
	"sync"
	"log"
	"github.com/PuerkitoBio/goquery"
	"github.com/bobesa/go-domain-util/domainutil"
	"strings"
)

/* The URL Consumer */
type URLConsumer struct {
	/* Inheritance from the consumer struct */
	Consumer
	/* channel of strings that the consumer consumes */
	urls chan string
	/* the channel of data that will be parsed by a parser */
	data chan DataCollection
	/* the parsing rules */
	rules URLParsingRules
}

/* Make a new URL consumer */
func NewURLConsumer(urls chan string, data chan DataCollection, quit chan int, rules URLParsingRules) *URLConsumer {
	c := &URLConsumer{
		Consumer: Consumer{
			quit: quit,
			waitGroup: &sync.WaitGroup{},
		},
		urls: urls,
		data: data,
		rules: rules,
	}
	c.waitGroup.Add(1)
	return c
}

/* Rules for parsing */
type URLParsingRules struct {
	sameDomain bool
	dataTags []string
}

/* Default URL Parsing rules */
func NewURLParsingRules() URLParsingRules {
	return URLParsingRules{sameDomain: true, dataTags: []string{body}}
}
/* Initialize URL Parsing rules defined by the user */
func InitURLParsingRules(sameDomain bool, dataTags []string) URLParsingRules {
	if dataTags == nil || len(dataTags) == 0 {
		dataTags = []string{body}
	}
	return URLParsingRules{sameDomain: sameDomain, dataTags: dataTags}
}

/* Consumption Loop */
func (consumer *URLConsumer) Consume() {
	defer consumer.waitGroup.Done()
	for {
		select {
		case <-consumer.quit:
			log.Println("url onsumer received the quit signal")
			break
		case url := <-consumer.urls:
			log.Println("url consumer consuming: ", url)
			/* Download the DOM */
			doc, err := goquery.NewDocument(url)
			if err != nil {
				log.Println(err)
			} else {
				/* consume the document in a separate thread */
				go consumer.consume(doc)
			}
		}
	}
}

/* Consume the url */
func (consumer *URLConsumer) consume(doc *goquery.Document) {
	/* Parse and enqueue the links */
	consumer.parseLinks(doc)

	/* enqueue the data */
	consumer.data <- Init(doc.Url.String(), doc)
}

/* Parse and enqueue the links from the document */
func (consumer *URLConsumer) parseLinks(doc *goquery.Document) {
	domain := domainutil.Domain(doc.Url.String())
	doc.Find(a).Each(func (_ int, sel *goquery.Selection) {
		href, exists := sel.Attr(href)
		if exists {
			/* there is an href attribute, try adding it to the urls channel */
			if consumer.rules.sameDomain {
				/* check that the domains are equal */
				if strings.EqualFold(domain, domainutil.Domain(href)) {
					/* the domains are equal, enqueue the href */
					consumer.urls <- href
				}
			} else {
				/* enqueue the href without checking the domain */
				consumer.urls <- href
			}
		}
	});
}
