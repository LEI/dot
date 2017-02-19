package conf

import (
	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
}

func New(name string, paths []string) *Config {
	c := &Config{viper.New()}
	c.SetConfigName(name)
	for _, p := range paths {
		c.AddConfigPath(p)
	}
	c.AutomaticEnv()
	return c
}

func (c *Config) SetFile(file string) {
	c.SetConfigFile(file)
}

func (c *Config) Read() (string, error) {
	err := c.ReadInConfig()
	used := c.ConfigFileUsed()
	return used, err
}

// func (c *Config) Key(key string, v interface{}) error {
// 	err := c.UnmarshalKey(key, &v)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
