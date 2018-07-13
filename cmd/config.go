package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config ...
type Config struct {
	// Name string
	Roles      []*Role
	IgnoreDeps bool
}

var config *Config

// ErrSkipDeps ...
// var ErrSkipDeps = fmt.Errorf("skip dependencies")

// NewConfig ...
func NewConfig() *Config {
	if config == nil {
		config = &Config{}
	}
	return config
}

// Read ...
func (c *Config) Read(name string) (string, error) {
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

// Prepare ...
func (c *Config) Prepare() error {
	// for _, n := range r.Deps {
	// 	fmt.Println("DEP", n)
	// }
	return nil
}

// AddRole ...
func (c *Config) AddRole(r *Role) error {
	for i, role := range c.Roles {
		if role.Name == r.Name {
			if err := r.Merge(role); err != nil {
				return err
			}
			// c.SetRoleIndex(i, r)
			c.Roles[i] = r
			return nil
			// break
		}
	}
	c.Roles = append(c.Roles, r)
	return nil
}

// RemoveRole ...
func (c *Config) RemoveRole(r *Role) (ret []*Role) {
	for _, role := range c.Roles {
		if role == r {
			continue
		}
		ret = append(ret, r)
	}
	return ret
}

// Require ...
func (c *Config) Require() error {
	if c.IgnoreDeps {
		return nil
	}
CHECK:
	for _, role := range c.Roles {
		if !role.IsEnabled() {
			continue
		}
		if len(role.Dependencies) < 0 {
			continue
		}
		for _, dep := range role.Dependencies {
			for _, r := range c.Roles {
				if r.Name == dep {
					fmt.Println(role.Name, "requires", r.Name)
					if !r.IsEnabled() {
						r.Enable()
					}
					continue CHECK
				}
			}
			return fmt.Errorf("Unable to resolve %s dependency: %s", role.Name, dep)
		}
	}
	// TODO: sort
	return nil
}

// FindConfig ...
func FindConfig(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	paths := []string{
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
