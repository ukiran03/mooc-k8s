package main

import "fmt"

type ErrNoTaskFound struct {
	ID int
}

func (e *ErrNoTaskFound) Error() string {
	return fmt.Sprintf("no task found with ID %d", e.ID)
}
