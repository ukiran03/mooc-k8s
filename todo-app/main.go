package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
)

//go:embed ui/*.tmpl
var templateFS embed.FS

type PageData struct {
	Title   string
	Message string
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Error: env PORT was unset")
	}

	addr := ":" + port
	http.HandleFunc("/", helloHandler)

	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:   "Hello from Todo-App",
		Message: "This is from Exercise: 1.5",
	}

	tmpl, err := template.ParseFS(templateFS, "ui/index.tmpl")
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}
