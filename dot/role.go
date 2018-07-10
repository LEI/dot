package dot

import (
	"fmt"
	// "reflect"
	"strings"

	"github.com/imdario/mergo"
)

// Role ...
type Role struct {
	Name string // Name of the role
	Path string // Local directory
	URL string // Repository URL
	OS []string // TODO
	Copy Paths
	Link Paths
	Template Paths
}

// Paths ...
type Paths map[string]string

// NewRole ...
func NewRole(name, p string) *Role {
	r := &Role{Name: name}
	r.Parse()
	return r
}

// Parse ...
func (r *Role) Parse() *Role {
	r.URL = ParseURL(r.Name)
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
	paths, err := ParsePath(s)
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
	paths, err := ParsePath(s)
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
	paths, err := ParsePath(s)
	if err != nil {
		return err
	}
	for src, dst := range paths {
		r.Template[src] = dst
	}
	return nil
}

// Init ...
func (r *Role) Init() error {
	fmt.Printf("Role [%s] %s (%s)\n", r.Name, r.Path, r.URL)
	fmt.Println("Copies", r.Copy)
	fmt.Println("Links", r.Link)
	fmt.Println("Templates", r.Template)
	return nil
}

// ParseURL ...
func ParseURL(p string) string {
	// if r.Name == "" {}
	// if p == "" {}
	if !strings.Contains(p, "http") {
		base := "https://github.com"
		p = base + "/" + p
	}
	return p
}

// ParsePath ...
func ParsePath(s string) (Paths, error) {
	paths := map[string]string{
		s: s,
	}

	return paths, nil
}
