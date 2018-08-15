package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
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
	Name string
	Path string
	URL  string
	OS   types.Slice // tasks.OS
	Env  types.Map   // tasks.Env
	// Vars  types.Map
	// IncludeVars types.IncludeMap
	Deps     types.Slice    `mapstructure:"dependencies"`
	Packages tasks.Packages `mapstructure:"pkg"`

	Dirs      tasks.Dirs      `mapstructure:"dir"`
	Files     tasks.Files     `mapstructure:"copy"`
	Links     tasks.Links     `mapstructure:"link"`
	Templates tasks.Templates `mapstructure:"template"`
	Lines     tasks.Lines     `mapstructure:"line"`

	// Hooks
	Install     tasks.Commands
	PostInstall tasks.Commands `mapstructure:"post_install"`
	Remove      tasks.Commands
	PostRemove  tasks.Commands `mapstructure:"post_remove"`

	Ignore []string
	Target string

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
		r.Packages.Parse(v["pkg"])

		r.Dirs.Parse(v["dir"])
		r.Files.Parse(v["copy"])
		r.Links.Parse(v["link"])
		r.Templates.Parse(v["template"])
		r.Lines.Parse(v["line"])

		r.Install.Parse(v["install"])
		r.PostInstall.Parse(v["post_install"])
		r.Remove.Parse(v["install"])
		r.PostRemove.Parse(v["post_install"])
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
		r.Packages.Parse(v["pkg"])

		r.Dirs.Parse(v["dir"])
		r.Files.Parse(v["copy"])
		r.Links.Parse(v["link"])
		r.Templates.Parse(v["template"])
		r.Lines.Parse(v["line"])

		r.Install.Parse(v["install"])
		r.PostInstall.Parse(v["post_install"])
		r.Remove.Parse(v["remove"])
		r.PostRemove.Parse(v["post_remove"])
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
		r.Packages.Parse(v["pkg"])

		r.Dirs.Parse(v["dir"])
		r.Files.Parse(v["copy"])
		r.Links.Parse(v["link"])
		r.Templates.Parse(v["template"])
		r.Lines.Parse(v["line"])

		r.Install.Parse(v["install"])
		r.PostInstall.Parse(v["post_install"])
		r.Remove.Parse(v["remove"])
		r.PostRemove.Parse(v["post_remove"])
	default:
		return fmt.Errorf("TODO NewRole type: %s", reflect.TypeOf(v))
	}
	return nil
}

// Prepare role
func (r *Role) Prepare() error {
	if err := r.PrepareDirs(); err != nil {
		return err
	}
	if err := r.PrepareFiles(); err != nil {
		return err
	}
	if err := r.PrepareLinks(); err != nil {
		return err
	}
	if err := r.PrepareTemplates(); err != nil {
		return err
	}
	if err := r.PrepareLines(); err != nil {
		return err
	}
	return nil
}

// PrepareDirs role
func (r *Role) PrepareDirs() error {
	dirs := &tasks.Dirs{}
	for _, d := range r.Dirs {
		dir := os.ExpandEnv(d.Path)
		if !filepath.IsAbs(dir) {
			dir = filepath.Join(r.Target, dir)
		}
		d.Path = dir
		dirs.Add(d)
	}
	r.Dirs = *dirs
	return nil
}

// PrepareFiles role
func (r *Role) PrepareFiles() error {
	files := &tasks.Files{}
	for _, f := range r.Files {
		src := os.ExpandEnv(f.Source)
		dst := os.ExpandEnv(f.Target)
		if !filepath.IsAbs(src) {
			src = filepath.Join(r.Path, src)
		}
		paths, err := preparePaths(r, src, dst)
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
		// 	dir = filepath.Join(r.Target, dir)
		// }
		// f.Path = dir
		// files.Add(f)
	}
	r.Files = *files
	return nil
}

