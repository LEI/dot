package shell

// https://github.com/magefile/mage/blob/master/sh/cmd.go

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Run command without specifying any environment variables.
func Run(cmd string, args ...string) error {
	// return RunWith(nil, cmd, args...)
	_, err := Exec(nil, Stdout, Stderr, cmd, args...)
	return err
}

// RunCmd uses Run and Exec underneath.
func RunCmd(cmd string, args ...string) func(args ...string) error {
	return func(args2 ...string) error {
		return Run(cmd, append(args, args2...)...)
	}
}

// RunWith executes a command within the given environment.
func RunWith(env map[string]string, cmd string, args ...string) error {
	// var output io.Writer
	// if Verbose {
	// 	output = Stdout
	// }
	_, err := Exec(env, Stdout, Stderr, cmd, args...)
	return err
}

// Output runs the command and returns the text from stdout.
func Output(cmd string, args ...string) (string, error) {
	buf := &bytes.Buffer{}
	_, err := Exec(nil, buf, Stderr, cmd, args...)
	return strings.TrimSuffix(buf.String(), "\n"), err
}

// OutputCmd uses Ouput and Exec underneath.
func OutputCmd(cmd string, args ...string) func(args ...string) (string, error) {
	return func(args2 ...string) (string, error) {
		return Output(cmd, append(args, args2...)...)
	}
}

// OutputWith is like RunWith, ubt returns what is written to stdout.
func OutputWith(env map[string]string, cmd string, args ...string) (string, error) {
	buf := &bytes.Buffer{}
	_, err := Exec(env, buf, Stderr, cmd, args...)
	return strings.TrimSuffix(buf.String(), "\n"), err
}

// CombinedOutput runs the command and returns the text from stdout and stderr.
func CombinedOutput(cmd string, args ...string) (string, error) {
	buf := &bytes.Buffer{}
	_, err := Exec(nil, buf, buf, cmd, args...)
	return strings.TrimSuffix(buf.String(), "\n"), err
}

// CombinedOutputCmd uses CombinedOutput and Exec underneath.
func CombinedOutputCmd(cmd string, args ...string) func(args ...string) (string, error) {
	return func(args2 ...string) (string, error) {
		return CombinedOutput(cmd, append(args, args2...)...)
	}
}

// CombinedOutputWith is like RunWith, ubt returns what is written to stdout. and stderr
func CombinedOutputWith(env map[string]string, cmd string, args ...string) (string, error) {
	buf := &bytes.Buffer{}
	_, err := Exec(env, buf, buf, cmd, args...)
	return strings.TrimSuffix(buf.String(), "\n"), err
}

// Exec executes the command.
func Exec(env map[string]string, stdout, stderr io.Writer, cmd string, args ...string) (ran bool, err error) {
	expand := func(s string) string {
		s2, ok := env[s]
		if ok {
			return s2
		}
		return os.Getenv(s)
	}
	cmd = os.Expand(cmd, expand)
	for i := range args {
		args[i] = os.Expand(args[i], expand)
	}
	ran, code, err := run(env, stdout, stderr, cmd, args...)
	if err == nil {
		return true, nil
	}
	if ran {
		return ran, fmt.Errorf(`exit code %d running "%s %s" failed with exit code %d`, code, cmd, strings.Join(args, " "), code)
	}
	return ran, fmt.Errorf(`failed to run "%s %s: %v"`, cmd, strings.Join(args, " "), err)
}

func run(env map[string]string, stdout, stderr io.Writer, cmd string, args ...string) (ran bool, code int, err error) {
	c := exec.Command(cmd, args...)
	// fmt.Println("run:", c)
	c.Env = os.Environ()
	for k, v := range env {
		c.Env = append(c.Env, k+"="+v)
	}
	c.Stderr = stderr
	c.Stdout = stdout
	c.Stdin = Stdin
	fmt.Println("exec:", cmd, strings.Join(args, " "))
	err = c.Run() // FIXME: stdout, stderr *bytes.Buffer
	return cmdRan(err), ExitStatus(err), err
}

func cmdRan(err error) bool {
	if err == nil {
		return true
	}
	ee, ok := err.(*exec.ExitError)
	if ok {
		return ee.Exited()
	}
	return false
}

type exitStatus interface {
	ExitStatus() int
}

// ExitStatus returns the exit status of the error if it is an exec.ExitError
// or if it implements ExitStatus() int.
// 0 if it is nil or 1 if it is a different error.
func ExitStatus(err error) int {
	if err == nil {
		return 0
	}
	if e, ok := err.(exitStatus); ok {
		return e.ExitStatus()
	}
	if e, ok := err.(*exec.ExitError); ok {
		if ex, ok := e.Sys().(exitStatus); ok {
			return ex.ExitStatus()
		}
	}
	return 1
}

// // Run a command
// func Run(cmd *exec.Cmd) (status int) {
// 	if err := cmd.Run(); err != nil {
// 		if exitError, ok := err.(*exec.ExitError); ok {
// 			if ws, ok := exitError.Sys().(syscall.WaitStatus); ok {
// 				status = ws.ExitStatus()
// 			} else {
// 				// fmt.Fprintf(os.Stderr, "Could not get exit status: %+v\n", cmd.Args)
// 				status = defaultStatusFailed
// 			}
// 		} else {
// 			// fmt.Fprintf(os.Stderr, "Could not get exit error: %+v\n", cmd.Args)
// 			status = defaultStatusFailed
// 			// if stdErr == "" { stdErr = err.Error() }
// 		}
// 	} else {
// 		if ws, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
// 			status = ws.ExitStatus()
// 		} else {
// 			// fmt.Fprintf(os.Stderr, "Could not get processed status: %+v\n", cmd.Args)
// 			status = defaultStatusFailed
// 		}
// 	}
// 	return
// }
