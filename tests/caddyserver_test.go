package tests

import (
	"os"
	"os/signal"
	"sync"
	"testing"
)

/**
 * This function was here for my early testing of caddy, it requires an interrupt to quit
 * and therefore doesn't play nice with other automated tests
 */
func DISABLED_TestCaddyServe(t *testing.T) {
	// Start the server
	quit := make(chan int, 1)
	wait := sync.WaitGroup{}
	ready := make(chan int, 1)
	go CaddyServe(&wait, quit, ready)
	// Wait for an OS interrupt
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	// Interrupt received, tell the servers to quit
	quit <- 1
	// Wait for the servers to quit
	wait.Wait()
}
