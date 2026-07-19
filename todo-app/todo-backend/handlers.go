package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"unicode/utf8"
)

func (app *backend) getTasks(w http.ResponseWriter, r *http.Request) {
	data, err := app.tasks.GetTasks()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.writeJSON(w, http.StatusOK, data)
}

func (app *backend) createTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var t Task
	if err := app.readJSON(w, r, &t); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	switch length := utf8.RuneCountInString(t.Title); {
	case length == 0:
		app.serverError(w, r, fmt.Errorf("title cannot be empty"))
		return
	case length > 100:
		app.serverError(
			w,
			r,
			fmt.Errorf("title is too long (maximum 100 characters)"),
		)
		return
	}
	_, err := app.tasks.Insert(t.Title, t.State)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *backend) readJSON(
	w http.ResponseWriter,
	r *http.Request,
	dst any,
) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // Disallow unknown fields

	err := dec.Decode(dst)
	if err != nil {
		return err
	}
	return nil
}

func (app *backend) writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("Error encoding JSON: %v", err) // to debug
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *backend) serverError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri)

	http.Error(w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
}
