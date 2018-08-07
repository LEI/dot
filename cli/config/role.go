package config

import (
	"fmt"
	"reflect"

	"github.com/LEI/dot/cli/config/tasks"
)

// type Roles []*Role
// func (roles *Roles) list() { }

// Role structure
type Role struct {
	Name string
	URL string
	OS []string
	Deps *tasks.Deps `mapstructure:"dependencies"`
	// Copy interface{} // []*tasks.Copy
	Link *tasks.Links // []*tasks.Link
	// Template interface{} // []*tasks.Template
}

// NewRole config
func NewRole(i interface{}) *Role {
	role := &Role{
		Deps: &tasks.Deps{},
		Link: &tasks.Links{},
		// Template: &tasks.Templates{},
	}
	switch r := i.(type) {
	case map[string]string:
		role.Name = r["name"]
		role.URL = r["url"]
		role.Deps.Parse(r["dependencies"])
		role.Link.Parse(r["link"])
	case map[string]interface{}:
		role.Name = r["name"].(string)
		role.URL = r["url"].(string)
		role.Deps.Parse(r["dependencies"])
		role.Link.Parse(r["link"])
	case map[interface{}]interface{}:
		role.Name = r["name"].(string)
		role.URL = r["url"].(string)
		role.Deps.Parse(r["dependencies"])
		role.Link.Parse(r["link"])
	default:
		fmt.Println("TODO NewRole type:", reflect.TypeOf(r))
	}
	return role
}

// Status role
func (r *Role) Status() bool {
	return true
}

// Init role
func (r *Role) Init() error {
	return nil
}
