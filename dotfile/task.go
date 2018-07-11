package dotfile

import (
	"os"
)

var (
	// DryRun ...
	DryRun bool

	// FileMode default
	FileMode os.FileMode = 0644
)

// Task interface
type Task interface {
	// Register(interface{}) error
	Install() error
	Remove() error
}
