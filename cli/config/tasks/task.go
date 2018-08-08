package tasks

// Task struct
type Task struct {
	Tasker
	execute bool
}

// Tasker interface
type Tasker interface {
	Check() error
	Execute() error
}
