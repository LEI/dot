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
	s := fmt.Sprintf("%s %s", name, strings.Join(args[:], " "))
	fmt.Println(s)
	if DryRun {
		// fmt.Println("DRY-RUN, unexpected behavior may occur.")
		return nil
	}
	c := exec.Command(name, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
