package dotfile

import (
	"os"
)

var (
	// DryRun ...
	DryRun bool

	// FileMode ...
	FileMode os.FileMode = 0644

	// DirMode ...
	DirMode os.FileMode = 0755
)

// Task interface
type Task interface {
	// Register(interface{}) error
	Install() error
	Remove() error
}
