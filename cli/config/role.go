package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/LEI/dot/cli/config/tasks"
	"github.com/LEI/dot/cli/config/types"
	"github.com/LEI/dot/pkg/git"
	"github.com/LEI/dot/system"
	"github.com/imdario/mergo"
)

// type Roles []*Role
// func (roles *Roles) list() { }

// Role structure
type Role struct {
	Name  string
	Path  string
	URL   string
	OS    types.Slice // tasks.OS
	Env   types.Map   // tasks.Env
	Deps  types.Slice `mapstructure:"dependencies"`
	Dirs  tasks.Dirs  `mapstructure:"dir"`
	Files tasks.Files `mapstructure:"copy"`
	Links tasks.Links `mapstructure:"link"`
	// Template interface{} // []*tasks.Template

	synced bool
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

// Sync role
func (r *Role) Sync() error {
	repo, err := git.NewRepo(r.Path, r.URL)
	if err != nil {
		return err
	}
	exists, err := system.IsDir(r.Path)
	if err != nil {
		return err
	}
	if exists {
		// fmt.Fprintf(dotCli.Out(), "Checking %s...\n", name)
		if err := repo.Status(); err != nil {
			return err
		}
		if err := repo.Pull(); err != nil {
			return err
		}
	} else {
		// fmt.Fprintf(dotCli.Out(), "Cloning %s into %s...\n", name, repo.Dir)
		if err := repo.Clone(); err != nil {
			return err
		}
	}
	return nil
}

// Parse role
func (r *Role) Parse(i interface{}) error {
	switch v := i.(type) {
	case map[string]string:
		r.Name = v["name"]
		r.URL = v["url"]
		r.Path = v["path"]
		r.OS.Parse(v["os"])
		r.Env.Parse(v["env"])
		r.Deps.Parse(v["dependencies"])
		r.Dirs.Parse(v["dir"])
		r.Files.Parse(v["copy"])
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
		r.Env.Parse(v["env"])
		r.Deps.Parse(v["dependencies"])
		r.Dirs.Parse(v["dir"])
		r.Files.Parse(v["copy"])
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
		r.Env.Parse(v["env"])
		r.Deps.Parse(v["dependencies"])
		r.Dirs.Parse(v["dir"])
		r.Files.Parse(v["copy"])
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
	if err := r.PrepareFiles(target); err != nil {
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

// PrepareFiles role
func (r *Role) PrepareFiles(target string) error {
	files := &tasks.Files{}
	for _, f := range r.Files {
		src := os.ExpandEnv(f.Source)
		dst := os.ExpandEnv(f.Target)
		if !filepath.IsAbs(src) {
			src = filepath.Join(r.Path, src)
		}
		paths, err := preparePaths(target, src, dst)
		if err != nil {
			return err
		}
		for k, v := range paths {
			ff := f
			// ff := &tasks.Copy{}
			// if err := mergo.Merge(ff, f); err != nil {
			// 	return err
			// }
			ff.Source = k
			ff.Target = v
			files.Add(*ff)
		}
		// dir := os.ExpandEnv(f.Path)
		// if !filepath.IsAbs(dir) {
		// 	dir = filepath.Join(target, dir)
		// }
		// f.Path = dir
		// files.Add(f)
	}
	r.Files = *files
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
		paths, err := preparePaths(target, src, dst)
		if err != nil {
			return err
		}
		for k, v := range paths {
			ll := l
			// ll := &tasks.Link{}
			// if err := mergo.Merge(ll, l); err != nil {
			// 	return err
			// }
			ll.Source = k
			ll.Target = v
			links.Add(*ll)
		}
	}
	r.Links = *links
	return nil
}

func preparePaths(target, src, dst string) (map[string]string, error) {
	ret := map[string]string{}
	//*links = append(*links, l)
	if strings.Contains(src, "*") {
		// fmt.Println("*", src, dst)
		glob, err := filepath.Glob(src)
		if err != nil {
			return ret, err
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
				return ret, err
			}
			ret[s] = t
		}
	} else {
		t, err := prepareTarget(target, src, dst)
		if err != nil {
			return ret, err
		}
		ret[src] = t
	}
	return ret, nil
}

func prepareTarget(target, src, dst string) (string, error) {
	//fmt.Println("+", src, dst)
	_, f := filepath.Split(src)
	if f == "" {
		return "", fmt.Errorf("error (no source file name) while parsing: %s / %s", src, dst)
	}
	if !filepath.IsAbs(dst) {
		dst = filepath.Join(target, dst)
	}
	// if _, err := dotfile.CreateDir(baseDir); err != nil {
	// 	return baseDir, err
	// }
	t := filepath.Join(dst, f)
	return t, nil
}
