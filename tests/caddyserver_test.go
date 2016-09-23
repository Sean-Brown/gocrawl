package tests

import (
	"os"
	"os/signal"
	"sync"
	"testing"
)

func TestCaddyServe(t *testing.T) {
	// Start the server
	quit := make(chan int, 1)
	wait := sync.WaitGroup{}
	go CaddyServe(&wait, quit)
	// Wait for an OS interrupt
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	// Interrupt received, tell the servers to quit
	quit <- 1
	// Wait for the servers to quit
	wait.Wait()
}
