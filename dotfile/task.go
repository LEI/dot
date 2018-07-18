package dotfile

import (
	"fmt"
	"os"
	// "path/filepath"
)

var (
	// DryRun ...
	DryRun bool

	// Verbose ...
	Verbose int

	// FileMode ...
	FileMode os.FileMode = 0644

	// DirMode ...
	DirMode os.FileMode = 0755
)

// Task interface
type Task interface {
	String() string
	// Register(interface{}) error
	Install() error
	Remove() error
	Do(string) error
}

func do(t Task, a string) (err error) {
	switch a {
	case "Install":
		err = t.Install()
	case "Remove":
		err = t.Remove()
	default: // Unhandled action
		return fmt.Errorf("Unknown task function: %s", a)
	}
	return
}
