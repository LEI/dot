package dotfile

import (
	"fmt"
	"os"
	"strings"
)

type File struct {
	path string
	stat os.FileInfo
}

// NewFile
func NewFile(path string) *File {
	return &File{path: path}
}

func (f *File) String() string {
	return fmt.Sprintf("File[%s]", f.path)
}

func (f *File) Path() string {
	return f.path
}

func (f *File) SetPath(path string) {
	f.path = path
}

func (f *File) Info() (os.FileInfo, error) {
	if f.stat != nil {
		return f.stat, nil
	}
	fi, err := f.Stat()
	f.stat = fi
	return f.stat, err
}

func (f *File) Stat() (os.FileInfo, error) {
	fi, err := os.Stat(f.path)
	f.stat = fi
	return f.stat, err
}

func (f *File) Exists() (bool, error) {
	fi, err := f.Stat()
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

func (f *File) IsDir() (bool, error) {
	fi, err := f.Info()
	if err != nil {
		return false, err
	}
	if fi != nil && fi.IsDir() {
		return true, nil
	}
	return false, err
}

func (f *File) Base() string {
	return Base(f.path)
}

func (f *File) Match(patterns ...string) (bool, error) {
	matched, err := Match(f.path, patterns...)
	return matched, err
}

func (f *File) BaseMatch(patterns ...string) (bool, error) {
	matched, err := Match(f.Base(), patterns...)
	return matched, err
}

func (f *File) Replace(old, new string) string {
	return strings.Replace(f.path, old, new, 1)
}
