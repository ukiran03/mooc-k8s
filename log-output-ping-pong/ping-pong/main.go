package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

type Ping struct {
	mu    sync.RWMutex
	count int
}

func (p *Ping) increment() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.count++
	return p.count
}

func (p *Ping) getCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.count
}

func main() {
	pingPort := os.Getenv("PING_PORT")
	if pingPort == "" {
		fmt.Println("env PING_PORT was unset\nUsing Port 3001 as pingPort")
		pingPort = "3001"
	}

	srv := &Ping{count: 0}
	mux := http.NewServeMux()

	mux.HandleFunc("/", srv.homeHandler)
	mux.HandleFunc("/pings", srv.getPings)
	mux.HandleFunc("/pingpong", srv.pingHandler)

	fmt.Printf("Starting Ping server v2 on port %s...\n", pingPort)
	if err := http.ListenAndServe(":"+pingPort, mux); err != nil {
		log.Fatalf("Server failed: %s\n", err)
	}
}

func (p *Ping) homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "visit /pingpong")
}

func (p *Ping) getPings(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ping %d", p.getCount())
}

func (p *Ping) pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ping %d", p.increment())
}
