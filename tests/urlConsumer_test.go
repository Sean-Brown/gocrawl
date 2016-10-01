package tests

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"testing"
	"github.com/Sean-Brown/gocrawl/gocrawl"
	"github.com/Sean-Brown/gocrawl/config"
	"strings"
)

func getConfig(host string, page string, sameDomain bool, depth int, dataSelector string) config.Config {
	return config.Config{
		StartUrl: fmt.Sprintf("http://%s:%d/%s.html", host, 2015, page),
		UrlParsingRules: config.CreateURLParsingRules(sameDomain, depth),
		DataParsingRules: []config.DataParsingRule{
			{
				UrlMatch:"host.*",
				DataSelector:dataSelector,
			},
		},
	}
}

func runTest(crawlConfig config.Config) (gocrawl.GoCrawl, chan int, *sync.WaitGroup) {
	quit := make(chan int)
	wait := sync.WaitGroup{}
	/* start the caddy server */
	go CaddyServe(&wait, quit)

	/* create a channel to receive OS interrupts on in case the user gets impatient */
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	/* create the gocrawler */
	gc := gocrawl.NewGoCrawl()
	/* start the crawler and wait for it to finish or for an OS interrupt */
	done := make(chan int, 1)
	go gc.Crawl(crawlConfig, quit, done)
loop:
	for {
		select {
		case interrupt := <-sig:
			fmt.Println(interrupt)
			break loop
		case <-done:
			fmt.Println("crawl completed")
			break loop
		}
	}
	return gc, quit, &wait
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

func Test_HostA_Page1_SameDomain_Depth1(t *testing.T) {
	/* run the test */
	crawler, quit, wait := runTest(getConfig(
		"hosta",
		"page1",
		true,
		1,
		"p#data",
	))

	/* assert that we got the data we expect to */
	ds := crawler.GetDS()
	data := ds.Get("http://hosta:2015/page1.html")
	expected := "Page 1 Data"
	if !dataAreEqual(data, expected) {
		t.Fatal(data, " != ", expected)
	}

	/* end the test */
	endTest(quit, wait)
}
