package config

// https://github.com/spf13/hugo/blob/master/config/configProvider.go

// Provider provides the configuration settings
type Provider interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	Get(key string) interface{}
	Set(key string, value interface{})
	IsSet(key string) bool
}
