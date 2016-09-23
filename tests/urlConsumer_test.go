package tests

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"testing"
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
		log.Fatal("Unable to create a new document for ", site)
	}
	resp, err := http.Get(site)
	if err != nil {
		log.Fatal("Failed to get ", site)
	}
	log.Println(resp.StatusCode)
	var body []byte
	resp.Body.Read(body)
	log.Println(string(body))
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
