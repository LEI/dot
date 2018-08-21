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
	return fmt.Sprintf("cp %s %s", tildify(c.Source), tildify(c.Target))
}

// UndoString string
func (c *Copy) UndoString() string {
	return fmt.Sprintf("rm %s", tildify(c.Target))
}

// Status check task
func (c *Copy) Status() error {
	exists, err := copyExists(c.Source, c.Target)
	if err != nil {
		return err
	}
	if exists {
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
		// fmt.Errorf("%s: no such file to copy to %s", src, dst)
		return false, &os.PathError{Op: "copy", Path: src, Err: ErrNotExist}
	}
	if !exists(dst) {
		// Stop here if the target does not exist
		return false, nil
	}
	return fileCompare(src, dst)
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

// fileCompare TODO read in chunks
func fileCompare(p1, p2 string) (bool, error) {
	a, err := ioutil.ReadFile(p1)
	if err != nil {
		return false, err
	}
	b, err := ioutil.ReadFile(p2)
	if err != nil {
		return false, err
	}
	return bytes.Compare(a, b) == 0, nil
}