// PrepareLinks role
func (r *Role) PrepareLinks() error {
	links := &tasks.Links{}
	for _, l := range r.Links {
		src := os.ExpandEnv(l.Source)
		dst := os.ExpandEnv(l.Target)
		if !filepath.IsAbs(src) {
			src = filepath.Join(r.Path, src)
		}
		paths, err := preparePaths(r, src, dst)
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

// PrepareTemplates role
func (r *Role) PrepareTemplates() error {
	templates := &tasks.Templates{}
	for _, t := range r.Templates {
		src := os.ExpandEnv(t.Source)
		dst := os.ExpandEnv(t.Target)
		if !filepath.IsAbs(src) {
			src = filepath.Join(r.Path, src)
		}
		paths, err := preparePaths(r, src, dst)
		if err != nil {
			return err
		}
		for k, v := range paths {
			tt := t
			// tt := &tasks.Template{}
			// if err := mergo.Merge(tt, t); err != nil {
			// 	return err
			// }
			tt.Source = k
			tt.Target = v
			templates.Add(*tt)
		}
	}
	r.Templates = *templates
	return nil
}

// PrepareLines role
func (r *Role) PrepareLines() error {
	lines := &tasks.Lines{}
	for _, l := range r.Lines {
		dst := os.ExpandEnv(l.File)
		if !filepath.IsAbs(dst) {
			dst = filepath.Join(r.Target, dst)
		}
		l.File = dst
		// l.Line = l.Line
		lines.Add(*l)
	}
	r.Lines = *lines
	return nil
}

func preparePaths(r *Role, src, dst string) (map[string]string, error) {
	ret := map[string]string{}
	//*links = append(*links, l)
	if hasMeta(src) { // strings.Contains(src, "*")
		// fmt.Println("*", src, dst)
		glob, err := filepath.Glob(src)
		if err != nil {
			return ret, err
		}
	GLOB:
		for _, s := range glob {
			// Extract source file name
			_, n := filepath.Split(s)
			for _, i := range r.Ignore {
				// Check for ignored patterns
				matched, err := filepath.Match(i, n)
				if err != nil {
					return ret, err
				}
				if matched {
					continue GLOB
				}
			}
			// fmt.Println("PREPARE GLOB", s, "/", dst)
			t, err := prepareTarget(r, s, dst)
			if err != nil {
				return ret, err
			}
			ret[s] = t
		}
	} else {
		t, err := prepareTarget(r, src, dst)
		if err != nil {
			return ret, err
		}
		ret[src] = t
	}
	return ret, nil
}

func prepareTarget(r *Role, src, dst string) (string, error) {
	//fmt.Println("+", src, dst)
	_, name := filepath.Split(src)
	if name == "" {
		return "", fmt.Errorf("no source file name for src / dst: %s / %s", src, dst)
	}
	if !filepath.IsAbs(dst) {
		dst = filepath.Join(r.Target, dst)
	}
	// if _, err := dotfile.CreateDir(baseDir); err != nil {
	// 	return baseDir, err
	// }
	// if isDir, _ := system.IsDir(dst); !isDir {
	// 	// Look for future directories
	// 	ok := false
	// 	for _, d := range r.Dirs {
	// 		// _, n := filepath.Split(d.Path)
	// 		n := strings.TrimPrefix(d.Path, r.Target+system.Separator)
	// 		fmt.Printf("TODO %s == %s / %s\n", n, name, r.Target)
	// 		if n == name {
	// 			ok = true
	// 			break
	// 		}
	// 	}
	// 	if !ok {
	// 		return dst, fmt.Errorf("%s: target directory does not exist and will not be created", dst)
	// 	}
	// }
	dst = filepath.Join(dst, name)
	return dst, nil
}

// Check magix chars recognized by Match
func hasMeta(path string) bool {
	magicChars := `*?[`
	if runtime.GOOS == "windows" {
		magicChars = `*?[\`
	}
	return strings.ContainsAny(path, magicChars)
}
