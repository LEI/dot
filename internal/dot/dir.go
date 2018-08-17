package dot

import (
	"fmt"
	"os"
)

// Dirs task list
type Dirs []*Dir

// Dir task
type Dir struct {
	Task
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
	// if err := d.Status(); err != ErrAlreadyExist {
	// 	return err
	// }
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
