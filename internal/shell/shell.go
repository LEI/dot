package shell

import (
	"fmt"
	"os"
	"strings"
)

var (
	defaultShell = "sh" // command -v sh

	// HomeDir of the current user.
	HomeDir = GetHomeDir()
)

// Key returns the env var name for the user's shell.
func Key() string {
	return "SHELL"
}

// Get returns the shell to use.
func Get() string {
	shell := os.Getenv(Key())
	if shell == "" {
		// fmt.Fprintf(Stderr, "Fallback to default shell: %s\n", defaultShell)
		shell = defaultShell
	}
	return shell
}

// GetShortcutString returns the variable to use in the native shell.
func GetShortcutString() string {
	return "$SHELL"
}

// FormatArgs formats command line aruments to a single string.
func FormatArgs(args []string) string {
	for i, a := range args {
		if strings.Contains(a, " ") {
			args[i] = fmt.Sprintf("%q", a)
			// windows? args[i] = syscall.EscapeArg(a)
		}
		// switch v := a.(type) {
		// case string:
		// }
	}
	return strings.Join(args, " ")
}
