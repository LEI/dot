package dot

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// Copy task
type Copy struct {
	Task   `mapstructure:",squash"` // Action, If, OS
	Source string
	Target string
}

func (c *Copy) String() string {
	return fmt.Sprintf("%s:%s", c.Source, c.Target)
}

// Type task name
func (c *Copy) Type() string {
	return "copy"
}

// DoString string
func (c *Copy) DoString() string {
	return fmt.Sprintf("cp %s %s", c.Source, c.Target)
}

// UndoString string
func (c *Copy) UndoString() string {
	return fmt.Sprintf("rm %s", c.Target)
}

// Status check task
func (c *Copy) Status() error {
	ok, err := copyExists(c.Source, c.Target)
	if err != nil {
		return err
	}
	if ok {
		return ErrAlreadyExist
	}
	return nil
}

// Do task
func (c *Copy) Do() error {
	if err := c.Status(); err != nil {
		switch err {
		case ErrAlreadyExist, ErrSkip:
			return nil
		default:
			return err
		}
	}
	in, err := os.Open(c.Source)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(c.Target)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	if err := out.Sync(); err != nil {
		return err
	}
	return nil
}

// Undo task
func (c *Copy) Undo() error {
	if err := c.Status(); err != nil {
		switch err {
		case ErrSkip:
			return nil
		case ErrAlreadyExist:
			// continue
		default:
			return err
		}
	}
	return os.Remove(c.Target)
}

// copyExists returns true if the file source and target have the same content.
func copyExists(src, dst string) (bool, error) {
	if !exists(src) {
		// return ErrIsNotExist
		return false, fmt.Errorf("%s: no such file to copy to %s", src, dst)
	}
	if !exists(dst) {
		// Stop here if the target does not exist
		return false, nil
	}
	a, err := ioutil.ReadFile(src)
	if err != nil {
		return false, err
	}
	b, err := ioutil.ReadFile(dst)
	if err != nil {
		return false, err
	}
	ok := bytes.Compare(a, b) == 0
	return ok, nil
}

// fileExists returns true if the name exists and is a not a directory.
func fileExists(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return !fi.IsDir()
}
