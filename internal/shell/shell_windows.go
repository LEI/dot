// +build windows

package shell

// %CYG_BASH% -lc ''

import "os"

// Key returns the env var name for the user's shell.
func Key() string {
	return "CYG_BASH"
}

// Get returns the shell to use.
func Get() string {
	return os.Getenv(Key())
}

// GetShortcutString returns the variable to use in the native shell.
func GetShortcutString() string {
	return "%CYG_BASH%"
}
