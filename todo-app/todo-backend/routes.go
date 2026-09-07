package main

import "net/http"

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
	mux.HandleFunc("PUT /api/tasks/{id}/done", app.markDoneTask)

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
