package main

import (
	_ "embed"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

type backend struct {
	tasks     *TaskModel
	logger    *slog.Logger
	broken    bool
	brokenMux sync.RWMutex
}

func main() {
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		fmt.Println("env PORT was unset\nUsing Port 3001 as Backend port")
		port = "3001"
	}

	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = initialseTable(db)
	if err != nil {
		log.Printf("Error initialseTable: %v", err)
		return
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &backend{
		tasks:  &TaskModel{DB: db},
		logger: logger,
	}

	addr := ":" + port
	logger.Info("Starting Todo-App Backend", "address", addr)
	log.Fatal(http.ListenAndServe(addr, app.routes()))
}

func (app *backend) routes() http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint with database connectivity
	mux.HandleFunc("GET /health", app.healthCheck)
	mux.HandleFunc("GET /api/health", app.healthCheck)
	mux.HandleFunc("GET /", app.healthCheck)

	// Break/fix endpoints for testing (accessible via ingress)
	mux.HandleFunc("POST /api/break", app.breakApp)
	mux.HandleFunc("POST /api/fix", app.fixApp)

	// API endpoints
	mux.HandleFunc("GET /api/tasks", app.getTasks)
	mux.HandleFunc("POST /api/tasks", app.createTask)

	return app.enableCORS(mux)
}

func (b *backend) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle the Pre-flight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Pass the request to the mux (where it will match GET or POST)
		next.ServeHTTP(w, r)
	})
}

func (app *backend) healthCheck(w http.ResponseWriter, r *http.Request) {
	app.brokenMux.RLock()
	isBroken := app.broken
	app.brokenMux.RUnlock()

	if isBroken {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Service is broken"))
		return
	}

	// Check database connectivity
	err := app.tasks.DB.Ping()
	if err != nil {
		app.logger.Error("Database health check failed", "error", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Database connection failed"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (app *backend) breakApp(w http.ResponseWriter, r *http.Request) {
	app.brokenMux.Lock()
	app.broken = true
	app.brokenMux.Unlock()
	app.logger.Warn("App has been broken")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("App broken"))
}

func (app *backend) fixApp(w http.ResponseWriter, r *http.Request) {
	app.brokenMux.Lock()
	app.broken = false
	app.brokenMux.Unlock()
	app.logger.Info("App has been fixed")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("App fixed"))
}
