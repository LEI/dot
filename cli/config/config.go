package config

import (
	// "fmt"
	"io"
	"os"
	"path/filepath"

	// "github.com/LEI/dot/cli/config/configfile"
	"github.com/spf13/viper"
)

// Config ...
type Config struct {
	// Filename string
	viper *viper.Viper
}

const (
	// ConfigFileName is the type of config file
	// ConfigFileType = "yaml"
	// ConfigFileName is the name of config file
	ConfigFileName = "dotrc"
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

// NewConfig initializes an empty configuration file
func NewConfig() *Config {
	v := viper.New()
	// v.SetConfigType(ConfigFileType)
	v.SetConfigName(ConfigFileName)
	v.AddConfigPath("/etc/"+configFileDir)
	// v.AddConfigPath("$HOME/.dot")
	// v.AddConfigPath(".")
	return &Config{
		// Filename: fn,
		viper: v,
	}
}

// LoadFromReader is a convenience function that creates a Config object from
// a reader
func LoadFromReader(configData io.Reader) (*Config, error) {
	config := Config{
	    viper: viper.New(),
	}
	// err := config.LoadFromReader(configData)
	err := config.viper.ReadConfig(configData)
	return &config, err
}

// Load reads the configuration files in the given directory, and sets up
// the auth config information and returns values.
// FIXME: use the internal golang config parser
func Load(configDir string) (*Config, error) {
	if configDir == "" {
		configDir = Dir()
	}
	config := Config{
	    viper: viper.New(),
	}
	config.viper.SetConfigName(ConfigFileName)
	config.viper.AddConfigPath(configDir)
	err := config.viper.ReadInConfig()
	if err != nil {
	    return &config, nil
	}
	// AuthConfigs: make(map[string]types.AuthConfig),
	// Filename:    filepath.Join(configDir, ConfigFileName),
	// if _, err := os.Stat(config.Filename); err == nil {
	// 	file, err := os.Open(config.Filename)
	// 	if err != nil {
	// 		return &config, fmt.Errorf("%s - %v", config.Filename, err)
	// 	}
	// 	defer file.Close()
	// 	err = config.LoadFromReader(file)
	// 	if err != nil {
	// 		err = fmt.Errorf("%s - %v", config.Filename, err)
	// 	}
	// 	return &config, err
	// } else if !os.IsNotExist(err) {
	// 	// if file is there but we can't stat it for any reason other
	// 	// than it doesn't exist then stop
	// 	return &config, fmt.Errorf("%s - %v", config.Filename, err)
	// }
	return &config, nil
}
