package executils

import (
	"fmt"
	"bytes"
	"os"
	"os/exec"
	// "strings"
	"syscall"
)

var (
	defaultStatusFailed = 1
	// dryRun = true // bool
)

// Run command
func Run(cmd *exec.Cmd) (stdout, stderr []byte, status int) {
	// fmt.Println("$", cmd.Args)
	// if dryRun {
	// 	fmt.Println("DRY-RUN: ", cmd.Args)
	// 	return []byte{}, []byte{}, 0
	// }
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if ws, ok := exitError.Sys().(syscall.WaitStatus); ok {
				status = ws.ExitStatus()
			} else {
				fmt.Fprintf(os.Stderr, "Could not get exit status: %+v\n", cmd.Args)
				status = defaultStatusFailed
			}
		} else {
			fmt.Fprintf(os.Stderr, "Could not get exit error: %+v\n", cmd.Args)
			status = defaultStatusFailed
			// if stderr == "" {
			// 	stderr = err.Error()
			// }
		}
	} else {
		if ws, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
			status = ws.ExitStatus()
		} else {
			fmt.Fprintf(os.Stderr, "Could not get processed status: %+v\n", cmd.Args)
			status = defaultStatusFailed
		}
	}
	stdout = outbuf.Bytes()
	stderr = errbuf.Bytes()
	return
}

// Execute command
func Execute(name string, args ...string) (stdout, stderr []byte, status int) {
	cmd := exec.Command(name, args...)
	return Run(cmd)
}

// ExecuteIn directory command
func ExecuteIn(dir, name string, args ...string) (stdout, stderr []byte, status int) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return Run(cmd)
}

// ExecuteEnv command
func ExecuteEnv(env []string, name string, args ...string) (stdout, stderr []byte, status int) {
	cmd := exec.Command(name, args...)
	cmd.Env = env
	return Run(cmd)
}
