package config

import (
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	configFileDir  = "" // HOME ".dot"
)

var (
	// ConfigFileType is the type of config file
	ConfigFileType = "yaml"
	// ConfigFileName is the name of config file
	ConfigFileName = ".dotrc"
	homeDir = os.Getenv("HOME")
	configDir = os.Getenv("DOT_CONFIG")
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
	config.addPaths(dir) // configDir
	err := config.v.ReadInConfig()
	if err != nil {
		return &config, err
	}
	return &config, nil
}

func addConfigPaths(v *viper.Viper, paths []string) {
	for _, p := range paths {
		v.AddConfigPath(p)
	}
}
