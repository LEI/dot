package config

// https://github.com/spf13/hugo/blob/master/config/configProvider.go

// Provider provides the configuration settings
type Provider interface {
	Get(key string) interface{}
	GetBool(key string) bool
	GetInt(key string) int
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringSlice(key string) []string
	IsSet(key string) bool
	Set(key string, value interface{})
	UnmarshalKey(key string, value interface{})
}
