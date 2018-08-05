package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/imdario/mergo"

	"github.com/LEI/dot/dotfile"
	"github.com/LEI/dot/parsers"
	"github.com/LEI/dot/utils"
)

// r.<Task>, r.Register<Task>
var defaultTasks = []string{
	"copy",
	"exec",
	"line",
	"link",
	"package",
	"template",
}

var ignore = []string{
	"*.json",
	"*.md",
	"*.yml",
	".git",
}

// RoleConfig ...
type RoleConfig struct {
	// TODO: Use top-level struct
	// for meta properties (dir, dst)
	// or change config format role:
	Role *Role
}

// Role ...
type Role struct {
	Name      string   // Name of the role
	Path      string   // Local directory
	URL       string   // Repository URL
	OS        []string // Allowed OSes
	Env       map[string]string
	Vars      map[string]interface{}
	Copies    parsers.Map       `yaml:"copy"`
	Lines     map[string]string `yaml:"line"`
	Links     parsers.Map       `yaml:"link"`
	Templates parsers.Templates `yaml:"template"`

	// Hooks
	Install     parsers.Commands
	PostInstall parsers.Commands `yaml:"post_install"`
	Remove      parsers.Commands
	PostRemove  parsers.Commands `yaml:"post_remove"`

	Pkg     parsers.Packages
	Deps    []string `yaml:"dependencies"`
	Enabled bool     // TODO `default:"true"`
}

// // Env ...
// type Env map[string]string

// // Copy ...
// type Copy struct {
// 	*parsers.Paths
// 	Format string
// }

// // Link ...
// type Link struct {
// 	*parsers.Paths
// 	Format string
// }

// // Template ...
// type Template struct {
// 	*parsers.Paths
// 	Format string
// }

// ErrEmptyRole ...
var ErrEmptyRole = fmt.Errorf("attempt to register an empty role")

// NewRole ...
func NewRole(name string) *Role {
	if name == "" {
		fmt.Println("No role name!!!")
		os.Exit(1)
	}
	// TODO switch name
	// "" -> default
	// "all" -> *
	url := ""
	// FIXME https://(url) can contain ":"
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
	return &Role{
		Name: name,
		URL:  url,
		// FIXME Enabled: true,
	}
}

// func (r *Role) String() string {
// 	return fmt.Sprintf("%s", r.Name)
// }

// PrintDeps ...
func (r *Role) PrintDeps() (s string) {
	if r.Deps != nil && len(r.Deps) > 0 {
		s += fmt.Sprintf("Deps: %d", len(r.Deps))
		for _, v := range r.Deps {
			s += fmt.Sprintf("  %+v", v)
		}
	}
	return s
}

// Register ...
func (r *Role) Register(cfg *Config) error {
	if (&Role{}) == r {
		return ErrEmptyRole
	}
	cfg.AddRole(r)
	return nil
}

// ApplyDeps ...
// func (r *Role) ApplyDeps(cfg *Config) error {
// 	for _, n := range r.Deps {
// 		fmt.Println("DEP", n)
// 	}
// 	return nil
// }

// Merge ...
func (r *Role) Merge(role *Role) error {
	// vr := reflect.ValueOf(r).Elem()
	// vrole := reflect.ValueOf(role).Elem()
	// fmt.Printf("%+v /// %+v\n", vr.Kind(), vrole.Kind())
	// reflect.TypeOf(r), reflect.TypeOf(role)
	// fmt.Printf("%+v\n%+v\n", r, role)
	return mergo.Merge(r, role)
}

// RegisterTask ...
func (r *Role) RegisterTask(name, s string) error {
	v := r.GetField(name)
	i := v.Interface()
	// switch t := i.(type) {
	// case Copy:
	// 	i = i.(Copy)
	// case Link:
	// 	i = i.(Link)
	// case Template:
	// 	i = i.(Template)
	// default:
	// 	return fmt.Errorf("??? %s", t)
	// }
	// if i == nil {
	// 	i = map[string]string{}
	// }
	// if v.IsValid() {
	// 	v.SetInterface(i)
	// }
	fmt.Println("REGISTER TASK", v, "i=", i)
	// paths := f.Interface().(map[string]string)
	// if paths == nil {
	// 	paths = map[string]string{}
	// }
	// src, dst := splitPath(s)
	// paths[src] = dst
	return nil
}

// RegisterCopy ...
func (r *Role) RegisterCopy(s string) error {
	// if r.Copies == nil {
	// 	r.Copies = map[string]string{}
	// }
	// src, dst := splitPath(s)
	// r.Copies[src] = dst
	return nil
}

