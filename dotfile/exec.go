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
	Task
}

// Do ...
func (t *ExecTask) Do(a string) error {
	return do(t, a)
}

// Install copy
func (t *ExecTask) Install() error {
	name, args := t.Cmd[0], t.Cmd[1:]
	fmt.Println(name, args)
	return nil
	// return execute(name, args...)
}

// Remove copy
func (t *ExecTask) Remove() error {
	name, args := t.Cmd[0], t.Cmd[1:]
	fmt.Println(name, args)
	return nil
	// return execute(name, args...)
}

var execWarned bool

func execute(name string, args ...string) error {
	if DryRun && !execWarned {
		fmt.Println("DRY-RUN, unexpected behavior may occur.")
		execWarned = true
	}
	fmt.Printf("%s %s\n", name, strings.Join(args[:], " "))
	if DryRun {
		return nil
	}
	c := exec.Command(name, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
