package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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
