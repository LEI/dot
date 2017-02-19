package config

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	*viper.Viper
}

func New(name string, paths []string) *Configuration {
	v := &Configuration{viper.New()}
	v.SetConfigName(name)
	for _, p := range paths {
		v.AddConfigPath(p)
	}
	v.AutomaticEnv()
	return v
}

func NewFile(file string) *Configuration {
	v := &Configuration{viper.New()}
	v.SetConfigFile(file)
	return v
}

func (v *Configuration) SetFile(file string) {
	v.SetConfigFile(file)
}

func (v *Configuration) Read() (string, error) {
	err := v.ReadInConfig()
	used := v.ConfigFileUsed()
	return used, err
}

// func (v *Configuration) Key(key string, v interface{}) error {
// 	err := v.UnmarshalKey(key, &v)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
