package tests

import (
	"fmt"
	"github.com/mholt/caddy"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func CaddyServe(wait *sync.WaitGroup, quit chan int) {
	wait.Add(1)
	defer wait.Done()
	defer removeCaddyfile()
	// Start the server
	instance := serve2()
	// Wait for the kill signal
	<-quit
	instance.Stop()
}

func removeCaddyfile() {
	// Remove the Caddyfile
	caddyfilePath := getCaddyfilePath()
	err := os.Remove(caddyfilePath)
	if err != nil {
		log.Println("Failed to remove the Caddyfile at ", caddyfilePath)
	}
}
func getCaddyfilePath() string {
	// Copy the Caddyfile to the Caddy installation directory
	caddyPath, err := exec.LookPath("caddy")
	if err != nil {
		log.Fatal("unable to find the Caddy installation directory")
	}
	caddyDir := filepath.Dir(caddyPath)
	if !os.IsPathSeparator(caddyDir[len(caddyDir)-1]) {
		caddyDir = caddyDir + string(os.PathSeparator)
	}
	return fmt.Sprintf("%sCaddyfile", caddyDir)
}

func loadCaddyfile() caddy.Input {
	// Copy the Caddyfile to the Caddy installation directory
	removeCaddyfile()
	caddyfilePath := getCaddyfilePath()
	err := os.Link("Caddyfile", caddyfilePath)
	if err != nil {
		log.Fatal("Failed to copy the Caddfyfile to ", caddyfilePath)
	}
	// Load the Caddyfile
	caddyfile, err := caddy.LoadCaddyfile("http")
	if err != nil {
		log.Fatal("Failed to load the Caddyfile")
	}
	return caddyfile
}

func serve2() *caddy.Instance {
	caddy.AppName = "gocrawl-tests"
	// Open the caddy file
	caddyfile := loadCaddyfile()
	// Start the server
	instance, err := caddy.Start(caddyfile)
	if err != nil {
		log.Fatal("Unable to start the caddy server: ", err)
	}
	// Wait for requests
	instance.Wait()
	return instance
}
