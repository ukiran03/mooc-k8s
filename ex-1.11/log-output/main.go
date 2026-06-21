package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var defaultPongFile = "pong.txt"

type LogEnt struct {
	timestamp string
	randStr   string
}

func (l LogEnt) String() string {
	return fmt.Sprintf("%s %s\n", l.timestamp, l.randStr)
}

func randomString() string {
	b := make([]byte, 4) // 4 bytes = 8 hex characters
	if _, err := rand.Read(b); err != nil {
		return "00000000"
	}
	return hex.EncodeToString(b)
}

type Data struct {
	file string
}

func (d *Data) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.URL.Path != "/log" {
		http.NotFound(w, r)
		return
	}

	fileBytes, err := os.ReadFile(d.file)

	var count string
	if err != nil {
		if os.IsNotExist(err) {
			count = "0"
		} else {
			log.Printf("Error reading file %s: %v", d.file, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		count = string(fileBytes)
	}

	entry := LogEnt{
		timestamp: time.Now().Format("2006-01-02 15:04:05"),
		randStr:   randomString(),
	}

	fmt.Fprintf(w, "%sPing / Pongs: %s", entry, count)
}

func main() {
	inputFile := os.Getenv("PONGFILE")
	if inputFile == "" {
		inputFile = defaultPongFile
	}
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("env PORT was unset\nUsing Port 3000")
		port = "3000"
	}
	addr := ":" + port

	srv := &Data{file: inputFile}

	http.Handle("/log", srv)

	log.Printf("Server starting on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
