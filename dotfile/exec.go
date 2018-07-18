package dotfile

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var (
	// Shell ...
	Shell = "bash"
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
	return execute(Shell, "-c", t.Cmd)
}

// Remove copy
func (t *ExecTask) Remove() error {
	return execute(Shell, "-c", t.Cmd)
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

const defaultStatusFailed = 1

// ShellExec ...
func ShellExec(c string) (stdout, stderr string, status int) {
	args := []string{"-c", c}
	stdout, stderr, status = ExecCommand(Shell, args...)
	return
}

// ExecCommand ...
func ExecCommand(name string, args ...string) (stdout, stderr string, status int) {
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if ws, ok := exitError.Sys().(syscall.WaitStatus); ok {
				status = ws.ExitStatus()
			} else {
				fmt.Fprintf(os.Stderr, "Could not get exit status: %+v\n", args)
				status = defaultStatusFailed
			}
		} else {
			fmt.Fprintf(os.Stderr, "Could not get exit error: %+v\n", args)
			status = defaultStatusFailed
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		if ws, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
			status = ws.ExitStatus()
		} else {
			fmt.Fprintf(os.Stderr, "Could not get processed status: %+v\n", args)
			status = defaultStatusFailed
		}
	}
	return
}

// AskConfirmation ...
func AskConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/n]: ", s)
		res, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not read input from stdin: %s\n", err)
			os.Exit(1)
		}
		res = strings.ToLower(strings.TrimSpace(res))
		if res == "y" || res == "yes" {
			return true
		} else if res == "n" || res == "no" {
			return false
		}
	}
}
