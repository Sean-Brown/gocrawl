package tests

import (
	"fmt"
	"github.com/hydrogen18/stoppableListener"
	"log"
	"net"
	"net/http"
	"sync"
)

type HostPort struct {
	Host string
	Port int
}

func Serve(wait *sync.WaitGroup, quit chan int, ports chan HostPort) {
	// The number of servers
	const numServers = 4
	// A channel that will indicate to the servers to quit
	subQuit := make(chan int, 4)
	// Add this thread to the main thread's wait
	wait.Add(1)
	defer wait.Done()
	// Create a wait mechanism for the servers and begin serving
	subWait := sync.WaitGroup{}
	serveAndWait("hosta", subQuit, &subWait, ports)
	serveAndWait("hostb", subQuit, &subWait, ports)
	serveAndWait("hostc", subQuit, &subWait, ports)
	serveAndWait("hostd", subQuit, &subWait, ports)
	/* wait for a quit signal */
	log.Println("Waiting for quit signal")
	<-quit
	/* received the quit signal, signal the other server threads to quit */
	log.Println("Received the quit signal")
	for ix := 0; ix < numServers; ix++ {
		subQuit <- 1
	}
	/* wait for the threads to quit */
	log.Println("Waiting for server threads")
	subWait.Wait()
	/* done waiting, now exit */
	log.Println("Done waiting for server threads")
}

func serveAndWait(host string, quit chan int, wait *sync.WaitGroup, ports chan HostPort) {
	// Increase the wait delta
	wait.Add(1)
	// Serve on the host in a separate go routine
	go func() {
		defer wait.Done()
		serve(host, quit, ports)
	}()
}

func getListener(host string) (net.Listener, int) {
	port, err := NewPort()
	if err != nil {
		log.Fatal(err)
	}
	listener, err2 := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err2 != nil {
		log.Fatal(err2)
	}
	return listener, port
}

func getStoppableListener(host string) (*stoppableListener.StoppableListener, int) {
	listener, port := getListener(host)
	stoppable, err := stoppableListener.New(listener)
	if err != nil {
		log.Fatal(err)
	}
	return stoppable, port
}

func serve(host string, quit chan int, ports chan HostPort) {
	// Get a stoppable HTTP listener and the port the listener is listening on
	listener, port := getStoppableListener(host)
	var wait sync.WaitGroup
	wait.Add(1)
	// Serve the files for that host in a separate go routine
	go func() {
		defer wait.Done()
		ports <- HostPort{host, port}
		http.Serve(listener, http.FileServer(http.Dir(fmt.Sprintf("./%s", host))))
	}()
	log.Printf("Host %s listening on port %d\n", host, port)
	// Wait for the signal to quit
	select {
	case <-quit:
		log.Printf("Host %s received the quit signal\n", host)
	}
	// Received the quit signal, stop the HTTP listener
	log.Printf("Stopping the listener (%s)\n", host)
	listener.Stop()
	// Wait for the HTTP listener to stop
	log.Printf("Waiting on the server (%s)\n", host)
	wait.Wait()
	// Exit the routine

}
