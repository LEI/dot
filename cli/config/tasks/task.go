package tasks

// Task interface
type Task interface {
	Check() error
	Execute() error
}
