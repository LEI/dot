package dotfile

import (
	// "fmt"
	"os"
)

func Exists(path string) (bool, error) {
	fi, err := os.Stat(path)
	switch {
	case err != nil && os.IsExist(err), err == nil && fi != nil:
		return true, err
	case err != nil && os.IsNotExist(err), fi == nil:
		return false, nil
	case fi == nil:
		return false, err
	}
	return true, err
}

func IsDir(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if fi != nil && fi.IsDir() {
		return true, nil
	}
	return false, err
}
