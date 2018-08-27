package shell

import (
	"github.com/docker/docker/pkg/homedir"
)

// HomeDirKey returns the env var name for the user's home directory.
func HomeDirKey() string {
	return homedir.Key() // "HOME"
}

// GetHomeDir returns the current user home directory.
func GetHomeDir() string {
	return homedir.Get() // os.Getenv(HomeDirKey())
}

// GetHomeDirShortcutString for the home directory.
func GetHomeDirShortcutString() string {
	return homedir.GetShortcutString() // "~"
}
