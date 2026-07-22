package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Ping struct {
	counter int
}

func (p *Ping) incrementCounter() int {
	p.counter++
	return p.counter
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("env PORT was unset\nUsing Port 3001")
		port = "3001"
	}
	addr := ":" + port
	srv := &Ping{counter: 0}

	http.HandleFunc("/pingpong", srv.pingHandler)

	log.Printf("Server starting on %s\n", addr)
	log.Fatal((http.ListenAndServe(addr, nil)))
}

func (p *Ping) pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	fmt.Fprintf(w, "ping %d", p.incrementCounter())
}
