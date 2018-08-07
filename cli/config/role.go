package config

import (
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/LEI/dot/cli/config/tasks"
	"github.com/imdario/mergo"
)

// type Roles []*Role
// func (roles *Roles) list() { }

// Role structure
type Role struct {
	Name string
	Dir string
	URL string
	OS *tasks.OS
	Deps *tasks.Deps `mapstructure:"dependencies"`
	// Copy interface{} // []*tasks.Copy
	Link *tasks.Links // []*tasks.Link
	// Template interface{} // []*tasks.Template
}

// NewRole config
func NewRole(i interface{}) (*Role, error) {
	r := &Role{
		// Name: "",
		// Dir: "",
		// URL: "",
	}
	if err := r.Parse(i); err != nil {
		return r, err
	}
	if r.Name == "" {
		return r, fmt.Errorf("missing name in role: %+v", r)
	}
	if r.Dir == "" {
		r.Dir = filepath.Join("/tmp/home", ".dot", r.Name)
	}
	return r, nil
}

// Parse role
func (r *Role) Parse(i interface{}) error {
	if r.OS == nil {
		r.OS = &tasks.OS{}
	}
	if r.Deps == nil {
		r.Deps = &tasks.Deps{}
	}
	if r.Link == nil {
		r.Link = &tasks.Links{}
	}
	switch v := i.(type) {
	case map[string]string:
		r.Name = v["name"]
		r.URL = v["url"]
		r.Dir = v["dir"]
		r.OS.Parse(v["os"])
		r.Deps.Parse(v["dependencies"])
		r.Link.Parse(v["link"])
	case map[string]interface{}:
		if name, ok := v["name"].(string); ok {
			r.Name = name
		}
		if dir, ok := v["dir"].(string); ok {
			r.Dir = dir
		}
		if url, ok := v["url"].(string); ok {
			r.URL = url
		}
		r.OS.Parse(v["os"])
		r.Deps.Parse(v["dependencies"])
		r.Link.Parse(v["link"])
	case map[interface{}]interface{}:
		if name, ok := v["name"].(string); ok {
			r.Name = name
		}
		if dir, ok := v["dir"].(string); ok {
			r.Dir = dir
		}
		if url, ok := v["url"].(string); ok {
			r.URL = url
		}
		r.OS.Parse(v["os"])
		r.Deps.Parse(v["dependencies"])
		r.Link.Parse(v["link"])
	default:
		return fmt.Errorf("TODO NewRole type: %s", reflect.TypeOf(v))
	}
	return nil
}

// Merge role
func (r *Role) Merge(i interface{}) error {
	role := &Role{}
	if err := role.Parse(i); err != nil {
		return err
	}
	// var role *Role
	// switch v := i.(type) {
	// case *Role:
	// 	role = v // .(*Role)
	// default:
	// 	return fmt.Errorf("?: %s", reflect.TypeOf(v))
	// }
	return mergo.Merge(r, role)
}

// // Init role
// func (r *Role) Init() error {
// 	return nil
// }

// Status role
func (r *Role) Status() bool {
	return true
}

// Sync role
func (r *Role) Sync() bool {
	return true
}
