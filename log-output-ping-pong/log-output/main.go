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

var ( // read fron ENV
	pingPort    string
	logPort     string
	message     string
	fileContent string
)

type logEntry struct {
	timestamp   string
	randStr     string
	pingCount   int
	message     string
	fileContent string
}

func (le logEntry) String() string {
	return fmt.Sprintf("%s\n%s\n%s %s\n Ping / Pongs: %d\n",
		le.message, le.fileContent, le.timestamp, le.randStr, le.pingCount)
}

func genRandomString() string {
	b := make([]byte, 4) // 4 bytes = 8 hex characters
	if _, err := rand.Read(b); err != nil {
		return "00000000"
	}
	return hex.EncodeToString(b)
}

func main() {
	pingPort = os.Getenv("PING_PORT")
	if pingPort == "" {
		fmt.Println("env PING_PORT was unset\nUsing Port 3001 as pingPort")
		pingPort = "3001"
	}

	logPort = os.Getenv("LOG_PORT")
	if logPort == "" {
		fmt.Println("env LOG_PORT was unset\nUsing Port 3000 as logPort")
		logPort = "3000"
	}

	message = os.Getenv("MESSAGE")
	if message == "" {
		log.Print("env MESSAGE was unset")
	}

	f := os.Getenv("INFO_FILE")
	if f == "" {
		log.Print("env INFO_FILE was unset")
	} else {
		data, err := os.ReadFile(f)
		if err != nil {
			log.Printf("Read INFO_FILE Error: %v", err)
		}
		fileContent = string(data)
	}

	http.HandleFunc("/", homeHandler)
	log.Printf("Starting log server on port %s...", logPort)
	if err := http.ListenAndServe(":"+logPort, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	pingCount, err := fetchPingCount(pingPort)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		log.Println(err.Error())
		return
	}

	entry := logEntry{
		timestamp:   time.Now().Format("2006-01-02 15:04:05"),
		randStr:     genRandomString(),
		pingCount:   pingCount,
		message:     message,
		fileContent: fileContent,
	}
	fmt.Fprint(w, entry)
}

func fetchPingCount(pingPort string) (int, error) {
	hostAddr := "http://localhost:" + pingPort + "/pings"
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	resp, err := client.Get(hostAddr)
	if err != nil {
		return 0, fmt.Errorf("failed to reach ping service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bad status from ping service: %s", resp.Status)
	}

	reader := io.LimitReader(resp.Body, 1024*1024)
	data, err := io.ReadAll(reader)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	var pingCount int
	_, err = fmt.Sscanf(string(data), "ping %d", &pingCount)
	if err != nil {
		return 0, nil
	}
	return pingCount, nil
}
