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

type HostPageSelData struct {
	host string
	page string
	sel string
	data string
}
const (
	// Hosts
	_hosta = "hosta"
	_hostb = "hostb"
	_hostc = "hostc"
	_hostd = "hostd"
	// Pages
	_page1 = "page1"
	_page2 = "page2"
	_page3 = "page3"
	_page4 = "page4"
	_page5 = "page5"
)
const (
	// Hosts
	hosta = 1 << 0
	hostb = hosta << 1
	hostc = hosta << 2
	hostd = hosta << 3
	// Pages
	page1 = hosta << 4
	page2 = hosta << 5
	page3 = hosta << 6
	page4 = hosta << 7
	page5 = hosta << 8
	// Combined
	hap1 = hosta | page1
	hap2 = hosta | page2
	hap3 = hosta | page3
	hap4 = hosta | page4
	hap5 = hosta | page5
	hbp1 = hostb | page1
	hbp2 = hostb | page2
	hcp1 = hostc | page1
	hcp2 = hostc | page2
	hcp3 = hostc | page3
	hdp1 = hostd | page1
)
var selectors map[int]string = map[int]string {
	hap1: "p#data",
	hap2: "div#data",
	hap3: "div#ultra-cool p.data",
	hap4: "select[title='cookies'] option",
	hap5: "div#help p#important",

	hbp1: "div#list ul li",
	hbp2: "h1#important",

	hcp1: "p span",
	hcp2: "div h1",
	hcp3: "td",

	hdp1: "h1",
}
var data map[int]string = map[int]string {
	hap1: "Page 1 Data",
	hap2: "Page 2 Data",
	hap3: "Page 3 Data",
	hap4: "sugar oatmeal raisin oreo",
	hap5: "world",

	hbp1: "Hello World",
	hbp2: "Page 2B Title",

	hcp1: "Page 1 Data",
	hcp2: "Page 2C Header",
	hcp3: "c1 data c2 data c3 data",

	hdp1: "Page D1 Title",
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
		_hosta,
		_page1,
		true,
		1,
		[]string{
			selectors[hap1],
			selectors[hap2],
			selectors[hap3],
		},
	))

	if crawler != nil {
		/* assert that we got the data we expect to */
		expectData(t, crawler, formatPage(_hosta, _page1), data[hap1])
		expectData(t, crawler, formatPage(_hosta, _page2), data[hap2])
		expectData(t, crawler, formatPage(_hosta, _page3), data[hap3])
	}

	/* end the test */
	endTest(quit, wait)
}

func Test_HostA_Page1_SameDomain_Depth2(t *testing.T) {
	/* run the test */
	crawler, quit, wait := runTest(getConfig(
		true,
		_hosta,
		_page1,
		true,
		2,
		[]string{
			selectors[hap1],
			selectors[hap2],
			selectors[hap3],
		},
	))

	if crawler != nil {
		/* assert that we got the data we expect to */
		expectData(t, crawler, formatPage(_hosta, _page1), data[hap1])
		expectData(t, crawler, formatPage(_hosta, _page2), data[hap2])
		expectData(t, crawler, formatPage(_hosta, _page3), data[hap3])
	}

	/* end the test */
	endTest(quit, wait)
}
func Test_DoesNotCrawlSamePageTwice(t *testing.T) {
	/* run the test */
	crawler, quit, wait := runTest(getConfig(
		false,
		_hosta,
		_page1,
		true,
		2,
		[]string{
			selectors[hap1],
			selectors[hap2],
			selectors[hap3],
		},
	))

	if crawler != nil {
		/* assert that hosta was only crawled once */
		assert.Equal(t, crawler.GetDS().Get(formatPage(_hosta, _page1)), "1")
	}

	/* end the test */
	endTest(quit, wait)
}

func Test_HostA_Page1_NotSameDomain_Depth2(t *testing.T) {
	/* run the test */
	crawler, quit, wait := runTest(getConfig(
		true,
		_hosta,
		_page1,
		false,
		2,
		[]string {
			selectors[hap1],
			selectors[hap2],
			selectors[hap3],

			selectors[hbp1],
			selectors[hbp2],

			selectors[hcp1],
			selectors[hcp2],
		},
	))

	if crawler != nil {
		fmt.Println(crawler.GetDS())
		/* assert that we got the data we expect to */
		expectData(t, crawler, formatPage(_hosta, _page1), data[hap1])
		expectData(t, crawler, formatPage(_hosta, _page2), data[hap2])
		expectData(t, crawler, formatPage(_hosta, _page3), data[hap3])

		expectData(t, crawler, formatPage(_hostb, _page1), data[hbp1])
		expectData(t, crawler, formatPage(_hostb, _page2), data[hbp2])

		expectData(t, crawler, formatPage(_hostc, _page1), data[hcp1])
		expectData(t, crawler, formatPage(_hostc, _page2), data[hcp2])
	}

	/* end the test */
	endTest(quit, wait)
}

func Test_HostA_Page1_NotSameDomain_Depth3(t *testing.T) {
	/* run the test */
	crawler, quit, wait := runTest(getConfig(
		true,
		_hosta,
		_page1,
		false,
		3,
		[]string {
			selectors[hap1],
			selectors[hap2],
			selectors[hap3],

			selectors[hbp1],
			selectors[hbp2],

			selectors[hcp1],
			selectors[hcp2],
		},
	))

	if crawler != nil {
		fmt.Println(crawler.GetDS())
		/* assert that we got the data we expect to */
		expectData(t, crawler, formatPage(_hosta, _page1), data[hap1])
		expectData(t, crawler, formatPage(_hosta, _page2), data[hap2])
		expectData(t, crawler, formatPage(_hosta, _page3), data[hap3])

		expectData(t, crawler, formatPage(_hostb, _page1), data[hbp1])
		expectData(t, crawler, formatPage(_hostb, _page2), data[hbp2])

		expectData(t, crawler, formatPage(_hostc, _page1), data[hcp1])
		expectData(t, crawler, formatPage(_hostc, _page2), data[hcp2])
	}

	/* end the test */
	endTest(quit, wait)
}

func Test_HostA_Page4_NotSameDomain_Depth2(t *testing.T) {
	/* run the test */
	crawler, quit, wait := runTest(getConfig(
		true,
		_hosta,
		_page4,
		false,
		3,
		[]string {
			selectors[hap4],
			selectors[hap5],

			selectors[hcp3],

			selectors[hdp1],
		},
	))

	if crawler != nil {
		fmt.Println(crawler.GetDS())
		/* assert that we got the data we expect to */
		expectData(t, crawler, formatPage(_hosta, _page4), data[hap4])
		expectData(t, crawler, formatPage(_hosta, _page5), data[hap5])

		expectData(t, crawler, formatPage(_hostc, _page3), data[hcp3])

		//expectData(t, crawler, formatPage(_hostd, _page1), data[hdp1])
	}

	/* end the test */
	endTest(quit, wait)
}
