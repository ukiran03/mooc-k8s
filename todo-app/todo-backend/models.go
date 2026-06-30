package main

import (
	"database/sql"
	"fmt"
	"log"
	"slices"
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
	stmt := `INSERT INTO tasks (title, state) VALUES (?, ?)`
	result, err := m.DB.Exec(stmt, title, state)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *TaskModel) Update(id int, state TaskState) error {
	stmt := `UPDATE tasks SET state = ?, where id = ?`
	result, err := m.DB.Exec(stmt, state, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf(
			"could not update: no task found with ID %d", id,
		)
	}
	return nil
}

func (m *TaskModel) Delete(id int) error {
	stmt := `DELETE FROM tasks WHERE id = ?`
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
