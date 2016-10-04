package gocrawl

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/Sean-Brown/gocrawl/config"
	"github.com/bobesa/go-domain-util/domainutil"
	"strings"
	"sync"
	"net/url"
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
	// make the 'crawled' make thread-safe
	mux sync.RWMutex
}

func (consumer *URLConsumer) addCrawled(url string) {
	if !consumer.isCrawled(url) {
		// Set the site to 'crawled' -- obtain a write lock
		consumer.mux.Lock()
		consumer.crawled[url] = true
		consumer.mux.Unlock()
	}
}

func (consumer *URLConsumer) isCrawled(url string) bool {
	var crawled bool
	// Determine if this url is already set, only use a read lock
	consumer.mux.RLock()
	crawled = consumer.crawled[url]
	consumer.mux.RUnlock()
	return crawled
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
		crawled: make(map[string]bool),
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
			fmt.Println("url consumer received the quit signal")
			break loop
		case urlData := <-consumer.urls:
			// Count this worker as working
			consumer.IncWorkers()
			/* Download the DOM */
			doc, err := goquery.NewDocument(urlData.URL)
			if err != nil {
				fmt.Println(err)
			} else if urlData.Depth < consumer.rules.MaxDepth && !consumer.isCrawled(urlData.URL) {
				/* consume the document in a separate thread, increment for that thread */
				consumer.IncWorkers()
				go func() {
					// Defer decrementing the number of workers
					defer consumer.DecWorkers()
					fmt.Println("url consumer consuming:", urlData.URL)
					consumer.consume(doc, urlData.Depth+1)
				}()
				/* don't crawl this link again */
				consumer.addCrawled(urlData.URL)
			}
			// Uncount this worker
			consumer.DecWorkers()
		}
	}
}

/* Consume the url */
func (consumer *URLConsumer) consume(doc *goquery.Document, depth int) {
	// check that we won't exceed the max depth
	if (depth + 1) <= consumer.rules.MaxDepth {
		// the next depth is within the crawl limit, parse and enqueue the links on this page
		consumer.parseLinks(doc, depth)
	}

	// enqueue the data
	consumer.data <- InitDataCollection(doc.Url.String(), doc)
}

/* Parse and enqueue the links from the document */
func (consumer *URLConsumer) parseLinks(doc *goquery.Document, depth int) {
	// set the doamin equal to the host in the URL
	domain := doc.Url.Host
	if domainutil.HasSubdomain(domain) {
		// the domain contains a subdomain, parse out the top-level domain
		domain = domainutil.Domain(domain)
	}
	doc.Find(a).Each(func(_ int, sel *goquery.Selection) {
		href, exists := sel.Attr(href)
		shouldAdd, href := consumer.shouldAddLink(domain, href, depth)
		if exists && shouldAdd {
			fmt.Println("adding href", href, ", depth =", depth)
			consumer.urls <- InitURLData(href, depth)
		}
	})
}

/* Add the href if there is no domain restriction or if the href is in the domain, returns a possibly modified href */
func (consumer *URLConsumer) shouldAddLink(domain string, href string, currentDepth int) (bool, string) {
	shouldAdd := false
	/* if the parsed href has no domain, add the current domain */
	fmt.Println("URLConsumer found href", href)
	if strings.Index(href, "http://") == -1 && strings.Index(href, "www.") == -1 {
		// There's no "http" or "www." prefix, assume we're on the given domain,
		// at this point assume href is for a page on the same domain
		if strings.HasPrefix(href, "/") {
			// cut the leading '/', i.e. href="/page2"
			href = strings.TrimPrefix(href, "/")
		}
		href = fmt.Sprintf("http://%s/%s", domain, href)
		fmt.Println("Modified href to", href)
	}
	/* see if the href should be added to the urls channel */
	if consumer.rules.SameDomain {
		/* check that the domains are equal */
		// Check if the domains are equal
		_url, err := url.Parse(href)
		if err != nil {
			fmt.Println("error parsing <", href, "> into a url.URL struct")
		} else if _url.Host == "" && strings.Contains(_url.String(), domain) {
			// The _url host is empty, so just see if the full path contains the domain
			shouldAdd = true
		} else if strings.EqualFold(domain, _url.Host) {
			/* the domain and _url.Host are equal, enqueue the href */
			shouldAdd = true
		}
	} else {
		/* enqueue the href without checking the domain */
		shouldAdd = true
	}
	if shouldAdd {
		// check that the href hasn't been crawled yet
		shouldAdd = !consumer.isCrawled(href)
	}
	return shouldAdd, href
}
