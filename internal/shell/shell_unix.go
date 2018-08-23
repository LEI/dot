// +build !windows

package shell

// $SHELL -c ''

import (
	"fmt"
	"os"
)

// Key returns the env var name for the user's shell.
func Key() string {
	return "SHELL"
}

// Get returns the shell to use.
func Get() string {
	shell := os.Getenv(Key())
	if shell == "" {
		fmt.Fprintf(os.Stderr, "Fallback to default shell: %s", defaultShell)
		shell = defaultShell
	}
	return shell
}

// GetShortcutString returns the variable to use in the native shell.
func GetShortcutString() string {
	return "$HOME"
}
