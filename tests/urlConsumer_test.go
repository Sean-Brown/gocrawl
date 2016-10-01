package tests

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"testing"
	"github.com/Sean-Brown/gocrawl/gocrawl"
	"github.com/Sean-Brown/gocrawl/config"
)

func getConfig(host string, sameDomain bool, depth int, dataSelector string) config.Config {
	return config.Config{
		StartUrl: fmt.Sprintf("http://%s:%d/page1.html", host, 2015),
		UrlParsingRules: config.CreateURLParsingRules(sameDomain, depth),
		DataParsingRules: []config.DataParsingRule{
			{
				UrlMatch:"host.*",
				DataSelector:dataSelector,
			},
		},
	}
}

func Test_HostA_SameDomain_Depth1(t *testing.T) {
	quit := make(chan int)
	wait := sync.WaitGroup{}
	go CaddyServe(&wait, quit)
	/* channel to receive OS interrupts on */
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	/* create the gocrawler */
	gc := gocrawl.NewGoCrawl()
	crawlConfig := getConfig("hosta", true, 1, "p#data, div#data, div#ultra-cool p.data")
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

	ds := gc.GetDS()
	if ds.Get("http://hosta:2015/page1.html") == "" {
		t.Fatal("Should have had hosta page1 data")
	}

	/* Tell the servers to quit */
	quit <- 0
	/* Wait for the servers to quit */
	fmt.Println("Wait for the caddy server threads")
	wait.Wait()
	fmt.Println("Done waiting, Goodbye!")
}
