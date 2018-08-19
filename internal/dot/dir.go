package dot

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	// ErrDirNotEmpty ...
	ErrDirNotEmpty = errors.New("directory not empty")

	defaultDirMode os.FileMode = 0755
)

// // DirError type
// type DirError struct {
// 	// taskError
// 	Action string
// 	Path   string
// 	Err    error
// 	skip   bool
// }

// func (e *DirError) Error() string {
// 	return e.Action + " " + e.Path + ": " + e.Err.Error()
// }

// Dir task
type Dir struct {
	Task `mapstructure:",squash"` // Action, If, OS
	Path string
}

func (d *Dir) String() string {
	return d.Path
}

// DoString string
func (d *Dir) DoString() string {
	return fmt.Sprintf("mkdir -p %s", d.Path)
}

// UndoString string
func (d *Dir) UndoString() string {
	return fmt.Sprintf("rmdir %s", d.Path)
}

// Status check task
func (d *Dir) Status() error {
	if dirExists(d.Path) {
		return ErrAlreadyExist
	}
	// fi, err := os.Stat(d.Path)
	// if err != nil && os.IsExist(err) {
	// 	return err
	// }
	// if fi != nil && fi.IsDir() {
	// 	return ErrAlreadyExist // fmt.Errorf("%s: directory exists", d.Path)
	// }
	return nil
}

// Do task
func (d *Dir) Do() error {
	// if err := d.Status(); err != nil {
	// 	if err == ErrAlreadyExist {
	// 		return nil
	// 	}
	// 	return err
	// }
	if err := d.Status(); err != nil {
		switch err {
		case ErrAlreadyExist, ErrSkip:
			return nil
		default:
			return err
		}
	}
	// if err := os.MkdirAll(d.Path, defaultDirMode); err != nil {
	// 	return err
	// }
	// return nil
	return os.MkdirAll(d.Path, defaultDirMode)
}

// Undo task
func (d *Dir) Undo() error {
	if err := d.Status(); err != nil {
		switch err {
		case ErrSkip:
			return nil
		case ErrAlreadyExist:
			// continue
		default:
			return err
		}
	}
	// if err := d.Status(); err != ErrAlreadyExist {
	// 	return err
	// }
	ok, err := dirIsEmpty(d.Path)
	if err != nil {
		return err // &DirError{"remove", d.Path, err}
	}
	if !ok {
		// return &taskError{Detail: d, Format: "", Message: ""}
		return ErrSkip // &DirError{"remove", d.Path, ErrDirNotEmpty, true}
	}
	// TODO dirOpts.empty
	return os.Remove(d.Path)
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
