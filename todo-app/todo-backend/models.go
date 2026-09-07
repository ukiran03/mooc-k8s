package main

import (
	"database/sql"
	"errors"
	"log"
	"slices"
)

type TaskState int

const (
	StateTodo TaskState = iota
	StateDone
)

func (s TaskState) String() string {
	switch s {
	case StateTodo:
		return "Created New Task"
	case StateDone:
		return "Task Marked Done"
	default:
		return "Unknown State"
	}
}

type Task struct {
	ID    int       `json:"id"`
	Title string    `json:"title"`
	State TaskState `json:"state"`
}

type TaskModel struct {
	DB *sql.DB
}

func (m *TaskModel) GetTasks() ([]Task, error) {
	stmt := `SELECT id, title, state FROM tasks`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]Task, 0, 10)
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.State)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Newest first
	slices.Reverse(tasks)
	return tasks, nil
}

func (m *TaskModel) Insert(title string, state TaskState) (int, error) {
	stmt := `INSERT INTO tasks (title, state) VALUES ($1, $2) RETURNING id`

	var id int
	err := m.DB.QueryRow(stmt, title, state).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *TaskModel) Update(id int, state TaskState) (string, error) {
	stmt := `UPDATE tasks SET state = $1 WHERE id = $2 RETURNING title`

	var title string
	err := m.DB.QueryRow(stmt, state, id).Scan(&title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", &ErrNoTaskFound{ID: id}
		}
		return "", err
	}

	return title, nil
}

func (m *TaskModel) Delete(id int) error {
	stmt := `DELETE FROM tasks WHERE id = $1`
	result, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("Deleted %d row(s)", rowsAffected)
	return nil
}
