package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/LEI/dot/utils"
)

// Config ...
type Config struct {
	// Name string
	Roles      []*Role
	IgnoreDeps bool
}

// ErrSkipDeps ...
// var ErrSkipDeps = fmt.Errorf("skip dependencies")

// Read ...
func (c *Config) Read(s string) error {
	if s == "" {
		return nil
	}
	if !utils.Exist(s) {
		return nil
	}
	cfg, err := ioutil.ReadFile(s)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(cfg, &c)
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
		}
	}
	c.Roles = append(c.Roles, r)
	return nil
}

// RemoveRole ...
func (c *Config) RemoveRole(r *Role) error {
	ret := []*Role{} // New slice
	for _, role := range c.Roles {
		// reflect.DeepEqual(role, r)
		if role.Name == r.Name {
			continue
		}
		ret = append(ret, role)
	}
	if len(ret) == len(c.Roles) {
		return fmt.Errorf("# Unable to remove role '%s' (not found)", r.Name)
	}
	c.Roles = ret
	// (*c).Roles = ret
	return nil
}

// Require ...
func (c *Config) Require() error {
	if c.IgnoreDeps {
		return nil
	}
	// CHECK:
	for _, role := range c.Roles {
		if err := c.RequireRole(role); err != nil {
			return err
		}
	}
	// TODO: sort
	return nil
}

// RequireRole ...
func (c *Config) RequireRole(role *Role) error {
	if c.IgnoreDeps {
		return nil
	}
	if !role.IsEnabled() {
		return nil
	}
	if len(role.Deps) < 0 {
		return nil
	}
	for _, dep := range role.Deps {
		for _, r := range c.Roles {
			if r.Name == dep {
				fmt.Println(role.Name, "requires", r.Name)
				if !r.IsEnabled() {
					r.Enable()
				}
				return nil // continue CHECK
			}
		}
		return fmt.Errorf("unable to resolve %s dependency: %s", role.Name, dep)
	}
	return nil
}

// FindConfig ...
func FindConfig(name string) (string, error) {
	if name == "" {
		return "", nil
	}
	cwd, err := filepath.Abs(name)
	if err != nil {
		return name, err
	}
	paths := []string{
		filepath.Join(os.Getenv("HOME"), name),
		filepath.Join(getConfigDir(), name),
		cwd, // Current working directory
		// Search in CWD last to prevent an eventual
		// local role config to be found before
		// TODO: relative to target directory
	}
	// fmt.Println("Config: searching in", paths)
	for _, p := range paths {
		//fmt.Println("# Config path", p, "->", utils.IsFile(p))
		if utils.IsFile(p) {
			// fmt.Println("Config: found config", p)
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
