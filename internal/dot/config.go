package dot

import (
	"bytes"
	"fmt"
	"net/url"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/LEI/dot/internal/conf"
	"github.com/LEI/dot/internal/git"
	"github.com/mitchellh/mapstructure"
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
	Git    *url.URL

	dirname  string // Role directory name
	filename string // Role config file name
	file     string // actual file used
}

// NewConfig ...
func NewConfig(path, dirname string) (*Config, error) {
	// if path == "" { ... }
	cfg := &Config{
		dirname: dirname,
		file:    path,
	}
	if err := cfg.Load(); err != nil {
		return cfg, fmt.Errorf("error loading config: %s", err)
	}
	return cfg, nil
}

// FileUsed path
func (c *Config) FileUsed() string {
	return c.file
}

// SetRoleFile name
func (c *Config) SetRoleFile(name string) {
	c.filename = name
}

// Load config from file
func (c *Config) Load() error {
	c.file = FindConfig(c.file, c.dirname)
	data, err := conf.ReadFile(c.file)
	if err != nil {
		return err
	}
	// var md mapstructure.Metadata
	// if err := mapstructure.WeakDecodeMetadata(data, &c, &md); err != nil {
	// 	return err
	// }
	// fmt.Printf("md: %+v\n", md)
	dc := &mapstructure.DecoderConfig{
		DecodeHook:       configDecodeHook,
		ErrorUnused:      DecodeErrorUnused,
		WeaklyTypedInput: DecodeWeaklyTypedInput,
		Result:           &c,
	}
	decoder, err := mapstructure.NewDecoder(dc)
	if err != nil {
		return err
	}
	err = decoder.Decode(data)
	if err != nil {
		return err
	}
	return nil
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

		// if r.URL == "" { r.URL = r.Name }
		// r.URL = git.ParseURL(r.Git.User, r.Git.Host, r.URL)

		r.SetConfigFile(c.filename)
		if f := r.GetConfigFile(); exists(f) {
			if err := git.CheckRemote(r.Path, r.URL); err != nil {
				return err
			}
			if err := r.Load(); err != nil {
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

// FindConfig searches a given file name or path
// relative to the home directory, or falls back
// to ~/.dot/config
func FindConfig(path, dirname string) string {
	// dirs := []string{".", homeDir}
	// /etc/dot, $HOME/.dot/config, $HOME/.config/dot...
	if filepath.IsAbs(path) {
		return path
	}
	// Current working directory
	if exists(path) {
		return path
	}
	// Home directory
	if rc := filepath.Join(homeDir, path); exists(rc) {
		return rc
	}
	// path = strings.TrimPrefix(path, ".")
	return filepath.Join(homeDir, dirname, "config")
}

func configDecodeHook(f reflect.Type, t reflect.Type, i interface{}) (interface{}, error) {
	switch val := i.(type) {
	case string:
		switch t {
		case reflect.TypeOf((*Role)(nil)):
			i = &Role{Name: val}
		}
	case *url.URL:
		// fmt.Println("DECODE URL", val)
	case *Role:
		if val.Name != "" && val.URL == "" {
			sep := ":" // os.PathListSeparator
			notURL := bytes.Count([]byte(val.Name), []byte(sep)) == 1
			if notURL {
				parts := strings.SplitN(val.Name, sep, 2)
				// if !strings.Contains(val.Name, string(os.PathSeparator)) {}
				val.Name = parts[0]
				val.URL = parts[1]
			}
		}
	}
	return i, nil
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
