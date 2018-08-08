package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config structure
type Config struct {
	Source string
	Target string
	Roles []*Role
	// Filename string
	// value interface{}
	v *viper.Viper
}

// Get a value
func (c *Config) Get(key string) interface{} {
	return c.v.Get(key)
}

// GetAll values
func (c *Config) GetAll() map[string]interface{} {
	return c.v.AllSettings()
}

// Parse into struct
func (c *Config) Parse(i interface{}) error {
	// c.value = &i
	return c.v.Unmarshal(&i)
}

// LoadRole config
func (c *Config) LoadRole(r *Role) error {
	ConfigFileName = ".dot" // -rc
	roleConfig, err := Load(r.Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Error loading config file: %v\n", err)
	}
	if roleConfig == nil {
		fmt.Fprintf(os.Stderr, "WARNING: nil role\n")
		return nil
	}
	// configFile := roleConfig.FileUsed()
	// if configFile != "" { // debug
	// 	fmt.Printf("Using role config file: %s\n", configFile)
	// }
	role := roleConfig.Get("role").(map[string]interface{})
	if err := r.Merge(role); err != nil {
		return err
	}
	return nil
}

// // Value config
// func (c *Config) Value() interface{} {
// 	return c.value
// }

// FileUsed by viper
func (c *Config) FileUsed() string {
	return c.v.ConfigFileUsed()
}

func (c *Config) setName(name string) {
	c.v.SetConfigName(name)
}

func (c *Config) setType(name string) {
	c.v.SetConfigType(name)
}

func (c *Config) addPaths(paths ...string) {
	addConfigPaths(c.v, paths)
}
