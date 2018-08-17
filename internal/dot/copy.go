package dot

import "fmt"

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

// Status check
func (c *Copy) Status() error {
	return nil
}

// Do task
func (c *Copy) Do(run bool) error {
	return nil
}

// Undo task
func (c *Copy) Undo(run bool) error {
	return nil
}
