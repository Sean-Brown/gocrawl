package testing

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"testing"
)

func TestServe(t *testing.T) {
	quit := make(chan int)
	wait := sync.WaitGroup{}
	go Serve(&wait, quit)
	/* channel to receive OS interrupts on */
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

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
