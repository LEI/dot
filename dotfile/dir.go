package dotfile

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	// RemoveEmptyDirs ...
	RemoveEmptyDirs bool

	// ErrDirShouldExist ...
	ErrDirShouldExist = fmt.Errorf("skip directory")
)

// CreateDir ...
func CreateDir(dir string) (bool, error) {
	fi, err := os.Stat(dir)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if dir == homeDir {
		return false, ErrDirShouldExist
	}
	if err == nil && fi.IsDir() {
		fmt.Printf("# mkdir -p %s\n", dir)
		return false, nil
	}
	fmt.Printf("mkdir -p %s\n", dir)
	if DryRun {
		return true, nil
	}
	return true, os.MkdirAll(dir, DirMode)
}

// RemoveDir ...
func RemoveDir(dir string) (bool, error) {
	fi, err := os.Stat(dir)
	if err != nil && !os.IsExist(err) {
		return false, err
	}
	if dir == homeDir {
		return false, ErrDirShouldExist
	}
	if err == nil && !fi.IsDir() {
		fmt.Printf("# rmdir %s\n", dir)
		return false, nil
	}
	fmt.Printf("rmdir %s\n", dir)
	if DryRun {
		return true, nil
	}
	// RemoveAll for recursive
	return true, os.Remove(dir)
}

// Cache directories for dry-run mode output
var dirs = map[string]map[string]bool{
	"created": {},
	"removed": {},
}

func createBaseDir(t string) error {
	t = filepath.Dir(t)
	if dirs["created"][t] {
		return nil
	}
	_, err := CreateDir(t)
	if err != nil {
		if err == ErrDirShouldExist {
			return nil
		}
		return err
	}
	dirs["created"][t] = true
	return nil
}

func removeBaseDir(t string) error {
	// TODO: check if empty and not $HOME?
	t = filepath.Dir(t)
	if dirs["removed"][t] {
		return nil
	}
	_, err := RemoveDir(t)
	if err != nil {
		if err == ErrDirShouldExist {
			return nil
		}
		return err
	}
	dirs["removed"][t] = true
	return nil
}
