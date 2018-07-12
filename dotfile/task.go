package dotfile

import (
	"fmt"
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
	Do(string) error
}

func do(t Task, a string) (err error) {
	switch a {
	case "Install":
		err = t.Install()
		break
	case "Remove":
		err = t.Remove()
		break
	default: // Unhandled action
		return fmt.Errorf("Unknown task function: %s", a)
	}
	return
}
