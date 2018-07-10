package dot

import (
	"fmt"
	"os"
	"path/filepath"
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
	Env Env
	Copy Paths
	Link Paths
	Template Paths
}

// Env ...
type Env map[string]string

// Paths ...
type Paths map[string]string

// ErrEmptyRole ...
var ErrEmptyRole = fmt.Errorf("Attempt to register an empty role")

// NewRole ...
func NewRole(name string) *Role {
	switch name {
	case "":
		name = "default"
		break
	case "all":
		name = "*"
		break
	}
	url := ""
	if strings.Contains(name, ":") {
		parts := strings.Split(name, ":")
		if len(parts) == 2 {
			name = parts[0]
			url = parts[1]
		} else {
			fmt.Println("Unhandled role spec", name)
			os.Exit(1)
		}
	}
	// if strings.Contains(name, "*") {
	// 	// find glob
	// }
	r := &Role{Name: name, URL: url}
	r.Parse()
	return r
}

// Parse ...
func (r *Role) Parse() *Role {
	r.URL = ParseURL(r.URL)
	// if r.Path == "" {
	// }
	return r
}

// Register ...
func (r *Role) Register(cfg *Config) error {
	if (&Role{}) == r {
		return ErrEmptyRole
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
func (r *Role) Init(target, roleDir string) error {
	if r.Path == "" {
		r.Path = filepath.Join(target, roleDir, r.Name)
	}
	r.Path = os.ExpandEnv(r.Path)
	fmt.Printf("Role [%s] %s (%s)\n", r.Name, r.Path, r.URL)
	fmt.Printf("Copies: %+v\n", r.Copy)
	// fmt.Printf("Links: %+v\n", r.Link)
	for s, t := range r.Link {
		fmt.Println("ln -s", s, t)
	}
	fmt.Printf("Templates: %+v\n", r.Template)
	return nil
}

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

// ParsePath ...
func ParsePath(s string) (Paths, error) {
	src := s
	dst := ""
	if strings.Contains(s, ":") {
		parts := strings.Split(s, ":")
		if len(parts) == 2 {
			src = parts[0]
			dst = parts[1]
		} else {
			fmt.Println("Unhandled path spec", s)
			os.Exit(1)
		}
	}
	paths := map[string]string{}
	if strings.Contains(src, "*") {
		glob, err := filepath.Glob(src)
		if err != nil {
			return paths, err
		}
		for _, s := range glob {
			_, f := filepath.Split(s)
			dst = filepath.Join(dst, f)
			paths[s] = dst
		}
	} else {
		paths[src] = dst
	}
	return paths, nil
}
