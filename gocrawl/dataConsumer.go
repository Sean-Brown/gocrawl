package gocrawl

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/Sean-Brown/gocrawl/config"
	"regexp"
	"sync"
	"fmt"
	"strings"
)

type DataConsumer struct {
	/* Compose with the Consumer struct */
	Consumer
	/* the channel of data that the consumer will parse */
	data chan DataCollection
	/* the rules for parsing the DOM */
	rules []config.DataParsingRule
	/* a dependency-injected data storage object to persist the data */
	storage config.DataStorage
}

/* Make a new Data consumer */
func NewDataConsumer(data chan DataCollection, quit chan int, rules []config.DataParsingRule, storage config.DataStorage) *DataConsumer {
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
loop:
	for {
		select {
		case <-consumer.Quit:
			fmt.Println("data consumer received the quit signal")
			break loop
		case data := <-consumer.data:
			fmt.Println("data consumer received data for ", data.URL)
			// increment the number of worker threads for this consumer
			consumer.IncWorkers()
			go func() {
				// Defer decrementing the number of workers
				defer consumer.DecWorkers()
				consumer.consume(data)
			}()
		}
	}
}

/* Consume the data */
func (consumer *DataConsumer) consume(data DataCollection) {
	// the text extracted from the DOM
	text := ""
	/* iterate the DOM-parsing rules */
	for _, rule := range consumer.rules {
		/* check if this rule applies to this url */
		matched, err := regexp.MatchString(rule.UrlMatch, data.URL)
		if err != nil {
			fmt.Println("Error matching url regex <", rule.UrlMatch, "> with ", data.URL)
		} else if matched {
			/* the rule does apply to this url, apply the rule */
			data.DOM.Find(rule.DataSelector).Each(func(_ int, sel *goquery.Selection) {
				/* store the data */
				tmp := strings.TrimSpace(sel.Text())
				if len(text) > 0 {
					if len(tmp) > 0 {
						// append the data with a space
						text = strings.Join([]string{text, tmp}, " ")
					}
				} else {
					text = tmp
				}
			})
		}
	}
	fmt.Println("Storing data", text, " for url", data.URL)
	consumer.storage.Store(data.URL, text)
}
