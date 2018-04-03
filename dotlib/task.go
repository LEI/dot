package dotlib

var (
	DryRun bool
	Verbose bool
)

// Task interface
type Task interface {
	// Register(interface{}) error
	Install() error
	Remove() error
}