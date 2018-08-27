package shell

import (
	"os"
	"runtime"

	"github.com/opencontainers/runc/libcontainer/user"
)

// github.com/docker/docker/pkg/homedir

// HomeDirKey returns the env var name for the user's home directory.
func HomeDirKey() string {
	return "HOME" // homedir.Key()
}

// GetHomeDir returns the current user home directory.
func GetHomeDir() string { // homedir.Get()
	home := os.Getenv(HomeDirKey())
	if home == "" && runtime.GOOS != "windows" {
		if u, err := user.CurrentUser(); err == nil {
			return u.Home
		}
	}
	return home
}

// GetHomeShortcutString for the home directory.
func GetHomeShortcutString() string {
	return "~" // homedir.GetShortcutString()
}