// RegisterLink ...
func (r *Role) RegisterLink(s string) error {
	if r.Links == nil {
		r.Links = *&parsers.Map{} // map[string]string{}
	}
	src, dst := splitPath(s)
	// fmt.Println("RegisterLink", s, "=", src, "+", dst)
	r.Links.Add(src, dst)
	// r.Links[src] = dst
	return nil
}

// RegisterTemplate ...
func (r *Role) RegisterTemplate(s string) error {
	// if r.Template == nil {
	// 	r.Template = map[string]string{}
	// }
	// src, dst := splitPath(s)
	// r.Template[src] = dst
	return nil
}

// Init ...
func (r *Role) Init() error {
	if !utils.Exist(target) {
		return fmt.Errorf("directory does not exist: %s", target)
	}
	if r.Path == "" {
		r.Path = filepath.Join(target, Options.RoleDir, r.Name)
	}
	// if r.Install == nil {
	// 	r.Install = make([]string, 0)
	// }
	// if r.PostInstall == nil {
	// 	r.PostInstall = make([]string, 0)
	// }
	// if r.Remove == nil {
	// 	r.Remove = make([]string, 0)
	// }
	// if r.PostRemove == nil {
	// 	r.PostRemove = make([]string, 0)
	// }
	// r.URL = ParseURL(r.URL)
	return nil
}

// Sync ...
func (r *Role) Sync() error {
	if r.URL == "" && !utils.Exist(r.Path) {
		return fmt.Errorf("# Role %s has no URL and could not be found in %s", r.Name, r.Path)
	}
	if !r.IsEnabled() {
		// if !utils.Exist(r.Path) { }
		return fmt.Errorf("not enabled: %s", r.Name)
	}
	repo := NewRepo(r.Path, r.URL)
	repoExists := utils.Exist(repo.Path)
	// Clone if the local directory does not exist
	if !repoExists {
		switch err := repo.Clone(); err {
		case nil:
			repoExists = true
		case ErrNetworkUnreachable:
			if !Options.NoSync {
				return err
			}
		default:
			return err
		}
	}
	switch err := repo.checkRepo(); err {
	case nil:
	case ErrNoGitDir:
		fmt.Println("CHECK REPO", repoExists, repo.Path)
		if !repoExists {
			return err
		}
		// Existing directory
		// but no .git: break
		// before exit status 128
		fmt.Printf("Using local %s: %s (no .git)\n", r.Name, r.Path)
		return nil
	case ErrDirtyRepo:
		if !Options.Force {
			return err
		}
	default:
		return err
	}
	switch err := repo.Pull(); err {
	case nil:
	case ErrNetworkUnreachable:
		if !Options.NoSync {
			return err
		}
	default:
		return err
	}
	return nil
}

// ReadConfig ...
func (r *Role) ReadConfig(name string) (string, error) {
	if r.Path == "" || name == "" {
		return "", nil
	}
	// if !filepath.IsAbs(name) {
	cfgPath := filepath.Join(r.Path, name)
	// }
	if !utils.Exist(cfgPath) {
		fmt.Printf("No role config file found: %s\n", cfgPath)
		return "", nil
	}
	cfg, err := utils.Read(cfgPath)
	if err != nil {
		return cfgPath, err
	}
	rc := &RoleConfig{} // Dir: r.Path
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
	}
	return cfgPath, nil // err
}

// IsEnabled ...
func (r *Role) IsEnabled() bool {
	return r.Enabled
}

// Enable ...
func (r *Role) Enable() error {
	if r.Enabled == true {
		return fmt.Errorf("already enabled: %s", r.Name)
	}
	r.Enabled = true
	return nil
}

// Disable ...
func (r *Role) Disable() error {
	if r.Enabled == false {
		return fmt.Errorf("already disabled: %s", r.Name)
	}
	r.Enabled = false
	return nil
}

// Prepare ...
func (r *Role) Prepare() error {
	if err := r.PreparePaths(&r.Copies); err != nil {
		return err
	}
	if err := r.PrepareLines(&r.Lines); err != nil {
		return err
	}
	if err := r.PreparePaths(&r.Links); err != nil {
		return err
	}
	if err := r.PrepareTemplates(&r.Templates); err != nil {
		return err
	}
	return nil
}

// PrepareLines ...
func (r *Role) PrepareLines(l *map[string]string) error {
	lines := make(map[string]string, 0)
	for file, line := range *l {
		// Prepend role directory to source path
		if !filepath.IsAbs(file) {
			file = filepath.Join(target, file)
		}
		file = dotfile.ExpandEnv(file)
		lines[file] = line
	}
	*l = lines
	return nil
}

