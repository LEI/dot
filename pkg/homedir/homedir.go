package homedir
// https://github.com/moby/moby/blob/master/pkg/homedir/homedir_unix.go

import (
	"log"
	"os"
)

// Key returns the env var name for the user's home dir based on
// the platform being run on
func Key() string {
	return "HOME"
}

// Get returns the home directory of the current user with the help of
// environment variables.
// Returned path should be used with "path/filepath" to form new paths.
func Get() string {
	home := os.Getenv(Key())
	if home == "" {
		log.Fatalf("invalid homedir '%s", os.ExpandEnv(Key()))
	}
	return home
}

// GetShortcutString returns the string that is shortcut to user's home directory
// in the native shell of the platform running on.
func GetShortcutString() string {
	return "~"
}
