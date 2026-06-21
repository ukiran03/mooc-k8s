package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var defaultPongFile = "pong.txt"

type Ping struct {
	file    string
	counter int
	mu      sync.Mutex
}

func (p *Ping) incrementAndSave() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.counter++
	data := []byte(strconv.Itoa(p.counter))
	err := os.WriteFile(p.file, data, 0o644)
	if err != nil {
		log.Printf("Error writing to file: %v", err)
	}
	return p.counter
}

// Implementing the http.Handler interface by naming this ServeHTTP
func (p *Ping) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Ensure we only respond to the exact path we want
	if r.URL.Path != "/pingpong" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "ping %d", p.incrementAndSave())
}

func main() {
	outputFile := os.Getenv("PONGFILE")
	if outputFile == "" {
		outputFile = defaultPongFile
	}
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("env PORT was unset\nUsing Port 3001")
		port = "3001"
	}
	addr := ":" + port

	srv := &Ping{
		file:    outputFile,
		counter: 0,
	}

	http.Handle("/pingpong", srv)

	log.Printf("Server starting on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
