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
	if err := d.Status(); err != nil {
		if err == ErrAlreadyExist {
			return nil
		}
		return err
	}
	if err := os.MkdirAll(d.Path, 0755); err != nil {
		return err
	}
	return nil
}

// Undo task
func (d *Dir) Undo() error {
	err := d.Status()
	if err != nil && err != ErrAlreadyExist {
	}
	// if err := d.Status(); err != ErrAlreadyExist {
	// 	return err
	// }
	if err == ErrAlreadyExist {
		ok, err := dirIsEmpty(d.Path)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Fprintf(os.Stderr, "%s: directory not empty\n", d.Path)
			return ErrSkip
		}
		// TODO dirOpts.empty
	}
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
