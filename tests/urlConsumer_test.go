package tests

import (
	//"github.com/Sean-Brown/gocrawl/gocrawl"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"testing"
	//"github.com/Sean-Brown/gocrawl/config"
	"github.com/PuerkitoBio/goquery"
	"github.com/Sean-Brown/gocrawl/config"
	"github.com/Sean-Brown/gocrawl/gocrawl"
)

func TestConsumesAllURLS(t *testing.T) {
	quit := make(chan int)
	wait := sync.WaitGroup{}
	go CaddyServe(&wait, quit)
	/* channel to receive OS interrupts on */
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	/* create a dummy goquery document to get data from a page */
	site := fmt.Sprintf("http://%s:%d", "hosta", 2015)
	log.Println("Creating document for site ", site)
	doc, err := goquery.NewDocument(site)
	if err != nil {
		log.Fatalf("Unable to create a new document for %s\n", site)
	}
	log.Println(doc.Text())

	/* Wait until the program receives an interrupt */
	interrupt := <-sig
	log.Println(interrupt)
	/* Tell the servers to quit */
	log.Println("Quit the server threads")
	quit <- 1
	/* Wait for the servers to quit */
	log.Println("Wait for the server threads")
	wait.Wait()
	log.Println("Done waiting, Goodbye!")
}

func tTestConsumesAllURLS(t *testing.T) {
	quit := make(chan int)
	wait := sync.WaitGroup{}
	ports := make(chan HostPort, 4)
	go Serve(&wait, quit, ports)
	/* channel to receive OS interrupts on */
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Wait to get all the ports
	listenPorts := [4]HostPort{}
	for ix := 0; ix < 4; ix++ {
		listenPorts[ix] = <-ports
		log.Printf("*Host %s listening on port %d\n", listenPorts[ix].Host, listenPorts[ix].Port)
	}
	/* create a dummy goquery document to get data from a page */
	site := fmt.Sprintf("http://%s:%d", listenPorts[0].Host, listenPorts[0].Port)
	log.Println("Creating document for site ", site)
	doc, err := goquery.NewDocument(site)
	if err != nil {
		log.Fatalf("Unable to create a new document for %s\n", site)
	}
	log.Println(doc.Text())

	/* Wait until the program receives an interrupt */
	interrupt := <-sig
	log.Println(interrupt)
	/* Tell the servers to quit */
	log.Println("Quit the server threads")
	quit <- 1
	/* Wait for the servers to quit */
	log.Println("Wait for the server threads")
	wait.Wait()
	log.Println("Done waiting, Goodbye!")
}

func ttTestConsumesAllURLS(t *testing.T) {
	quitServ := make(chan int)
	wait := sync.WaitGroup{}
	ports := make(chan HostPort, 4)
	go Serve(&wait, quitServ, ports)

	// Wait to get all the ports
	listenPorts := [4]HostPort{}
	for ix := 0; ix < 4; ix++ {
		listenPorts[ix] = <-ports
		log.Printf("*Host %s listening on port %d\n", listenPorts[ix].Host, listenPorts[ix].Port)
	}
	// create the urls channel and add the start url
	urls := make(chan gocrawl.URLData, 1)
	urls <- gocrawl.InitURLData(
		//fmt.Sprintf("http://%s:%d", listenPorts[0].Host, listenPorts[0].Port),
		fmt.Sprintf("http://%s:%d", "127.0.0.1", listenPorts[0].Port),
		0,
	)
	/* channel to funnel URL DOM data through */
	data := make(chan gocrawl.DataCollection)
	/* channel to signal the consumer to quit */
	quitCons := make(chan int, 1)

	/* *** Consume URLs *** */
	urlConsumer := gocrawl.NewURLConsumer(urls, data, quitCons, config.InitURLParsingRules(false, 2))
	go urlConsumer.Consume()

	go func() {
	loop:
		for {
			select {
			case dat := <-data:
				log.Println("Received data: ")
				log.Println("-- url: ", dat.URL)
				log.Println("-- fetched: ", dat.Fetched)
				log.Println("-- data: ", dat.DOM.Text())
			case <-quitCons:
				break loop
			}
		}
	}()
	/* channel to receive OS interrupts on */
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	/* Wait until the program receives an interrupt */
	interrupt := <-sig
	log.Println(interrupt)
	/* signal the consumer to quit */
	quitCons <- 1
	/* wait for the consumer threads to exit */
	urlConsumer.WaitGroup.Wait()
	/* Tell the servers to quit */
	log.Println("Quit the server threads")
	quitServ <- 1
	/* Wait for the servers to quit */
	log.Println("Wait for the server threads")
	wait.Wait()
	log.Println("Done waiting, Goodbye!")
}
