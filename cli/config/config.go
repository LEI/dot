package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/LEI/dot/cli/config/configfile"
)

const (
	// ConfigFileName is the name of config file
	ConfigFileName = "dotrc.yaml"
	configFileDir  = ".dot"
)

var (
	configDir = os.Getenv("DOCKER_CONFIG")
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

// LoadFromReader is a convenience function that creates a ConfigFile object from
// a reader
func LoadFromReader(configData io.Reader) (*configfile.ConfigFile, error) {
	configFile := configfile.ConfigFile{
		// AuthConfigs: make(map[string]types.AuthConfig),
	}
	err := configFile.LoadFromReader(configData)
	return &configFile, err
}

// Load reads the configuration files in the given directory, and sets up
// the auth config information and returns values.
// FIXME: use the internal golang config parser
func Load(configDir string) (*configfile.ConfigFile, error) {
	if configDir == "" {
		configDir = Dir()
	}
	configFile := configfile.ConfigFile{
		// AuthConfigs: make(map[string]types.AuthConfig),
		Filename:    filepath.Join(configDir, ConfigFileName),
	}
	if _, err := os.Stat(configFile.Filename); err == nil {
		file, err := os.Open(configFile.Filename)
		if err != nil {
			return &configFile, fmt.Errorf("%s - %v", configFile.Filename, err)
		}
		defer file.Close()
		err = configFile.LoadFromReader(file)
		if err != nil {
			err = fmt.Errorf("%s - %v", configFile.Filename, err)
		}
		return &configFile, err
	} else if !os.IsNotExist(err) {
		// if file is there but we can't stat it for any reason other
		// than it doesn't exist then stop
		return &configFile, fmt.Errorf("%s - %v", configFile.Filename, err)
	}
	return &configFile, nil
}
