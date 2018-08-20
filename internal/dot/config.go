package dot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	toml "github.com/pelletier/go-toml"
	"gopkg.in/yaml.v2"
)

var (
	// DecodeErrorUnused mapstructure decode option
	DecodeErrorUnused = true // TODO false if opts.Force

	// DecodeWeaklyTypedInput mapstructure decode option
	DecodeWeaklyTypedInput = true
)

// Config struct
type Config struct {
	Source string
	Target string
	Roles  []*Role

	dirname  string // Role directory name
	filename string // Role config file name
}

// SetDir name
func (c *Config) SetDir(name string) {
	c.dirname = name
}

// SetRoleFile name
func (c *Config) SetRoleFile(name string) {
	c.filename = name
}

// ParseRoles config
func (c *Config) ParseRoles() error {
	roles := []*Role{}
	for _, r := range c.Roles {
		if r.Path == "" {
			r.Path = filepath.Join(c.dirname, r.Name)
			if !filepath.IsAbs(r.Path) {
				r.Path = filepath.Join(c.Source, r.Path)
			}
		}
		// if ok := r.Ignore(); ok {
		// 	continue
		// }
		r.SetConfigFile(c.filename)
		if exists(r.GetConfigFile()) {
			if err := r.LoadConfig(); err != nil {
				return err
			}
		}
		if err := r.Parse(c.Target); err != nil {
			return err
		}
		roles = append(roles, r)
	}
	c.Roles = roles
	return nil
}

// NewConfig ...
func NewConfig(path string) (*Config, error) {
	// if path == "" {}
	cfgPath := FindConfig(path)
	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		return &cfg, fmt.Errorf("error loading config: %s", err)
	}
	return &cfg, nil
}

// FindConfig ...
func FindConfig(path string) string {
	return path
}

// LoadConfig ...
func LoadConfig(path string) (Config, error) {
	cfg := Config{}
	data, err := ReadConfigFile(path)
	if err != nil {
		return cfg, err
	}
	// var md mapstructure.Metadata
	// if err := mapstructure.WeakDecodeMetadata(data, &cfg, &md); err != nil {
	// 	return cfg, err
	// }
	// fmt.Printf("md: %+v\n", md)
	dc := &mapstructure.DecoderConfig{
		// DecodeHook:       ...,
		ErrorUnused:      DecodeErrorUnused,
		WeaklyTypedInput: DecodeWeaklyTypedInput,
		Result:           &cfg,
	}
	decoder, err := mapstructure.NewDecoder(dc)
	if err != nil {
		return cfg, err
	}
	err = decoder.Decode(data)
	return cfg, err
}

// ReadConfigFile ...
func ReadConfigFile(path string) (map[string]interface{}, error) {
	var data map[string]interface{}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return data, err
	}
	// TODO if Verbose fmt.Println("## Loaded config file", path)
	cfgType := detectType(path)
	switch cfgType {
	case "toml":
		if err := toml.Unmarshal(b, &data); err != nil {
			return data, err
		}
	case "yaml":
		if err := yaml.Unmarshal(b, &data); err != nil {
			return data, err
		}
	case "json":
		if err := json.Unmarshal(b, &data); err != nil {
			return data, err
		}
	default:
		return data, fmt.Errorf("%s: unknown config type", path)
	}
	return data, nil
}

func detectType(path string) string {
	var fileType string
	switch {
	case strings.HasSuffix(path, ".toml"):
		fileType = "toml"
	case strings.HasSuffix(path, ".yaml"), strings.HasSuffix(path, ".yml"):
		fileType = "yaml"
	case strings.HasSuffix(path, ".json"):
		fileType = "json"
	}
	return fileType
}

// func weaklyTypedHook(
// 	f reflect.Kind,
// 	t reflect.Kind,
// 	data interface{}) (interface{}, error) {
// 	dataVal := reflect.ValueOf(data)
// 	switch t {
// 	case reflect.String:
// 		switch f {
// 		case reflect.Bool:
// 			if dataVal.Bool() {
// 				return "1", nil
// 			}
// 			return "0", nil
// 		case reflect.Float32:
// 			return strconv.FormatFloat(dataVal.Float(), 'f', -1, 64), nil
// 		case reflect.Int:
// 			return strconv.FormatInt(dataVal.Int(), 10), nil
// 		case reflect.Slice:
// 			dataType := dataVal.Type()
// 			elemKind := dataType.Elem().Kind()
// 			if elemKind == reflect.Uint8 {
// 				return string(dataVal.Interface().([]uint8)), nil
// 			}
// 		case reflect.Uint:
// 			return strconv.FormatUint(dataVal.Uint(), 10), nil
// 		}
// 	}
// 	return data, nil
// }
