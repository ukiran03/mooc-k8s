package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaSQL string

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

	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database forcefully reset to demo state.")
	defer db.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &backend{
		tasks:  &TaskModel{DB: db},
		logger: logger,
	}

	addr := ":" + port
	log.Printf("Backend Server started on port %s", port)
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

func initDB() (*sql.DB, error) {
	dbDir := os.Getenv("DATABASE_DIR")
	if dbDir == "" {
		log.Printf("env DATABASE_DIR is unset")
		dbDir = "."
	}

	// Ensure the directory exists (crucial for K8s volumes)
	if err := os.MkdirAll(dbDir, 0o755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dbDir, "tasks.sqlite3")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	// // Split the schema by semicolon and execute each statement sequentially
	// queries := strings.Split(schemaSQL, ";")
	// for _, query := range queries {
	// 	query = strings.TrimSpace(query)
	// 	if query == "" {
	// 		continue
	// 	}

	// }

	if _, err := db.Exec(schemaSQL); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
