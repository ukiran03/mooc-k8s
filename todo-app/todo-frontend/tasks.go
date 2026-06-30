package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TaskState int

const (
	StateTodo TaskState = iota
	StateDone
)

type Task struct {
	ID    int       `json:"id"`
	Title string    `json:"title"`
	State TaskState `json:"state"`
}

var myClient = &http.Client{
	Timeout: 10 * time.Second,
}

func fetchTasksFromBackend(url string) (*[]Task, error) {
	resp, err := myClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	var tasks []Task
	err = json.NewDecoder(resp.Body).Decode(&tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}
	return &tasks, nil
}
