package dot

import (
	// "fmt"
	// "reflect"
	"strings"

	"github.com/imdario/mergo"
)

// Role ...
type Role struct {
	Name string
	URL string
	OS []string
	Copy Paths
	Link Paths
	Template Paths
}

// Paths ...
type Paths map[string]string

// ParseURL ...
func ParseURL(url string) string {
	// if r.Name == "" {}
	// if url == "" {}
	if !strings.Contains(url, "http") {
		base := "https://github.com"
		url = base + "/" + url
	}
	return url
}

// NewRole ...
func NewRole(name, url string) *Role {
	r := &Role{Name: name, URL: url}
	r.URL = ParseURL(r.URL)
	return r
}

// Register ...
func (r *Role) Register(cfg *Config) error {
	if (&Role{}) == r {
		return nil
	}
	for i, cfgRole := range cfg.Roles {
		if cfgRole.Name == r.Name {
			if err := r.Merge(cfgRole); err != nil {
				return err
			}
			cfg.Roles[i] = r
			return nil
			// break
		}
	}
	cfg.Roles = append(cfg.Roles, r)
	return nil
}

// Merge ...
func (r *Role) Merge(role *Role) error {
	// vr := reflect.ValueOf(r).Elem()
	// vrole := reflect.ValueOf(role).Elem()
	// fmt.Printf("%+v /// %+v\n", vr.Kind(), vrole.Kind())
	// reflect.TypeOf(r), reflect.TypeOf(role)
	// fmt.Printf("%+v\n%+v\n", r, role)
	return mergo.Merge(r, role)
}

// RegisterCopy ...
func (r *Role) RegisterCopy(s string) error {
	if r.Copy == nil {
		r.Copy = map[string]string{}
	}
	paths, err := parseString(s)
	if err != nil {
		return err
	}
	for src, dst := range paths {
		r.Copy[src] = dst
	}
	return nil
}

// RegisterLink ...
func (r *Role) RegisterLink(s string) error {
	if r.Link == nil {
		r.Link = map[string]string{}
	}
	paths, err := parseString(s)
	if err != nil {
		return err
	}
	for src, dst := range paths {
		r.Link[src] = dst
	}
	return nil
}

// RegisterTemplate ...
func (r *Role) RegisterTemplate(s string) error {
	if r.Template == nil {
		r.Template = map[string]string{}
	}
	paths, err := parseString(s)
	if err != nil {
		return err
	}
	for src, dst := range paths {
		r.Template[src] = dst
	}
	return nil
}

func parseString(s string) (Paths, error) {
	paths := map[string]string{
		s: s,
	}

	return paths, nil
}
