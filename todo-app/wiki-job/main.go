package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Task struct {
	Title string `json:"title"`
	State int    `json:"state"`
}

var fallbackURL = "http://todo-backend-svc"

func main() {
	wikiURL := getWikiLink()
	data := Task{
		Title: "Read" + wikiURL,
		State: 0,
	}

	jsData, _ := json.Marshal(data)
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		log.Printf("env BACKEND_URL was unset, using default: %q", fallbackURL)
		backendURL = fallbackURL
	}
	log.Print("Started Job")

	resp, err := http.Post(
		backendURL+"/api/tasks",
		"application/json",
		bytes.NewBuffer(jsData),
	)
	if err != nil {
		log.Fatalf("Failed to connect to backend: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusCreated {
		log.Fatalf("Backend returned non-success code: %d", resp.StatusCode)
	}
	log.Printf("Successfully added todo. Status: %s", resp.Status)
}

func getWikiLink() string {
	client := &http.Client{}

	req, err := http.NewRequest(
		"GET",
		"https://en.wikipedia.org/wiki/Special:Random",
		nil,
	)
	if err != nil {
		log.Println("Error creating request:", err)
		return ""
	}

	// Wikipedia requires a descriptive User-Agent
	req.Header.Set("User-Agent", "MyK3sTodoBot/1.0 (contact: your@email.com)")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error getWikiLink:", err)
		return ""
	}
	defer resp.Body.Close()

	// Now it should follow the redirect to the final article URL
	return resp.Request.URL.String()
}