// PrepareTemplates ...
func (r *Role) PrepareTemplates(t *parsers.Templates) error {
	templates := make(parsers.Templates, 0)
	for _, v := range *t {
		if v.Ext == "" {
			v.Ext = "tpl"
		}
		if v.Target == "" {
			s, t := splitPath(v.Source)
			v.Source = s
			v.Target = t
		}
		// fmt.Println("src", v.Source, "dst", v.Target)
		v.Source = os.ExpandEnv(v.Source)
		v.Target = os.ExpandEnv(v.Target)
		// Prepend role directory to source path
		if !filepath.IsAbs(v.Source) {
			v.Source = filepath.Join(r.Path, v.Source)
		}
		// Check frob globs
		if strings.Contains(v.Source, "*") {
			// fmt.Println("*", v.Source, v.Target)
			glob, err := filepath.Glob(v.Source)
			if err != nil {
				return err
			}
		GLOB:
			for _, s := range glob {
				// Extract source file name
				_, n := filepath.Split(s)
				for _, i := range ignore {
					// Check for ignored patterns
					matched, err := filepath.Match(i, n)
					if err != nil {
						return err
					}
					if matched {
						continue GLOB
					}
				}
				t, err := prepareTarget(s, v.Target)
				if err != nil {
					return err
				}
				v.Source = s
				v.Target = t
				templates.Append(v) // paths[s] = t
			}
		} else {
			t, err := prepareTarget(v.Source, v.Target)
			if err != nil {
				return err
			}
			v.Target = t
			templates.Append(v) // paths[src] = t
		}
	}
	*t = templates
	return nil
}

// PreparePaths ...
func (r *Role) PreparePaths(p *parsers.Map) error {
	// in interface{} p := in.(*parsers.Map)
	// var paths = make(map[string]string, 0)
	paths := make(parsers.Map)
	for src, dst := range *p {
		if dst == "" {
			s, t := splitPath(src)
			src = s
			dst = t
		}
		// fmt.Println("src", src, "dst", dst)
		src = os.ExpandEnv(src)
		dst = os.ExpandEnv(dst)
		// Prepend role directory to source path
		if !filepath.IsAbs(src) {
			src = filepath.Join(r.Path, src)
		}
		// Check frob globs
		if strings.Contains(src, "*") {
			// fmt.Println("*", src, dst)
			glob, err := filepath.Glob(src)
			if err != nil {
				return err
			}
		GLOB:
			for _, s := range glob {
				// Extract source file name
				_, n := filepath.Split(s)
				for _, i := range ignore {
					// Check for ignored patterns
					matched, err := filepath.Match(i, n)
					if err != nil {
						return err
					}
					if matched {
						continue GLOB
					}
				}
				// fmt.Println("PREPARE GLOB", s, "/", dst)
				t, err := prepareTarget(s, dst)
				if err != nil {
					return err
				}
				paths.Add(s, t) // paths[s] = t
			}
		} else {
			// fmt.Println("PREPARE", src, "/", dst)
			t, err := prepareTarget(src, dst)
			if err != nil {
				return err
			}
			paths.Add(src, t) // paths[src] = t
		}
	}
	// *p = *(p.Merge(paths))
	*p = paths
	return nil
}

func prepareTarget(src, dst string) (string, error) {
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

// GetField ...
func (r *Role) GetField(key string) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(r)).FieldByName(key)
}

func splitPath(s string) (src, dst string) {
	parts := filepath.SplitList(s)
	switch len(parts) {
	case 1:
		src = s
	case 2:
		src = parts[0]
		dst = parts[1]
	default:
		fmt.Println("Unhandled path spec", src)
		os.Exit(1)
	}
	// src = s
	// if strings.Contains(src, ":") {
	// 	parts := strings.Split(src, ":")
	// 	if len(parts) == 2 {
	// 		src = parts[0]
	// 		dst = parts[1]
	// 	} else {
	// 		fmt.Println("Unhandled path spec", src)
	// 		os.Exit(1)
	// 	}
	// }
	return src, dst
}

/*
func ParsePath(src, dst string) (parsers.Paths, error) {
	if strings.Contains(src, ":") {
		parts := strings.Split(src, ":")
		if len(parts) == 2 {
			src = parts[0]
			dst = filepath.Join(dst, parts[1])
		} else {
			fmt.Println("Unhandled path spec", src)
			os.Exit(1)
		}
	}
	src = os.ExpandEnv(src)
	dst = os.ExpandEnv(dst)
	paths := map[string]string{}
	if strings.Contains(src, "*") {
		fmt.Println("GLOB", src)
		glob, err := filepath.Glob(src)
		if err != nil {
			return paths, err
		}
		for _, s := range glob {
			_, f := filepath.Split(s)
			t := filepath.Join(dst, f)
			paths[s] = t
		}
	} else {
		_, f := filepath.Split(src)
		t := filepath.Join(dst, f)
		paths[src] = t
	}
	return paths, nil
}
*/
