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

	execDir string
)

// ExecTask struct
type ExecTask struct {
	Cmd string
	Dir string
	Task
}

// Status exec
func (t *ExecTask) Status() bool {
	return true
}

// Do exec
func (t *ExecTask) Do(a string) (string, error) {
	return do(t, a)
}

// List exec
func (t *ExecTask) List() (string, error) {
	str := fmt.Sprintf("Exec: %s", t.Cmd)
	return str, nil
}

// Install copy
func (t *ExecTask) Install() (string, error) {
	fmt.Println(t.Cmd)
	// return "", executeIn(t.Dir, Shell, "-c", t.Cmd)
	stdout, stderr, status := ExecCommandIn(t.Dir, Shell, "-c", t.Cmd)
	stdout = strings.TrimSuffix(stdout, "\n")
	stderr = strings.TrimSuffix(stderr, "\n")
	if status != 0 || stderr != "" {
		return stdout, fmt.Errorf(stderr)
	}
	return stdout, nil
}

// Remove copy
func (t *ExecTask) Remove() (string, error) {
	fmt.Println(t.Cmd)
	// return "", executeIn(t.Dir, Shell, "-c", t.Cmd)
	stdout, stderr, status := ExecCommandIn(t.Dir, Shell, "-c", t.Cmd)
	stdout = strings.TrimSuffix(stdout, "\n")
	stderr = strings.TrimSuffix(stderr, "\n")
	if status != 0 || stderr != "" {
		return stdout, fmt.Errorf(stderr)
	}
	return stdout, nil
}

var execWarned bool

func execute(name string, args ...string) error {
	if DryRun && !execWarned {
		fmt.Println("DRY-RUN, unexpected behavior may occur.")
		execWarned = true
	}
	// fmt.Printf("%s %s\n", name, strings.Join(args[:], " "))
	if DryRun {
		return nil
	}
	c := exec.Command(name, args...)
	if execDir != "" {
		c.Dir = execDir
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func executeIn(dir, name string, args ...string) error {
	execDir = dir
	err := execute(name, args...)
	execDir = ""
	return err
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
	if execDir != "" {
		cmd.Dir = execDir
	}
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
			// if stderr == "" {
			// 	stderr = err.Error()
			// }
		}
	} else {
		if ws, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
			status = ws.ExitStatus()
		} else {
			fmt.Fprintf(os.Stderr, "Could not get processed status: %+v\n", args)
			status = defaultStatusFailed
		}
	}
	stdout = outbuf.String()
	stderr = errbuf.String() // + stderr
	// fmt.Println("CMD:", args)
	// fmt.Println("OUT:", stdout)
	// fmt.Println("ERR:", stderr)
	return
}

// ExecCommandIn ...
func ExecCommandIn(dir, name string, args ...string) (string, string, int) {
	execDir = dir
	stdout, stderr, status := ExecCommand(name, args...)
	execDir = ""
	return stdout, stderr, status
}

// AskConfirmation ...
func AskConfirmation(s string) (ret bool) {
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
			ret = true
			break
		} else if res == "n" || res == "no" {
			ret = false
			break
		}
	}
	// FIXME: no new line if enter is pressed before the last fmt.Printf
	// fmt.Printf("\n")
	return
}
