package dot

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
