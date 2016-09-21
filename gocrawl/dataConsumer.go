package gocrawl

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"regexp"
	"sync"
	"github.com/Sean-Brown/gocrawl/config"
)

type DataConsumer struct {
	/* Compose with the Consumer struct */
	Consumer
	/* the channel of data that the consumer will parse */
	data chan DataCollection
	/* the rules for parsing the DOM */
	rules []config.DataParsingRule
	/* a dependency-injected data storage object to persist the data */
	storage DataStorage
}

/* Make a new Data consumer */
func NewDataConsumer(data chan DataCollection, quit chan int, rules []config.DataParsingRule, storage DataStorage) *DataConsumer {
	if rules == nil {
		rules = []config.DataParsingRule{}
	}
	c := &DataConsumer{
		Consumer: Consumer{
			Quit:      quit,
			WaitGroup: &sync.WaitGroup{},
		},
		data:    data,
		rules:   rules,
		storage: storage,
	}
	c.WaitGroup.Add(1)
	return c
}

/* Consumption Loop */
func (consumer *DataConsumer) Consume() {
	defer consumer.WaitGroup.Done()
	for {
		select {
		case <-consumer.Quit:
			log.Println("data consumer received the quit signal")
			break
		case data := <-consumer.data:
			log.Println("data consumer received data for ", data.url)
			go consumer.consume(data)
		}
	}
}

/* Consume the data */
func (consumer *DataConsumer) consume(data DataCollection) {
	/* iterate the DOM-parsing rules */
	for _, rule := range consumer.rules {
		/* check if this rule applies to this url */
		matched, err := regexp.MatchString(rule.UrlMatch, data.url)
		if err != nil {
			log.Println("Error matching url regex <", rule.UrlMatch, "> with ", data.url)
		} else if matched {
			/* the rule does apply to this url, apply the rule */
			log.Println("Matched <", rule.UrlMatch, "> to ", data.url)
			data.dom.Find(rule.DataSelector).Each(func(_ int, sel *goquery.Selection) {
				/* store the data */
				consumer.storage.Store(data.url, sel.Text())
			})
		}
	}
}
