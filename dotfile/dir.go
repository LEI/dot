package dotfile

import (
	"os"
)

// NewDir
func NewDir(path string, mode os.FileMode) (*File, error) {
	f := NewFile(path)
	err := f.CreateDir(mode)
	return f, err
}

// CreateDir creates a new directory named path along with necessary parents,
// and permission bits are used for all directories that os.MkdirAll creates.
// If there is an error, it will be of type *PathError.
// If path is already a directory, os.MkdirAll does nothing
func (f *File) CreateDir(mode os.FileMode) error {
	return os.MkdirAll(f.path, mode)
}

// RemoveDir removes the named path and fails if it is a file.
// If there is an error, it will be of type *PathError.
func (f *File) RemoveDir() error {
	_, err := f.Stat()
	if err != nil && os.IsExist(err) {
		return err
	}
	if err != nil && os.IsNotExist(err) {
		return nil
	}
	return os.Remove(f.path)
}
