package dot

import (
	"fmt"
	"os"
	"path/filepath"
)

// Files task list
type Files []*Copy

// Copy task
type Copy struct {
	Task
	Source string
	Target string
}

func (c *Copy) String() string {
	return fmt.Sprintf("%s:%s", c.Source, c.Target)
}

// DoString string
func (c *Copy) DoString() string {
	return fmt.Sprintf("cp %s %s", c.Source, c.Target)
}

// UndoString string
func (c *Copy) UndoString() string {
	return fmt.Sprintf("rm %s", c.Target)
}

// Prepare task
func (c *Copy) Prepare(target string) error {
	// if !filepath.IsAbs(c.Source) {
	// 	c.Target = filepath.Join(source, c.Source)
	// }
	if !filepath.IsAbs(c.Target) {
		c.Target = filepath.Join(target, c.Target)
	}
	return nil
}

// Status check task
func (c *Copy) Status() error {
	if fileExists(c.Target) {
		return ErrAlreadyExist
	}
	return nil
}

// Do task
func (c *Copy) Do() error {
	if err := c.Status(); err != nil {
		if err == ErrAlreadyExist {
			return nil
		}
		return err
	}
	fmt.Println("todo", c)
	return nil
}

// Undo task
func (c *Copy) Undo() error {
	fmt.Println("toundo", c)
	return nil
	// return os.Remove(c.Target)
}

// fileExists returns true if the file has the same content.
func fileExists(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()
	if _, err := f.Stat(); err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
