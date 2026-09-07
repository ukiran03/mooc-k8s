package main

import (
	"fmt"
	"net/http"
	"strconv"
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
		app.logger.Error("failed to decode task JSON", "error", err)
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	switch length := utf8.RuneCountInString(t.Title); {
	case length == 0:
		app.serverError(w, r, fmt.Errorf("title cannot be empty"))
		app.logger.Error("title cannot be empty", "length", length)
		return
	case length > 100:
		app.serverError(
			w,
			r,
			fmt.Errorf("title is too long (maximum 100 characters)"),
		)
		app.logger.Error(
			"title is too long (maximum 100 characters)",
			"length",
			length,
		)
		return
	}

	id, err := app.tasks.Insert(t.Title, t.State)
	if err != nil {
		app.logger.Error("database insertion failed",
			"error", err,
			"title", t.Title,
		)
		app.serverError(w, r, err)
		return
	}

	app.logger.Info("task created successfully", "id", id)
	w.WriteHeader(http.StatusCreated)
}

func (app *backend) markDoneTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Task ID", http.StatusBadRequest)
		return
	}

	err = app.tasks.Update(id, StateDone)
	if err != nil {
		app.logger.Error(err.Error())
		app.serverError(w, r, err)
		return
	}

	app.logger.Info("task marked done", "id", id)
	w.WriteHeader(http.StatusAccepted)
}
