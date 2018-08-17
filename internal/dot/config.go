package dot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mitchellh/mapstructure"
	toml "github.com/pelletier/go-toml"
	"gopkg.in/yaml.v2"
)

// Config struct
type Config struct {
	// Source string
	// Target string
	Roles []*Role
}

// Prepare config roles
func (c *Config) Prepare() error {
	roles := []*Role{}
	for _, r := range c.Roles {
		if err := r.Prepare(); err != nil {
			return err
		}
		// if ok := r.Ignore(); ok {
		// 	continue
		// }
		if err := r.LoadConfig(); err != nil {
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
	return &cfg, err
}

// FindConfig ...
func FindConfig(path string) string {
	return path
}

// LoadConfig ...
func LoadConfig(path string) (Config, error) {
	cfg := Config{}
	data, err := Read(path)
	if err != nil {
		return cfg, err
	}
	var md mapstructure.Metadata
	// decoderConfig := &mapstructure.DecoderConfig{
	// 	WeaklyTypedInput: true,
	// 	Result:           &cfg,
	// }
	// decoder, err := mapstructure.NewDecoder(decoderConfig)
	// if err != nil {
	// 	return err
	// }
	// return decoder.Decode(data)
	if err := mapstructure.WeakDecodeMetadata(data, &cfg, &md); err != nil {
		return cfg, err
	}
	// fmt.Printf("md: %+v\n", md)
	return cfg, nil
}

// Read ...
func Read(path string) (map[string]interface{}, error) {
	// fmt.Println("Loading config file", path)
	var data map[string]interface{}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return data, err
	}
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
