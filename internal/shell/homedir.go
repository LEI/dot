package shell

import (
	"os"

	"github.com/mitchellh/go-homedir"
)

// HomeDirKey returns the env var name for the user's home directory.
// docker/pkg/homedir.Key()
func HomeDirKey() string {
	return "HOME"
}

// GetHomeDir returns the current user home directory.
// docker/pkg/homedir.Get()
func GetHomeDir() string {
	home := os.Getenv(HomeDirKey())
	/* if home == "" && runtime.GOOS != "windows" {
		if u, err := user.CurrentUser(); err == nil {
			return u.Home
		}
	} */
	if home == "" {
		/* if usr, err := user.Current(); err != nil {
			return usr.HomeDir
		} */
		if dir, err := homedir.Dir(); err != nil {
			return dir
		}
	}
	return home
}

// GetHomeShortcutString for the home directory.
// docker/pkg/homedir.GetShortcutString()
func GetHomeShortcutString() string {
	return "~" // homedir.Expand()
}
