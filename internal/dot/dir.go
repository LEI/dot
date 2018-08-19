package dot

import (
	"fmt"
	"io"
	"os"
)

var (
	defaultDirMode os.FileMode = 0755
)

// Dir task
type Dir struct {
	Task `mapstructure:",squash"` // Action, If, OS
	Path string
}

func (d *Dir) String() string {
	return d.Path
}

// Type task name
func (d *Dir) Type() string {
	return "dir"
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
		return ErrAlreadyExist // &TaskError{"check dir", d, ErrAlreadyExist}
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
	if err := d.Status(); err != nil {
		// terr, ok := err.(*TaskError)
		// if !ok {
		// 	return err
		// }
		// switch terr {}
		switch err {
		case ErrAlreadyExist, ErrSkip:
			return nil
		default:
			return err
		}
	}
	if err := os.MkdirAll(d.Path, defaultDirMode); err != nil {
		return err // &TaskError{"mkdir", d, err}
	}
	return nil
}

// Undo task
func (d *Dir) Undo() error {
	if err := d.Status(); err != nil {
		// terr, ok := err.(*TaskError)
		// if !ok {
		// 	return err
		// }
		// switch terr {
		switch err {
		case ErrAlreadyExist:
		// continue
		case ErrSkip:
			return nil
		default:
			return err
		}
	}
	ok, err := dirIsEmpty(d.Path)
	if err != nil {
		return err // &DirError{"remove", d.Path, err}
	}
	if !ok {
		return &TaskError{"undo dir", d, ErrNotEmpty}
	}
	// TODO dirOpts.empty
	if err := os.Remove(d.Path); err != nil {
		return err // &TaskError{"rmdir", d, err}
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
