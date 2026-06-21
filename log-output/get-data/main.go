// Reads that file and provides the content in the HTTP GET endpoint
// for the user to see

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var defaultOutFile = "output.txt"

type Data struct {
	file    string
	content []byte
}

func main() {
	outFile := os.Getenv("OUTFILE")
	if outFile == "" {
		outFile = defaultOutFile
	}

	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("env PORT was unset\nUsing Port 3000")
		port = "3000"
	}
	addr := ":" + port

	logdata := &Data{
		file:    outFile,
		content: []byte{},
	}

	http.HandleFunc("/", logdata.homeHandler)
	log.Printf("Server starting on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (d *Data) readData() error {
	data, err := os.ReadFile(d.file)
	if err != nil {
		return err
	}
	d.content = data
	return nil
}

func (d *Data) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := d.readData(); err != nil {
		log.Printf("Error reading file: %v", err)
		http.Error(w, "Could not read file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(d.content)
}
