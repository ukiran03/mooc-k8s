package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Error: env PORT was unset")
	}

	addr := ":" + port
	http.HandleFunc("/", helloHandler)

	log.Printf("Server started in port %s", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Todo App: Hello world!"))
}
