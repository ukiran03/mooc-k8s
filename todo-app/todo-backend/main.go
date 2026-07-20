package main

import (
	_ "embed"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type backend struct {
	tasks  *TaskModel
	logger *slog.Logger
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
