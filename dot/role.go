package dot

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/imdario/mergo"
)

var tasks = []string{"Copy", "Link", "Template"}

var noCheck bool // Ignore uncomitted changes in repository
var noSync bool // Do not attempt to git clone or pull

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

	// v0.0.2
	Pkg []interface{}
	Line map[string]string
}

// Env ...
type Env map[string]string

// Paths ...
type Paths map[string]string

// UnmarshalYAML ...
func (p *Paths) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Avoid assignment to entry in nil map
	if *p == nil {
		*p = make(Paths)
	}
	var i interface{}
	if err := unmarshal(&i); err != nil {
		return err
	}
	switch val := i.(type) {
	case []string:
		for _, v := range val {
			(*p)[v] = v
		}
		break
	// case interface{}:
	// 	s := val.(string)
	// 	(*p)[s] = s
	case []interface{}:
		for _, v := range val {
			s := v.(string)
			(*p)[s] = s
		}
		break
	case map[string]string:
		p = i.(*Paths)
		break
	case map[interface{}]interface{}:
		for _, v := range val {
			s := v.(string)
			(*p)[s] = s
		}
		break
	default:
		t := reflect.TypeOf(val)
		T := t.Elem()
		if t.Kind() == reflect.Map {
			T = reflect.MapOf(t.Key(), t.Elem())
		}
		return fmt.Errorf("Unable to unmarshal %s into struct: %+v", T, val)
	}
	return nil
}

// // UnmarshalFlag ...
// func (p *Paths) UnmarshalFlag(value string) error {
// 	fmt.Println("UnmarshalFlag", value)
// 	return nil
// }

// // MarshalFlag ...
// func (p Paths) MarshalFlag() (string, error) {
// 	return fmt.Sprintf("MarshalFlag: %+v", p), nil
// }

// ErrEmptyRole ...
var ErrEmptyRole = fmt.Errorf("Attempt to register an empty role")

// NewRole ...
func NewRole(name string) *Role {
	switch name {
	case "":
		fmt.Println("No role name")
		os.Exit(1)
		// name = "default"
		break
	// case "all":
	// 	name = "*"
	// 	break
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
	paths, err := ParsePath(s, r.Path)
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
	paths, err := ParsePath(s, r.Path)
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
	paths, err := ParsePath(s, r.Path)
	if err != nil {
		return err
	}
	for src, dst := range paths {
		r.Template[src] = dst
	}
	return nil
}

// Init ...
func (r *Role) Init(target string) error {
	target = os.ExpandEnv(target)
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return fmt.Errorf("Directory does not exist: %s", target)
	}
	if r.Path == "" {
		r.Path = filepath.Join(target, r.Name)
	}
	// r.URL = ParseURL(r.URL)
	fmt.Printf("# Syncing role %s [%s] %s\n", r.Name, r.Path, r.URL)
	if err := r.Sync(); err != nil {
		return err
	}
	for _, c := range tasks {
		if err := r.InitPaths(c); err != nil {
			return err
		}
	}
	return nil
}

// Sync ...
func (r *Role) Sync() error {
	repo := NewRepo(r.Path, r.URL)
	// Clone if the local directory does not exist
	if _, err := os.Stat(repo.Path); os.IsNotExist(err) {
		switch err := repo.Clone(); err {
		case nil:
			break
		case ErrNetworkUnreachable:
			if !noSync {
				return err
			}
		default:
			return err
		}
	}
	switch err := repo.Clone(); err {
	case nil:
		break
	case ErrNetworkUnreachable:
		if !noSync {
			return err
		}
	default:
		return err
	}
	// TODO: flag ignore dirty
	switch err := repo.checkRepo(); err {
	case nil:
		break
	case ErrDirtyRepo:
		if !noCheck {
			return err
		}
	default:
		return err
	}
	// TODO: skip if just cloned
	switch err := repo.Pull(); err {
	case nil:
		break
	case ErrNetworkUnreachable:
		if !noSync {
			return err
		}
	default:
		return err
	}
	return nil
}

// RoleConfig ...
type RoleConfig struct {
	Role *Role
}

// LoadConfig ...
func (r *Role) LoadConfig(name string) (string, error) {
	if r.Path == "" || name == "" {
		return "", nil
	}
	cfgPath := filepath.Join(r.Path, name)
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		fmt.Println("No role config file found:", cfgPath)
		return "", nil
	}
	cfg, err := readConfig(cfgPath)
	if err != nil {
	    return cfgPath, err
	}
	rc := &RoleConfig{}
	// fmt.Printf("+++\n%v\n+++\n", string(cfg))
	err = yaml.Unmarshal(cfg, &rc)
	// fmt.Printf("---\n%v\n---\n", rc.Role)
	if err != nil {
		// fmt.Fprintf(os.Stderr, "Error while parsing %s:\n%s\n", cfgPath, err)
		return cfgPath, err
	}
	if rc.Role != nil {
		if err := r.Merge(rc.Role); err != nil {
			return cfgPath, err
		}
		// fmt.Printf("---\n%v\n---\n", r)
	}
	return cfgPath, nil // err
}

// GetField ...
func (r *Role) GetField(key string) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(r)).FieldByName(key)
}

// InitPaths ...
func (r *Role) InitPaths(key string) error {
	f := r.GetField(key)
	val := f.Interface().(Paths)
	// fmt.Printf("%s: %+v\n", key, val)
	var paths Paths = make(map[string]string, len(val))
	for s, t := range val {
		s = filepath.Join(r.Path, s)
		// fmt.Printf("%s '%s' '%s'\n", key, s, t)
		paths[s] = t
	}
	r.Copy = paths
	return nil
}

// Execute ...
func (r *Role) Execute() error {
	fmt.Println("# Role", r.Name)
	for s, t := range r.Copy {
		fmt.Printf("cp '%s' '%s'\n", s, t)
	}
	for s, t := range r.Link {
		fmt.Printf("ln -s '%s' '%s'\n", s, t)
	}
	for s, t := range r.Template {
		fmt.Printf("tpl '%s' '%s'\n", s, t)
	}
	return nil
}

// ParseURL ...
// func ParseURL(url string) string {
// 	// if r.Name == "" {}
// 	// if url == "" {}
// 	if !strings.Contains(url, "http") {
// 		base := "https://github.com"
// 		url = base + "/" + url
// 	}
// 	return url
// }

// ParsePath ...
func ParsePath(s, baseDir string) (Paths, error) {
	source := s
	target := baseDir
	if strings.Contains(s, ":") {
		parts := strings.Split(s, ":")
		if len(parts) == 2 {
			source = parts[0]
			target = filepath.Join(target, parts[1])
		} else {
			fmt.Println("Unhandled path spec", s)
			os.Exit(1)
		}
	}
	// fmt.Println("TARGET", target, baseDir)
	paths := map[string]string{}
	if strings.Contains(source, "*") {
		glob, err := filepath.Glob(source)
		if err != nil {
			return paths, err
		}
		for _, s := range glob {
			_, f := filepath.Split(s)
			t := filepath.Join(target, f)
			paths[s] = t
		}
	} else {
		paths[source] = target
	}
	return paths, nil
}
