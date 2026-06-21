package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var OUTFILE string = "output.txt"

func randomString() string {
	b := make([]byte, 4) // 4 bytes = 8 hex characters
	if _, err := rand.Read(b); err != nil {
		return "00000000"
	}
	return hex.EncodeToString(b)
}

func main() {
	outFile := os.Getenv("OUTFILE")
	if outFile == "" {
		outFile = OUTFILE
	}

	f, err := os.OpenFile(outFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(5 * time.Second)

	defer func() {
		ticker.Stop()
		f.Sync()
		f.Close()
		fmt.Println("\nFile closed, Exiting gracefully.")
	}()

	str := randomString()
	fmt.Printf("Logger started with ID: %s. Press Ctrl+C to stop.\n", str)

	for {
		select {
		case <-ticker.C:
			data := fmt.Sprintf(
				"%s %s\n", (time.Now().Format("2006-01-02 15:04:05")), str,
			)
			if _, err := f.WriteString(data); err != nil {
				log.Print(err)
				return
			}
		case <-sigChan:
			return
		}
	}
}
