package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/LEI/dot/cli/config/tasks"
	"github.com/imdario/mergo"
)

// type Roles []*Role
// func (roles *Roles) list() { }

// Role structure
type Role struct {
	Name string
	Path string
	URL string
	OS tasks.OS
	Deps tasks.Deps `mapstructure:"dependencies"`
	Dirs tasks.Dirs `mapstructure:"dir"`
	// Copy interface{} // []*tasks.Copy
	Links tasks.Links `mapstructure:"link"`
	// Template interface{} // []*tasks.Template
}

// NewRole config
func NewRole(i interface{}) (*Role, error) {
	r := &Role{
		// Name: "",
		// Path: "",
		// URL: "",
	}
	if err := r.Parse(i); err != nil {
		return r, err
	}
	if r.Name == "" {
		return r, fmt.Errorf("missing name in role: %+v", r)
	}
	if r.Path == "" {
		r.Path = filepath.Join("/tmp/home", ".dot", r.Name)
	}
	return r, nil
}

// // Init role
// func (r *Role) Init() error {
// 	return nil
// }

// // Status role
// func (r *Role) Status() bool {
// 	return true
// }

// // Sync role
// func (r *Role) Sync() bool {
// 	return true
// }

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

// Parse role
func (r *Role) Parse(i interface{}) error {
	// if r.OS == nil {
	// 	r.OS = tasks.OS{}
	// }
	// if r.Deps == nil {
	// 	r.Deps = tasks.Deps{}
	// }
	// if r.Dirs == nil {
	// 	r.Dirs = tasks.Dirs{}
	// }
	if r.Links == nil {
		r.Links = tasks.Links{}
	}
	switch v := i.(type) {
	case map[string]string:
		r.Name = v["name"]
		r.URL = v["url"]
		r.Path = v["path"]
		r.OS.Parse(v["os"])
		r.Deps.Parse(v["dependencies"])
		r.Dirs.Parse(v["dir"])
		r.Links.Parse(v["link"])
	case map[string]interface{}:
		if name, ok := v["name"].(string); ok {
			r.Name = name
		}
		if dir, ok := v["path"].(string); ok {
			r.Path = dir
		}
		if url, ok := v["url"].(string); ok {
			r.URL = url
		}
		r.OS.Parse(v["os"])
		r.Deps.Parse(v["dependencies"])
		r.Dirs.Parse(v["dir"])
		r.Links.Parse(v["link"])
	case map[interface{}]interface{}:
		if name, ok := v["name"].(string); ok {
			r.Name = name
		}
		if path, ok := v["path"].(string); ok {
			r.Path = path
		}
		if url, ok := v["url"].(string); ok {
			r.URL = url
		}
		r.OS.Parse(v["os"])
		r.Deps.Parse(v["dependencies"])
		r.Dirs.Parse(v["dir"])
		r.Links.Parse(v["link"])
	default:
		return fmt.Errorf("TODO NewRole type: %s", reflect.TypeOf(v))
	}
	return nil
}

// Prepare role
func (r *Role) Prepare(target string) error {
	if err := r.PrepareDirs(target); err != nil {
		return err
	}
	if err := r.PrepareLinks(target); err != nil {
		return err
	}
	return nil
}

// PrepareDirs role
func (r *Role) PrepareDirs(target string) error {
	dirs := &tasks.Dirs{}
	for _, d := range r.Dirs {
		dir := os.ExpandEnv(d.Path)
		if !filepath.IsAbs(dir) {
			dir = filepath.Join(target, dir)
		}
		d.Path = dir
		dirs.Add(d)
	}
	r.Dirs = *dirs
	return nil
}

// PrepareLinks role
func (r *Role) PrepareLinks(target string) error {
	links := &tasks.Links{}
	for _, l := range r.Links {
		src := os.ExpandEnv(l.Source)
		dst := os.ExpandEnv(l.Target)
		if !filepath.IsAbs(src) {
			src = filepath.Join(r.Path, src)
		}
		//*links = append(*links, l)
		if strings.Contains(src, "*") {
			// fmt.Println("*", src, dst)
			glob, err := filepath.Glob(src)
			if err != nil {
				return err
			}
		// GLOB:
			for _, s := range glob {
				// // Extract source file name
				// _, n := filepath.Split(s)
				// for _, i := range ignore {
				// 	// Check for ignored patterns
				// 	matched, err := filepath.Match(i, n)
				// 	if err != nil {
				// 		return err
				// 	}
				// 	if matched {
				// 		continue GLOB
				// 	}
				// }
				// fmt.Println("PREPARE GLOB", s, "/", dst)
				t, err := prepareTarget(target, s, dst)
				if err != nil {
					return err
				}
				// FIXME copy clone struct
				ll := &tasks.Link{}
				if err := mergo.Merge(ll, l); err != nil {
					return err
				}
				ll.Source = s
				ll.Target = t
				links.Add(ll)
			}
		} else {
			t, err := prepareTarget(target, src, dst)
			if err != nil {
				return err
			}
			l.Source = src
			l.Target = t
			links.Add(l)
		}
	}
	r.Links = *links
	return nil
}

func prepareTarget(dir, src, dst string) (string, error) {
	//fmt.Println("+", src, dst)
	_, f := filepath.Split(src)
	if f == "" {
		return "", fmt.Errorf("error (no source file name) while parsing: %s / %s", src, dst)
	}
	if !filepath.IsAbs(dst) {
		dst = filepath.Join(dir, dst)
	}
	// if _, err := dotfile.CreateDir(baseDir); err != nil {
	// 	return baseDir, err
	// }
	t := filepath.Join(dst, f)
	return t, nil
}
