package tasks

import (
	"fmt"
)

var (
	// ErrSkip ...
	ErrSkip = fmt.Errorf("skip")
)

// Task struct
type Task struct {
	Tasker
	// execute bool
	toDo bool
}

// DoInstall ...
func (t *Task) DoInstall() bool {
	return !t.toDo
}

// DoRemove ...
func (t *Task) DoRemove() bool {
	return t.toDo
}

// Tasker interface
type Tasker interface {
	Check() error
	// Execute() error
	Install() error
	Remove() error
}
