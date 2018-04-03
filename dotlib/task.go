package dotlib

var (
	// DryRun ...
	DryRun bool
	// Verbose ...
	Verbose bool
)

// Task interface
type Task interface {
	// Register(interface{}) error
	Install() error
	Remove() error
}
