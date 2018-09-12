package dot

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/LEI/dot/internal/conf"
	"github.com/LEI/dot/internal/git"
	"github.com/LEI/dot/internal/host"
	"github.com/mitchellh/mapstructure"
)

var (
	// DecodeErrorUnused mapstructure decode option
	DecodeErrorUnused = true // TODO false if opts.Force

	// DecodeWeaklyTypedInput mapstructure decode option
	DecodeWeaklyTypedInput = true

	// Stdout writer
	Stdout io.Writer = os.Stdout
	// Stderr writer
	Stderr io.Writer = os.Stderr
	// Stdin reader
	Stdin io.Reader = os.Stdin
)

// Config struct
type Config struct {
	Source    string
	Target    string
	Roles     []*Role
	Platforms map[string][]*Role
	Git       *url.URL

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
		return cfg, err // fmt.Errorf("error loading config: %s", err)
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

// PrepareRoles config
func (c *Config) PrepareRoles() error {
	roles := []*Role{}
	for _, r := range c.Roles {
		// if ok := r.Ignore(); ok {
		// 	continue
		// }
		if r.Path == "" {
			r.Path = filepath.Join(c.dirname, r.Name)
			if !filepath.IsAbs(r.Path) {
				r.Path = filepath.Join(c.Source, r.Path)
			}
		}
		if len(r.OS) > 0 && !host.HasOS(r.OS...) {
			fmt.Fprintf(os.Stderr, "## Skip load %s (OS: %+v)", r.Name, r.OS)
			continue
		}
		roles = append(roles, r)
	}
	c.Roles = roles
	return nil
}

// ParseRoles config
func (c *Config) ParseRoles() error {
	roles := []*Role{}
	// Clone repository, load config and parse it
	for _, r := range c.Roles {
		// if r.URL == "" { r.URL = r.Name }
		// r.URL = git.ParseURL(r.Git.User, r.Git.Host, r.URL)
		// Verify repository state
		if err := git.CheckRemote(r.Path, r.URL); err != nil {
			return err
		}
		r.SetConfigFile(c.filename)
		// Load role config if found
		if f := r.GetConfigFile(); exists(f) {
			if err := r.Load(); err != nil {
				return err
			}
		}
		if len(r.OS) > 0 && !host.HasOS(r.OS...) {
			fmt.Fprintf(os.Stderr, "## Skip parse %s (OS: %+v)", r.Name, r.OS)
			continue
		}
		if err := r.Parse(c.Target); err != nil {
			return err
		}
		roles = append(roles, r)
	}
	if err := checkDeps(roles); err != nil {
		return err
	}
	// Update config
	c.Roles = roles
	return nil
}

// Verify role dependencies
func checkDeps(roles []*Role) error {
	for i, ro := range roles {
	DEPS:
		for _, name := range ro.Deps {
			for j, r := range roles {
				if name == r.Name {
					if j > i {
						return fmt.Errorf("%s: should be loaded before %s", ro.Name, r.Name)
					}
					continue DEPS
				}
			}
			return fmt.Errorf("%s: requires %s", ro.Name, name)
		}
	}
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
	// fmt.Printf("DECODE %T (%s -> %s)\n", i, f, t)
	switch val := i.(type) {
	case string:
		switch t {
		case reflect.TypeOf((*Role)(nil)):
			i = &Role{Name: val}
		case reflect.TypeOf((*url.Userinfo)(nil)):
			i = url.User(val) // &url.URL{}
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
