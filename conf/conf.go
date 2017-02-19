package conf

import (
	"github.com/spf13/viper"
)

type Conf struct {
	*viper.Viper
}

func New(name string, paths []string) *Conf {
	v := &Conf{viper.New()}
	v.SetConfigName(name)
	for _, p := range paths {
		v.AddConfigPath(p)
	}
	v.AutomaticEnv()
	return v
}

func NewFile(file string) *Conf {
	v := &Conf{viper.New()}
	v.SetConfigFile(file)
	return v
}

func (v *Conf) SetFile(file string) {
	v.SetConfigFile(file)
}

func (v *Conf) Read() (string, error) {
	err := v.ReadInConfig()
	used := v.ConfigFileUsed()
	return used, err
}

// func (v *Conf) Key(key string, v interface{}) error {
// 	err := v.UnmarshalKey(key, &v)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
