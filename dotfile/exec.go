package dotfile

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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

func execute(name string, args ...string) error {
	fmt.Printf("$ %s %s\n", name, strings.Join(args[:], " "))
	c := exec.Command(name, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
