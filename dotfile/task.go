package dotfile

import (
	"fmt"
	"os"
	"path/filepath"
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

// SplitPath ...
func SplitPath(s string) (src, dst string) {
	parts := filepath.SplitList(s)
	switch len(parts) {
	case 1:
		src = s
		break
	case 2:
		src = parts[0]
		dst = parts[1]
		break
	default:
		fmt.Println("Unhandled path spec", src)
		os.Exit(1)
	}
	// src = s
	// if strings.Contains(src, ":") {
	// 	parts := strings.Split(src, ":")
	// 	if len(parts) == 2 {
	// 		src = parts[0]
	// 		dst = parts[1]
	// 	} else {
	// 		fmt.Println("Unhandled path spec", src)
	// 		os.Exit(1)
	// 	}
	// }
	return src, dst
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
