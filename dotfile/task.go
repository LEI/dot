/*
Package dotfile tasks
*/
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
	Status() bool
	List() (string, error)
	Install() (string, error)
	Remove() (string, error)
	Do(string) (string, error)
}

func do(t Task, a string) (str string, err error) {
	switch a {
	case "List":
		str, err = t.List()
	case "Install":
		str, err = t.Install()
	case "Remove":
		str, err = t.Remove()
	default: // Unhandled action
		err = fmt.Errorf("unknown task function: %s", a)
	}
	return
}
