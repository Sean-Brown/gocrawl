package tests

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"testing"
	"github.com/Sean-Brown/gocrawl/gocrawl"
	"github.com/Sean-Brown/gocrawl/config"
	"github.com/stretchr/testify/assert"
	"strings"
)

type CountingDS struct {
	config.DataStorage
	// Map of url to # of times it was written to
	store map[string] int
}
func (ds *CountingDS) Store(url string, data string) {
	ds.store[url]++
}
func (ds *CountingDS) Get(url string) string {
	return fmt.Sprint(ds.store[url])
}
func (ds *CountingDS) NumItems() int {
	return len(ds.store)
}
func NewCountingDS() *CountingDS {
	return &CountingDS{store: make(map[string]int)}
}

/**
 Return a string formatted as http://<host>:2015/<page>.html
 */
func formatPage(host string, page string) string {
	return fmt.Sprintf("http://%s:%d/%s.html", host, 2015, page)
}

func getConfig(inMemory bool, host string, page string, sameDomain bool, depth int, dataSelectors []string) config.Config {
	var dataParsingRules []config.DataParsingRule
	for _, selector := range dataSelectors {
		dataParsingRules = append(dataParsingRules, config.DataParsingRule{
			UrlMatch: "host.*",
			DataSelector: selector,
		})
	}
	var ds config.DataStorage
	if inMemory {
		ds = gocrawl.CreateInMemoryDataStore()
	} else {
		ds = NewCountingDS()
	}
	return config.Config {
		StartUrl: formatPage(host, page),
		UrlParsingRules: config.CreateURLParsingRules(sameDomain, depth),
		DataParsingRules: dataParsingRules,
		DataStore: ds,
	}
}

func runTest(crawlConfig config.Config) (*gocrawl.GoCrawl, chan int, *sync.WaitGroup) {
	/* create a channel to receive OS interrupts on in case the user gets impatient */
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	/* start the caddy server */
	quit := make(chan int)
	wait := sync.WaitGroup{}
	ready := make(chan int, 1)
	go CaddyServe(&wait, quit, ready)
	/* wait for the server to start serving */
readyLoop:
	for {
		fmt.Println("in the readyLoop")
		select {
		case interrupt := <-sig:
			fmt.Println(interrupt)
			return nil, quit, &wait
		case <-ready:
			fmt.Println("caddyserver running, now crawl")
			break readyLoop
		}
	}
	fmt.Println("now crawling")
	/* create the gocrawler */
	gc := gocrawl.NewGoCrawl()
	/* start the crawler and wait for it to finish or for an OS interrupt */
	done := make(chan int, 1)
	go gc.Crawl(crawlConfig, quit, done)
loop:
	for {
		fmt.Println("in the loop")
		select {
		case interrupt := <-sig:
			fmt.Println(interrupt)
			break loop
		case <-done:
			fmt.Println("crawl completed")
			break loop
		}
	}
	return &gc, quit, &wait
}

func endTest(quit chan int, wait *sync.WaitGroup) {
	/* Tell the servers to quit */
	quit <- 0
	/* Wait for the servers to quit */
	fmt.Println("Wait for the caddy server threads")
	wait.Wait()
	fmt.Println("Done waiting, Goodbye!")
}

func dataAreEqual(actual string, expected string) bool {
	return strings.Compare(actual, expected) == 0
}

func expectData(t *testing.T, crawler *gocrawl.GoCrawl, url string, expected string) {
	data := (*crawler).GetDS().Get(url)
	assert.True(t, dataAreEqual(data, expected), fmt.Sprintf("%s != %s", data, expected))
}

func Test_HostA_Page1_SameDomain_Depth1(t *testing.T) {
	/* run the test */
	crawler, quit, wait := runTest(getConfig(
		true,
		"hosta",
		"page1",
		true,
		1,
		[]string{"p#data", "div#data", "div#ultra-cool p.data"},
	))

	if crawler != nil {
		/* assert that we got the data we expect to */
		expectData(t, crawler, formatPage("hosta", "page1"), "Page 1 Data")
		expectData(t, crawler, formatPage("hosta", "page2"), "Page 2 Data")
		expectData(t, crawler, formatPage("hosta", "page3"), "Page 3 Data")
	}

	/* end the test */
	endTest(quit, wait)
}

