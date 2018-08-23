// +build windows

package shell

import (
	"fmt"
	"os"
)

// C:\cygwin64\bin\bash -c ''

// Key returns the env var name for the user's shell.
func Key() string {
	return "" // CYG_BASH
}

// Get returns the shell to use.
func Get() string {
	// return os.Getenv(Key())
	fmt.Fprintf(os.Stderr, "Default shell: %s", defaultShell)
	return defaultShell // "/cygdrive/c/cygwin64/bin/bash"
}

// GetShortcutString returns the variable to use in the native shell.
func GetShortcutString() string {
	return "" // %CYG_BASH%
}
