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
	s := d.Path
	switch d.GetAction() {
	case "install":
		s = fmt.Sprintf("mkdir -p %s", tildify(d.Path))
	case "remove":
		s = fmt.Sprintf("rmdir %s", tildify(d.Path))
	}
	return s // d.Path
}

// Status check task
func (d *Dir) Status() error {
	if !dirExists(d.Path) {
		return nil
	}
	if d.GetAction() == "remove" {
		empty, err := dirIsEmpty(d.Path)
		if err != nil {
			return err // &DirError{"remove", d.Path, err}
		}
		if !empty {
			// ErrExist would indicate that the directory should be removed
			return &os.PathError{Op: "rmdir", Path: d.Path, Err: ErrDirNotEmpty}
		}
	}
	return ErrExist // &OpError{"check dir", d, ErrExist}
	// fi, err := os.Stat(d.Path)
	// if err != nil && os.IsExist(err) {
	// 	return err
	// }
	// if fi != nil && fi.IsDir() {
	// 	return ErrExist // fmt.Errorf("%s: directory exists", d.Path)
	// }
	// return nil
}

// Do task
func (d *Dir) Do() error {
	if err := d.Status(); err != nil {
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
	if err := os.MkdirAll(d.Path, defaultDirMode); err != nil {
		return err // &OpError{"mkdir", d, err}
	}
	return nil
}

// Undo task unless the directory is not empty.
func (d *Dir) Undo() error {
	if err := d.Status(); err != nil {
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
	if err := os.Remove(d.Path); err != nil {
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
