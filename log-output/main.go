package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
)

type Pair struct {
	key   string
	value string
}

func secureRandomString() string {
	b := make([]byte, 4) // 4 bytes = 8 hex characters
	rand.Read(b)
	return hex.EncodeToString(b)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("env PORT was unset\nUsing Port 3000")
		port = "3000"
	}
	addr := ":" + port
	keyStr := secureRandomString()
	srv := &Pair{key: keyStr}

	http.ListenAndServe(addr, srv)
}

func (p *Pair) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	val := secureRandomString()
	fmt.Fprintf(w, "%s: %s", p.key, val)
}
