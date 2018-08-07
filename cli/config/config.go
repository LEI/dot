package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	// "github.com/LEI/dot/cli/config/configfile"
	"github.com/spf13/viper"
)

// Config structure
type Config struct {
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

const (
	// ConfigFileType is the type of config file
	ConfigFileType = "yaml"
	// ConfigFileName is the name of config file
	ConfigFileName = ".dotrc"
	configFileDir  = "" // HOME ".dot"
)

var (
	configDir = os.Getenv("DOT_CONFIG")
	homeDir = os.Getenv("HOME")
)

func init() {
	// https://github.com/moby/moby/blob/17.05.x/pkg/homedir/homedir.go
	if configDir == "" {
		configDir = filepath.Join(homeDir, configFileDir)
	}
}

// Dir returns the directory the configuration file is stored in
func Dir() string {
	return configDir
}

// SetDir sets the directory the configuration file is stored in
func SetDir(dir string) {
	configDir = dir
}

// NewConfig initializes an empty configuration file
func NewConfig(name string) *Config {
	c := &Config{
		// Filename: fn,
		v: viper.New(),
	}
	c.setName(name)
	c.setType(ConfigFileType)
	c.addPaths(configDir) // "/etc/"+configFileDir
	// c.v.AddConfigPath("$HOME/.dot")
	// c.v.AddConfigPath(".")
	return c
}

// LoadFromReader is a convenience function that creates a Config object from
// a reader
func LoadFromReader(configData io.Reader) (*Config, error) {
	config := Config{
		v: viper.New(),
	}
	err := config.v.ReadConfig(configData)
	return &config, err
}

// Load reads the configuration files in the given directory
func Load(dir string) (*Config, error) {
	if dir == "" {
		dir = Dir()
	}
	config := Config{
		v: viper.New(),
	}
	config.setName(ConfigFileName)
	config.setType(ConfigFileType)
	config.addPaths(configDir)
	err := config.v.ReadInConfig()
	if err != nil {
		return &config, err
	}
	configFile := config.FileUsed()
	if configFile != "" { // debug
		fmt.Printf("Using config file: %s\n", configFile)
	}
	return &config, nil
}

func addConfigPaths(v *viper.Viper, paths []string) {
	for _, p := range paths {
		v.AddConfigPath(p)
	}
}
