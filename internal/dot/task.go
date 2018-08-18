package dot

import "os"

// Tasker interface
type Tasker interface {
	String() string
	DoString() string
	UndoString() string
	Status() error
	// Sync() error
	Do() error
	Undo() error
}

// Task struct
type Task struct {
	Tasker
}

// IsOk status
func IsOk(err error) bool {
	return err == ErrAlreadyExist
}

// exists checks if a file is present
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
