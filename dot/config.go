package dot

import (
	// "fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config ...
type Config struct {
	// Name string
	Roles []*Role
}

// Load ...
func (c *Config) Load(name string) (string, error) {
	// configFile, err := dot.Load(c, name)
	// return configFile, err
	if name == "" {
		return "", nil
	}
	cfgPath := name
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return "", nil
	}
	cfg, err := ioutil.ReadFile(cfgPath)
	if err != nil {
	    return cfgPath, err
	}
	err = yaml.Unmarshal(cfg, &c)
	return cfgPath, err
}

// Execute ...
func (c *Config) Execute() error {
	// fmt.Println("Executing...")
	for _, r := range c.Roles {
		if err := r.Execute(); err != nil {
			return err
		}
	}
	return nil
}

// FindConfig ...
func FindConfig(s string) (string, error) {
	if s == "" {
	    return "", nil
	}

	paths := []string {
		s, // Current working directory
		filepath.Join(os.Getenv("HOME"), s),
		filepath.Join(getConfigDir(), s),
	}

	for _, p := range paths {
		if isFile(p) {
			return p, nil
		}
	}

	return "", nil
}

// shibukawa/configdir
func getConfigDir() string {
	dir := ""
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		dir = os.Getenv("XDG_CONFIG_HOME")
	} else {
		dir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	// XDG_CONFIG_DIRS /etc/xdg
	return dir
}

func readConfig(s string) ([]byte, error) {
	bytes, err := read(s)
	// str := string(bytes)
	// if err != nil {
	// 	return str, err
	// }
	return bytes, err
}

// func exists(s string) bool {
// 	_, err := os.Stat(s)
// 	return !os.IsNotExist(err)
// }

func isFile(s string) bool {
	fi, err := os.Stat(s)

	return !os.IsNotExist(err) && !fi.IsDir()
}

func read(s string) ([]byte, error) {
	return ioutil.ReadFile(s)
}
