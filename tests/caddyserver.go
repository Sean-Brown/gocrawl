package tests

import (
	"github.com/mholt/caddy"
	_ "github.com/mholt/caddy/caddyhttp/browse"
	_ "github.com/mholt/caddy/caddyhttp/httpserver"
	_ "github.com/mholt/caddy/caddyhttp/root"
	"log"
	"os"
	"sync"
)

func CaddyServe(wait *sync.WaitGroup, quit chan int) {
	wait.Add(1)
	defer wait.Done()
	// Create the server instance
	instance := serve2()
	// Start the server in a separate go-routine
	go func() {
		// Wait for requests
		instance.Wait()
	}()
	<-quit
	// Got a quit signal, stop waiting for requests
	instance.Stop()
}

func loadCaddyfile() caddy.Input {
	// Get a handle to the Caddyfile
	file, err := os.Open("Caddyfile")
	if err != nil {
		log.Fatal(err)
	}
	// Pipe the Caddyfile into Caddy
	caddyfile, err := caddy.CaddyfileFromPipe(file, "http")
	if err != nil {
		log.Fatal(err)
	}
	return caddyfile
}

func serve2() *caddy.Instance {
	caddy.AppName = "gocrawl-tests"
	// Load the caddy file
	caddyfile := loadCaddyfile()
	// Start the server
	instance, err := caddy.Start(caddyfile)
	if err != nil {
		log.Fatal("Unable to start the caddy server: ", err)
	}
	return instance
}
