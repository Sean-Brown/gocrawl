package tests

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"testing"
)

func TestCaddyServe(t *testing.T) {
	// Start the server
	quit := make(chan int)
	wait := sync.WaitGroup{}
	go CaddyServe(&wait, quit)
	// Wait for an OS interrupt
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	interrupt := <-sig
	// Interrupt received, tell the servers to quit
	log.Println(interrupt)
	log.Println("Quit the server threads")
	quit <- 1
	// Wait for the servers to quit
	log.Println("Wait for the server threads")
	wait.Wait()
	log.Println("Done waiting, Goodbye!")
}
