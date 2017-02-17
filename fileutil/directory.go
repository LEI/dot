package fileutil

import (
	"os"
)

func MakeDir(path string) error {
	fi, err := os.Stat(path)
	if err != nil && os.IsExist(err) {
		return err
	}
	if fi != nil {
		return nil
	}
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	return nil
}

func RemoveDir(path string) error {
	fi, err := os.Stat(path)
	if err != nil && os.IsExist(err) {
		return err
	}
	if err != nil && os.IsNotExist(err) || fi == nil {
		return nil
	}
	err = os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}
