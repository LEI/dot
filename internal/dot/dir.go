package dot

import (
	"fmt"
	"io"
	"os"
)

// Dir task
type Dir struct {
	Task `mapstructure:",squash"` // Action, If, OS
	Path string
}

// NewDir task
func NewDir(s string) *Dir {
	return &Dir{Path: s}
}

func (t *Dir) String() string {
	s := t.Path
	switch Action {
	case "install":
		s = fmt.Sprintf("mkdir -p %s", tildify(t.Path))
	case "remove":
		s = fmt.Sprintf("rmdir %s", tildify(t.Path))
	}
	return s // t.Path
}

// Init task
func (t *Dir) Init() error {
	// ...
	return nil
}

// Status check task
func (t *Dir) Status() error {
	if !dirExists(t.Path) {
		return nil
	}
	if Action == "remove" {
		empty, err := dirIsEmpty(t.Path)
		if err != nil {
			return err // &DirError{"remove", t.Path, err}
		}
		if !empty {
			// ErrExist would indicate that the directory should be removed
			return &os.PathError{Op: "rmdir", Path: t.Path, Err: ErrDirNotEmpty}
		}
	}
	return ErrExist // &OpError{"check dir", t, ErrExist}
	// fi, err := os.Stat(t.Path)
	// if err != nil && os.IsExist(err) {
	// 	return err
	// }
	// if fi != nil && fi.IsDir() {
	// 	return ErrExist // fmt.Errorf("%s: directory exists", t.Path)
	// }
	// return nil
}

// Do task
func (t *Dir) Do() error {
	if err := t.Status(); err != nil {
		// terr, ok := err.(*OpError)
		// if !ok {
		// 	return err
		// }
		// switch terr {}
		switch err {
		case ErrExist, ErrSkip:
			return nil
		default:
			return err
		}
	}
	if err := os.MkdirAll(t.Path, defaultDirMode); err != nil {
		return err // &OpError{"mkdir", t, err}
	}
	return nil
}

// Undo task unless the directory is not empty.
func (t *Dir) Undo() error {
	if err := t.Status(); err != nil {
		// terr, ok := err.(*OpError)
		// if !ok {
		// 	return err
		// }
		// switch terr {
		switch err {
		case ErrExist:
			// continue
		case ErrSkip:
			return nil
		default:
			return err
		}
	}
	if err := os.Remove(t.Path); err != nil {
		return err
	}
	return nil
}

// dirExists returns true if the name exists and is a directory.
func dirExists(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func dirIsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()
	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}
