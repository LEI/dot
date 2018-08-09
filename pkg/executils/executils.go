package executils

import (
	"fmt"
	"bytes"
	"io"
	"os"
	"os/exec"
	// "strings"
	"syscall"
)

var (
	// Stdout ...
	Stdout = os.Stdout
	// Stderr ...
	Stderr = os.Stderr

	defaultStatusFailed = 1
)

// Run executes the command
func Run(cmd *exec.Cmd) (status int) {
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if ws, ok := exitError.Sys().(syscall.WaitStatus); ok {
				status = ws.ExitStatus()
			} else {
				// fmt.Fprintf(os.Stderr, "Could not get exit status: %+v\n", cmd.Args)
				status = defaultStatusFailed
			}
		} else {
			// fmt.Fprintf(os.Stderr, "Could not get exit error: %+v\n", cmd.Args)
			status = defaultStatusFailed
			// if stdErr == "" { stdErr = err.Error() }
		}
	} else {
		if ws, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
			status = ws.ExitStatus()
		} else {
			// fmt.Fprintf(os.Stderr, "Could not get processed status: %+v\n", cmd.Args)
			status = defaultStatusFailed
		}
	}
	return
}

// RunStd command
func RunStd(cmd *exec.Cmd) (status int) {
	return RunStream(Stdout, Stderr, cmd)
}

// RunBuf command
func RunBuf(cmd *exec.Cmd) (stdOut, stdErr []byte, status int) {
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	status = Run(cmd)
	stdOut = outbuf.Bytes()
	stdErr = errbuf.Bytes()
	return
}

// RunStream command
func RunStream(stdOut, stdErr io.Writer, cmd *exec.Cmd) (status int) {
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
	return Run(cmd)
}

// Execute command
func Execute(name string, args ...string) (status int) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = Stdout
	cmd.Stderr = Stderr
	return Run(cmd)
	// return RunStd(cmd)
}

// ExecuteBuf command
func ExecuteBuf(name string, args ...string) (stdOut, stdErr []byte, status int) {
	cmd := exec.Command(name, args...)
	return RunBuf(cmd)
}

// ExecuteStream command
func ExecuteStream(stdOut, stdErr io.Writer, name string, args ...string) (status int) {
	cmd := exec.Command(name, args...)
	return RunStream(stdOut, stdErr, cmd)
}

// ExecuteInDir command
func ExecuteInDir(dir, name string, args ...string) (stdOut, stdErr []byte, status int) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return RunBuf(cmd)
}

// ExecuteEnv command
func ExecuteEnv(env []string, name string, args ...string) (stdOut, stdErr []byte, status int) {
	cmd := exec.Command(name, args...)
	cmd.Env = env
	return RunBuf(cmd)
}
