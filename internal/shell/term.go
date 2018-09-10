package shell

import (
	"io"
	"os"
)

var (
	// Stdout writer
	Stdout io.Writer = os.Stdout
	// Stderr writer
	Stderr io.Writer = os.Stderr
	// Stdin reader
	Stdin io.Reader = os.Stdin
)

/*
// import "github.com/mattn/go-isatty"

// ioctlReadTermios...

// IsTerminal return true if the file descriptor is terminal.
func IsTerminal(fd uintptr) bool {
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, ioctlReadTermios, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return err == 0
}
*/

/*
// import "github.com/docker/docker/pkg/term"

// getTermios, setTermios...

// term/term.go
// \+build !windows

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal(fd uintptr) bool {
	var termios Termios
	return tcget(fd, &termios) == 0
}

// term/tc.go
// \+build !windows

func tcget(fd uintptr, p *Termios) syscall.Errno {
	_, _, err := unix.Syscall(unix.SYS_IOCTL, fd, uintptr(getTermios), uintptr(unsafe.Pointer(p)))
	return err
}

func tcset(fd uintptr, p *Termios) syscall.Errno {
	_, _, err := unix.Syscall(unix.SYS_IOCTL, fd, setTermios, uintptr(unsafe.Pointer(p)))
	return err
}
*/
