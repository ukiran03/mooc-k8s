package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	pingPort string
	logPort  string
)

type LogEnt struct {
	timestamp string
	randStr   string
	pingCount int
}

func (l LogEnt) String() string {
	return fmt.Sprintf("%s %s\n Ping / Pongs: %d\n",
		l.timestamp, l.randStr, l.pingCount)
}

func randomString() string {
	b := make([]byte, 4) // 4 bytes = 8 hex characters
	if _, err := rand.Read(b); err != nil {
		return "00000000"
	}
	return hex.EncodeToString(b)
}

func init() {
	pingPort = os.Getenv("PING_PORT")
	if pingPort == "" {
		fmt.Println("env PING_PORT was unset\nUsing Port 3001 as pingPort")
		pingPort = "3001"
	}

	logPort = os.Getenv("PORT")
	if logPort == "" {
		fmt.Println("env PORT was unset\nUsing Port 3000 as logPort")
		logPort = "3000"
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	hostAddr := "http://localhost:" + pingPort + "/pings"

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	resp, err := client.Get(hostAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("bad status: %s", resp.Status), resp.StatusCode)
		return
	}

	reader := io.LimitReader(resp.Body, 1024*1024)
	data, err := io.ReadAll(reader)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	var pongCount int
	_, err = fmt.Sscanf(string(data), "ping %d", &pongCount)
	if err != nil {
		pongCount = 0
	}

	entry := LogEnt{
		timestamp: time.Now().Format("2006-01-02 15:04:05"),
		randStr:   randomString(),
		pingCount: pongCount,
	}
	fmt.Fprint(w, entry)
}

func main() {
	http.HandleFunc("/", homeHandler)
	fmt.Printf("Starting Log server v2 on port %s...\n", logPort)
	if err := http.ListenAndServe(":"+logPort, nil); err != nil {
		log.Fatalf("Server failed: %s\n", err)
	}
}
