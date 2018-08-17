package dot

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
)

// RoleConfig struct
type RoleConfig struct {
	Role Role // `mapstructure:",squash"`
}

// Role struct
type Role struct {
	Name string
	Path string
	URL  string
	// Tasks []string

	OS  []string
	Env map[string]string
	// Vars  types.Map
	// IncludeVars types.IncludeMap

	Deps []string `mapstructure:"dependencies"`
	Pkgs Pkgs     `mapstructure:"pkg"`

	Dirs      Dirs      `mapstructure:"dir"`
	Files     Files     `mapstructure:"copy"`
	Links     Links     `mapstructure:"link"`
	Templates Templates `mapstructure:"template"`
	Lines     Lines     `mapstructure:"line"`

	Install     Hooks
	PostInstall Hooks `mapstructure:"post_install"`
	Remove      Hooks
	PostRemove  Hooks `mapstructure:"post_remove"`

	// Ignore []string
	// Target string

	// synced bool
}

func (r *Role) String() string {
	// return fmt.Sprintf("%s %s", r.Name, r.URL)
	s := ""
	s += fmt.Sprintf("%s\n", r.Name)
	s += fmt.Sprintf("  Path: %s\n", r.Path)
	s += fmt.Sprintf("  URL: %s\n", r.URL)

	if r.Deps != nil {
		s += fmt.Sprintf("  Deps: %s\n", r.Deps)
	}
	if r.Env != nil {
		s += fmt.Sprintf("  Env: %s\n", r.Env)
	}
	if r.OS != nil {
		s += fmt.Sprintf("  OS: %s\n", r.OS)
	}
	if r.Pkgs != nil {
		s += fmt.Sprintf("  Pkgs: %s\n", r.Pkgs)
	}
	if t := tasksPrefix("  ", r); t != "" {
		s += t
	}
	if r.Install != nil {
		s += fmt.Sprintf("  Install: %s\n", r.Install)
	}
	if r.PostInstall != nil {
		s += fmt.Sprintf("  PostInstall: %s\n", r.PostInstall)
	}
	if r.Remove != nil {
		s += fmt.Sprintf("  Remove: %s\n", r.Remove)
	}
	if r.PostRemove != nil {
		s += fmt.Sprintf("  PostRemove: %s\n", r.PostRemove)
	}
	return strings.TrimRight(s, "\n")
}

func tasksPrefix(prefix string, r *Role) string {
	s := ""
	if r.Dirs != nil {
		s += fmt.Sprintf("%sDirs: %s\n", prefix, r.Dirs)
	}
	if r.Files != nil {
		s += fmt.Sprintf("%sFiles: %s\n", prefix, r.Files)
	}
	if r.Lines != nil {
		s += fmt.Sprintf("%sLines: %s\n", prefix, r.Lines)
	}
	if r.Links != nil {
		s += fmt.Sprintf("%sLinks: %s\n", prefix, r.Links)
	}
	if r.Templates != nil {
		s += fmt.Sprintf("%sTemplates: %s\n", prefix, r.Templates)
	}
	return s
}

// LoadConfig ...
func (r *Role) LoadConfig() error {
	cfgPath := filepath.Join(r.Path, ".dot.yml")
	role, err := LoadRole(cfgPath)
	if err != nil {
		return fmt.Errorf("%s: %s", r.Name, err)
	}
	// fmt.Printf("MERGE %+v\n", role.Env)
	return mergo.Merge(r, role)
}

// Prepare role tasks
func (r *Role) Prepare() error {
	if r.Path == "" {
		r.Path = filepath.Join(os.ExpandEnv("$HOME"), ".dot", r.Name)
	}
	// fmt.Println("prepare", r.Name)
	// for _, t := range r.Tasks {
	// 	if err := t.Prepare(); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

// // RunSync role tasks
// func (r *Role) RunSync() error {
// 	fmt.Println("Sync", r.Name)
// 	return nil
// }

// // RunInstall role tasks
// func (r *Role) RunInstall() error {
// 	fmt.Println("Install", r.Name)
// 	return nil
// }

// // RunRemove role tasks
// func (r *Role) RunRemove() error {
// 	fmt.Println("Remove", r.Name)
// 	return nil
// }

// NewRole ...
func NewRole() *Role {
	return &Role{}
}

// LoadRole ...
func LoadRole(path string) (Role, error) {
	rc := &RoleConfig{}
	data, err := Read(path)
	if err != nil {
		return rc.Role, err
	}
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook:       roleDecodeHook,
		WeaklyTypedInput: true,
		Result:           &rc,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return rc.Role, err
	}
	if err := decoder.Decode(data); err != nil {
		return rc.Role, err
	}
	// if err := mapstructure.WeakDecode(data, &rc); err != nil {
	// 	return rc.Role, err
	// }
	return rc.Role, nil
}

func roleDecodeHook(f reflect.Type, t reflect.Type, i interface{}) (interface{}, error) {
	// input := i.(map[string]interface{})
	// f == reflect.TypeOf("")
	switch val := i.(type) {
	case string:
		switch {
		// case t == reflect.TypeOf((*Dirs)(nil)):
		case t == reflect.TypeOf((*Dir)(nil)):
			i = &Dir{Path: val}
		case t == reflect.TypeOf((*Hook)(nil)):
			i = &Hook{Command: val}
		case t == reflect.TypeOf((*Pkg)(nil)):
			i = &Pkg{Name: val}
		case t == reflect.TypeOf((*Link)(nil)):
			i = &Link{Source: val}
		case t == reflect.TypeOf((*Template)(nil)):
			i = &Template{Source: val}
			// default:
			// 	fmt.Println("sss", val)
		}
		// case map[interface{}]interface{}:
		// 	// case map[string]interface{}:
		// 	switch {
		// 	case t == reflect.TypeOf((*Line)(nil)):
		// 		fmt.Println("LINE", val)
		// 		lines := &Lines{} // []*Line{}
		// 		for k, v := range val {
		// 			*lines = append(*lines, &Line{
		// 				File: k.(string),
		// 				Line: v.(string),
		// 			})
		// 		}
		// 		i = *lines
		// 		fmt.Println("=======", i)
		// 	case t == reflect.TypeOf((*Lines)(nil)):
		// 		fmt.Println("LINES", val)
		// 	}
		// default:
		// 	fmt.Println("->", f, t)
	}
	// switch t {
	// case reflect.TypeOf(&Dir{}):
	// 	fmt.Println("DIR", t, "=>", i)
	// default:
	// 	fmt.Println("???", t, "=>", reflect.TypeOf(i))
	// }
	return i, nil
}
