package dotfile

import (
	"fmt"
	"os"
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
