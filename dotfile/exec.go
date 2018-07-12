package dotfile

import (
	"fmt"
)

// ExecTask struct
type ExecTask struct {
	Cmd string
}

// Do ...
func (t *ExecTask) Do(a string) error {
	return do(t, a)
}

// Install copy
func (t *ExecTask) Install() error {
	fmt.Printf("%s\n", t.Cmd)
	return nil
}

// Remove copy
func (t *ExecTask) Remove() error {
	fmt.Printf("%s\n", t.Cmd)
	return nil
}
