package dotfile

import (
	"os"
)

// CreateDir creates a new directory named path along with necessary parents,
// and permission bits are used for all directories that os.MkdirAll creates.
// If there is an error, it will be of type *PathError.
// If path is already a directory, os.MkdirAll does nothing
func CreateDir(path string, mode os.FileMode) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if err == nil {
		if fi.IsDir() {
			return false, nil
		}
	}
	if DryRun {
		return true, nil
	}
	return true, os.MkdirAll(path, mode)
}

// RemoveDir removes the named path and fails if it is a file.
// If there is an error, it will be of type *PathError.
func RemoveDir(path string) (bool, error) {
	_, err := os.Stat(path)
	switch {
	case err != nil && os.IsExist(err):
		return false, err
	case err != nil && os.IsNotExist(err):
		return false, nil
	// case fi != nil && !fi.IsDir():
	// 	return false, &os.PathError{"dir", path, syscall.ENOTDIR}
	// case fi == nil:
	// 	return false, nil
	}
	if DryRun {
		return true, nil
	}
	err = os.Remove(path)
	if err != nil {
		return false, err
	}
	return true, err
}

func ReadDir(path string) ([]os.FileInfo, error) {
	d, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer d.Close()
	di, err := d.Readdir(-1)
	return di, err
}
