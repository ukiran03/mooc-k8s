package main

import (
	"errors"
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
		http.Error(w, "title cannot be empty", http.StatusBadRequest)
		app.logger.Error("title cannot be empty", "length", length)
		return
	case length > 100:
		http.Error(
			w,
			"title is too long (maximum 100 characters)",
			http.StatusBadRequest,
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

	// Assign the newly generated DB ID to the task before broadcasting
	t.ID = id

	msg, err := NewBroadcastMsg(t, t.State).Bytes()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if err := app.natscon.Publish(app.natsSub, msg); err != nil {
		app.logger.Error("failed to publish nats message", "error", err)
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

	title, err := app.tasks.Update(id, StateDone)
	if err != nil {
		if _, ok := errors.AsType[*ErrNoTaskFound](err); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		app.logger.Error("failed to update task", "error", err)
		app.serverError(w, r, err)
		return
	}

	msg, err := NewBroadcastMsg(Task{ID: id, Title: title}, StateDone).Bytes()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if err := app.natscon.Publish(app.natsSub, msg); err != nil {
		app.logger.Error("failed to publish nats message", "error", err)
	}

	app.logger.Info("task marked done", "id", id)
	w.WriteHeader(http.StatusAccepted)
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
