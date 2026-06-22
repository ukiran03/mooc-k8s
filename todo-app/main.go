package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

//go:embed ui/*.tmpl
var templateFS embed.FS

var templates *template.Template

// Global image cache instance
var currentImage = &Image{}

const (
	pathname = "./image.jpg"
	url      = "https://placeholdpicsum.dev/photo/category/nature/400/400"
)

type Image struct {
	mu      sync.Mutex
	name    string
	modTime time.Time
}

type PageData struct {
	Title   string
	Message string
	Image   string
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Parse templates once at startup
	var err error
	templates, err = template.ParseFS(templateFS, "ui/index.tmpl")
	if err != nil {
		log.Fatalf("Critical: Could not parse templates: %v", err)
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/image.jpg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, pathname)
	})

	addr := ":" + port
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	isCached, img := GetImage(currentImage)
	log.Printf("Image requested, Cache hit: %t, File: %s", isCached, img.name)

	data := PageData{
		Title:   "Todo App",
		Message: "This is from Exercise: 1.12",
		Image:   "/image.jpg", // the route pointing to our handler above
	}

	err := templates.ExecuteTemplate(w, "index.tmpl", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func GetImage(img *Image) (bool, *Image) {
	img.mu.Lock()
	defer img.mu.Unlock()

	// Cache Hit check
	if !img.modTime.IsZero() && time.Since(img.modTime) < 10*time.Minute {
		return true, img
	}

	// Cache miss -> Download
	filename, err := downloadImage()
	if err != nil {
		log.Printf("Download failed: %v", err)
		return true, img // Fallback to stale data if download fails
	}

	info, err := os.Stat(filename)
	if err != nil {
		log.Printf("Stat failed: %v", err)
		return true, img
	}

	img.name = filename
	img.modTime = info.ModTime()

	return false, img
}

func downloadImage() (string, error) {
	resp, err := http.Get(
		fmt.Sprintf("%s?t=%d", url, time.Now().UnixMilli()),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	f, err := os.Create(pathname)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", err
	}
	return pathname, nil
}