func Test_HostA_Page1_SameDomain_Depth2(t *testing.T) {
	/* run the test */
	crawler, quit, wait := runTest(getConfig(
		true,
		"hosta",
		"page1",
		true,
		2,
		[]string{"p#data", "div#data", "div#ultra-cool p.data"},
	))

	if crawler != nil {
		/* assert that we got the data we expect to */
		expectData(t, crawler, formatPage("hosta", "page1"), "Page 1 Data")
		expectData(t, crawler, formatPage("hosta", "page2"), "Page 2 Data")
		expectData(t, crawler, formatPage("hosta", "page3"), "Page 3 Data")
	}

	/* end the test */
	endTest(quit, wait)
}
func Test_DoesNotCrawlSamePageTwice(t *testing.T) {
	/* run the test */
	crawler, quit, wait := runTest(getConfig(
		false,
		"hosta",
		"page1",
		true,
		2,
		[]string{"p#data", "div#data", "div#ultra-cool p.data"},
	))

	if crawler != nil {
		/* assert that hosta was only crawled once */
		assert.Equal(t, crawler.GetDS().Get(formatPage("hosta", "page1")), "1")
	}

	/* end the test */
	endTest(quit, wait)
}

func Test_HostA_Page1_NotSameDomain_Depth2(t *testing.T) {
	/* run the test */
	crawler, quit, wait := runTest(getConfig(
		true,
		"hosta",
		"page1",
		false,
		2,
		[]string {
			"p#data", // hosta/page1
			"div#data", // hosta/page2
			"div#ultra-cool p.data", // hosta/page3
			"div#list ul li", // hostb/page1
			"h1#important", // hostb/page2
			"p span", // hostc/page1
			"div h1", // hostc/page2
		},
	))

	if crawler != nil {
		fmt.Println(crawler.GetDS())
		/* assert that we got the data we expect to */
		expectData(t, crawler, formatPage("hosta", "page1"), "Page 1 Data")
		expectData(t, crawler, formatPage("hosta", "page2"), "Page 2 Data")
		expectData(t, crawler, formatPage("hosta", "page3"), "Page 3 Data")
		expectData(t, crawler, formatPage("hostb", "page1"), "Hello World")
		expectData(t, crawler, formatPage("hostb", "page2"), "Page 2B Title")
		expectData(t, crawler, formatPage("hostc", "page1"), "Page 1 Data")
		expectData(t, crawler, formatPage("hostc", "page2"), "Page 2C Header")
	}

	/* end the test */
	endTest(quit, wait)
}

//func Test_HostA_Page1_NotSameDomain_Depth3(t *testing.T) {
//	/* run the test */
//	crawler, quit, wait := runTest(getConfig(
//		true,
//		"hosta",
//		"page1",
//		false,
//		3,
//		[]string {
//			"p#data", // hosta/page1
//			"div#data", // hosta/page2
//			"div#ultra-cool p.data", // hosta/page3
//			"div#list ul li", // hostb/page1
//			"h1#important", // hostb/page2
//			"p span", // hostc/page1
//			"div h1", // hostc/page2
//		},
//	))
//
//	if crawler != nil {
//		fmt.Println(crawler.GetDS())
//		/* assert that we got the data we expect to */
//		expectData(t, crawler, formatPage("hosta", "page1"), "Page 1 Data")
//		expectData(t, crawler, formatPage("hosta", "page2"), "Page 2 Data")
//		expectData(t, crawler, formatPage("hosta", "page3"), "Page 3 Data")
//		expectData(t, crawler, formatPage("hostb", "page1"), "Hello World")
//		expectData(t, crawler, formatPage("hostb", "page2"), "Page 2B Title")
//		expectData(t, crawler, formatPage("hostc", "page1"), "Page 1 Data")
//		expectData(t, crawler, formatPage("hostc", "page2"), "Page 2C Header")
//	}
//
//	/* end the test */
//	endTest(quit, wait)
//}
