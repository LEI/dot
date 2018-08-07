package config

import (
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/LEI/dot/cli/config/tasks"
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
		OS: &tasks.OS{},
		Deps: &tasks.Deps{},
		// Copy: &tasks.Map{},
		Link: &tasks.Links{},
		// Template: &tasks.Templates{},
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
		r.Name = v["name"].(string)
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
		r.Name = v["name"].(string)
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
		return r, fmt.Errorf("TODO NewRole type: %s", reflect.TypeOf(v))
	}
	if r.Name == "" {
		return r, fmt.Errorf("missing name in role: %+v", r)
	}
	if r.Dir == "" {
		r.Dir = filepath.Join(targetDir, ".dot", r.Name)
	}
	return r, nil
}
// // Init role
// func (r *Role) Init() error {
// 	return nil
// }

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
