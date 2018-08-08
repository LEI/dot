package system

import (
	"os"
)

var (
	// FileMode ...
	FileMode os.FileMode = 0644
)

// Exists checks if a file is present
func Exists(file string) bool {
	_, err := os.Stat(file)
	// return !os.IsNotExist(err)
	return err == nil || os.IsExist(err)
}

// IsDir checks if a file is a directory
func IsDir(file string) (bool, error) {
	fi, err := os.Stat(file)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if fi == nil {
		return false, nil
	}
	return fi.IsDir(), nil
}
