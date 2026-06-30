package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

// Embed the entire ui directory (templates, js, css, etc.)
//
//go:embed all:ui
var uiFS embed.FS

// Create a sub-filesystem rooted at the "ui" folder
var templateFS, _ = fs.Sub(uiFS, "ui")

const (
	webImagePath = "/images/image.jpg"
	imagePath    = "./images/image.jpg"

	imageURL = "https://placeholdpicsum.dev/photo/category/nature/250/250"

	pageTitle   = "Todo App"
	pageMessage = "This is from Exercise: 2.2"
)

var (
	currentImage = &Image{} // Global image cache instance

	// The full URL used by the Go backend container to fetch data
	// internally within the cluster (e.g., http://todo-backend-svc/api/tasks)
	backendURL string

	// The relative path injected into <form data-url="..."> for the browser
	// to make client-side requests routed through the Ingress controller
	backendRoute = "/api/tasks"
)

type frontend struct {
	image     string
	logger    *slog.Logger
	templates *template.Template
}

func main() {
	backendEnvUrl := os.Getenv("BACKEND_URL")
	if backendEnvUrl == "" {
		backendEnvUrl = "http://localhost:3001" // Local fallback
	}
	backendEnvUrl = strings.TrimSuffix(backendEnvUrl, "/")
	backendURL = backendEnvUrl + backendRoute

	log.Printf("[DEBUG] Server dynamically configured backend target to: %s", backendURL)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Parse templates once at startup
	templates, err := newTemplateCache()
	if err != nil {
		log.Fatalf("error: could not parse templates: %v", err)
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &frontend{
		image:     imagePath,
		logger:    logger,
		templates: templates,
	}

	addr := ":" + port
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(addr, app.routes()))
}

func (app *frontend) routes() http.Handler {
	mux := http.NewServeMux()

	// Create a file server from your embedded UI filesystem: This maps the
	// "/static/" URL route directly to the "ui/static" folder in your embed
	staticFS, _ := fs.Sub(uiFS, "ui")
	fileServer := http.FileServer(http.FS(staticFS))

	// Do NOT use StripPrefix here if the folder inside staticFS is named "static"
	mux.Handle("/static/", fileServer)

	imageServer := http.FileServer(http.Dir("./images/"))
	mux.Handle("/images/", http.StripPrefix("/images/", imageServer))

	mux.HandleFunc("/{$}", app.homeHandler)
	return mux
}

func (app *frontend) homeHandler(w http.ResponseWriter, r *http.Request) {
	isCached, img := GetImage(currentImage)
	log.Printf("Image requested, Cache hit: %t, File: %s", isCached, img.name)

	tasks, err := fetchTasksFromBackend(backendURL) // this call is done via the cluster
	if err != nil {
		app.backendError(w, r, err)
		return
	}

	data := PageData{
		Title:      pageTitle,
		Message:    pageMessage,
		Image:      webImagePath,
		Tasks:      *tasks,
		BackendURL: backendRoute, // this call is done via the browser
	}
	app.logger.Info(imagePath)

	err = app.templates.ExecuteTemplate(w, "index.tmpl", data)
	if err != nil {
		app.logger.Error("template execution error: " + err.Error())
		app.serverError(w, r, err)
		return
	}
}

func (app *frontend) backendError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method   = r.Method
		uri      = r.URL.RequestURI()
		errorMsg = fmt.Sprintf("Backed Error: %v", err)
	)
	app.logger.Error(errorMsg, "method", method, "uri", uri)
	http.Error(w, errorMsg, http.StatusInternalServerError)
}

func (app *frontend) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(
		w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}
