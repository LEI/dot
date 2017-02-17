package fileutil

import (
	"os"
)

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func MakeDir(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	return nil
}

func MakeDirs(paths []string) error {
	for _, path := range paths {
		err := MakeDir(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func RemoveDir(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}
