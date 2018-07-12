package dotfile

import (
	"fmt"
	"os"
	"path/filepath"
)

// CreateDir ...
func CreateDir(dir string) (bool, error) {
	fi, err := os.Stat(dir)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if err == nil && fi.IsDir() {
		return false, nil
	}
	// if Verbose {}
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
	if err == nil && !fi.IsDir() {
		return false, nil
	}
	// if Verbose {}
	if DryRun {
		return true, nil
	}
	// RemoveAll for recursive
	return true, os.Remove(dir)
}

// Cache directories for dry-run mode output
var dirs = map[string]map[string]bool{
	"created": map[string]bool{},
	"removed": map[string]bool{},
}

func createBaseDir(t string) error {
	t = filepath.Dir(t)
	if dirs["created"][t] {
		return nil
	}
	changed, err := CreateDir(t)
	if err != nil {
		return err
	}
	dirs["created"][t] = true
	prefix := "# "
	if changed {
		prefix = ""
	}
	fmt.Printf("%smkdir -p %s\n", prefix, t)
	return nil
}

func removeBaseDir(t string) error {
	// TODO: check if empty and not $HOME?
	t = filepath.Dir(t)
	if dirs["removed"][t] {
		return nil
	}
	changed, err := RemoveDir(t)
	if err != nil {
		return err
	}
	dirs["removed"][t] = true
	prefix := "# "
	if changed {
		prefix = ""
	}
	fmt.Printf("%srmdir %s\n", prefix, t)
	return nil
}
