package cmd

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
)

// r.<Task>, r.Register<Task>
var defaultTasks = []string{
	"copy",
	"exec",
	"line",
	"link",
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
	Name     string   // Name of the role
	Path     string   // Local directory
	URL      string   // Repository URL
	OS       []string // Allowed OSes
	Env      Env
	Copy     parsers.Paths
	Line     map[string]string
	Link     parsers.Paths
	Template parsers.Paths

	// Hooks
	Install     []string
	PostInstall []string `yaml:"post_install"`
	Remove      []string
	PostRemove  []string `yaml:"post_remove"`

	Pkg          parsers.Packages
	Dependencies []string
	Enabled      bool // TODO `default:"true"`
}

// Env ...
type Env map[string]string

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
var ErrEmptyRole = fmt.Errorf("Attempt to register an empty role")

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

func (r *Role) String() string {
	if Verbose > 2 {
		return fmt.Sprintf("[%s](%s)%s", r.Name, r.URL, r.Path)
	}
	return fmt.Sprintf("%s", r.Name)
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
	// 	break
	// case Link:
	// 	i = i.(Link)
	// 	break
	// case Template:
	// 	i = i.(Template)
	// 	break
	// default:
	// 	return fmt.Errorf("??? %s", t)
	// }
	// if i == nil {
	// 	i = map[string]string{}
	// }
	// if v.IsValid() {
	// 	v.SetInterface(i)
	// }
	fmt.Println(v, i)
	// paths := f.Interface().(map[string]string)
	// if paths == nil {
	// 	paths = map[string]string{}
	// }
	// src, dst := dotfile.SplitPath(s)
	// paths[src] = dst
	return nil
}

// RegisterCopy ...
func (r *Role) RegisterCopy(s string) error {
	if r.Copy == nil {
		r.Copy = map[string]string{}
	}
	src, dst := dotfile.SplitPath(s)
	r.Copy[src] = dst
	return nil
}

// RegisterLink ...
func (r *Role) RegisterLink(s string) error {
	if r.Link == nil {
		r.Link = map[string]string{}
	}
	src, dst := dotfile.SplitPath(s)
	r.Link[src] = dst
	return nil
}

// RegisterTemplate ...
func (r *Role) RegisterTemplate(s string) error {
	if r.Template == nil {
		r.Template = map[string]string{}
	}
	src, dst := dotfile.SplitPath(s)
	r.Template[src] = dst
	return nil
}

// Init ...
func (r *Role) Init() error {
	if !exist(target) {
		return fmt.Errorf("Directory does not exist: %s", target)
	}
	if r.Path == "" {
		r.Path = filepath.Join(target, Options.RoleDir, r.Name)
	}
	// r.URL = ParseURL(r.URL)
	if Verbose > 0 {
		fmt.Printf("# [%s] Syncing %s %s\n", r.Name, r.Path, r.URL)
	}
	if err := r.Sync(); err != nil {
		return err
	}
	return nil
}

// Sync ...
func (r *Role) Sync() error {
	if r.URL == "" && !exist(r.Path) {
		return fmt.Errorf("# Role %s has no URL and could not be found in %s", r.Name, r.Path)
	}
	repo := NewRepo(r.Path, r.URL)
	exists := exist(repo.Path)
	// Clone if the local directory does not exist
	if !exist(repo.Path) {
		switch err := repo.Clone(); err {
		case nil:
			exists = true
			break
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
		break
	case ErrNoGitDir:
		if !exists {
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
		break
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
	cfgPath := filepath.Join(r.Path, name)
	if !exist(cfgPath) {
		fmt.Printf("No role config file found: %s\n", cfgPath)
		return "", nil
	}
	cfg, err := readConfig(cfgPath)
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
		return fmt.Errorf("Already enabled: %s", r.Name)
	}
	r.Enabled = true
	return nil
}

// Disable ...
func (r *Role) Disable() error {
	if r.Enabled == false {
		return fmt.Errorf("Already disabled: %s", r.Name)
	}
	r.Enabled = false
	return nil
}

// Prepare ...
func (r *Role) Prepare() error {
	if err := r.PreparePaths(&r.Copy); err != nil {
		return err
	}
	if err := r.PrepareLines(&r.Line); err != nil {
		return err
	}
	if err := r.PreparePaths(&r.Link); err != nil {
		return err
	}
	if err := r.PreparePaths(&r.Template); err != nil {
		return err
	}
	return nil
}

// PreparePaths ...
func (r *Role) PreparePaths(p *parsers.Paths) error {
	// in interface{} p := in.(*parsers.Paths)
	var paths parsers.Paths = make(map[string]string, len(*p))
	for src, dst := range *p {
		//fmt.Println("PREPARE", src, dst)
		// Prepend role directory to source path
		src = filepath.Join(r.Path, src)
		// Check frob globs
		if strings.Contains(src, "*") {
			//fmt.Println("*", src, dst)
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
				t, err := prepareTarget(s, dst)
				if err != nil {
					return err
				}
				paths[s] = t
			}
		} else {
			t, err := prepareTarget(src, dst)
			if err != nil {
				return err
			}
			paths[src] = t
		}
	}
	*p = paths
	return nil
}

func prepareTarget(src, dst string) (string, error) {
	//fmt.Println("+", src, dst)
	_, f := filepath.Split(src)
	if f == "" {
		return "", fmt.Errorf("Error (no source file name) while parsing: %s / %s", src, dst)
	}
	baseDir := filepath.Join(target, dst)
	// if _, err := dotfile.CreateDir(baseDir); err != nil {
	// 	return baseDir, err
	// }
	t := filepath.Join(baseDir, f)
	return t, nil
}

// PrepareLines ...
func (r *Role) PrepareLines(l *map[string]string) error {
	lines := make(map[string]string, len(*l))
	for file, line := range *l {
		// Prepend role directory to source path
		file = filepath.Join(target, file)
		lines[file] = line
	}
	*l = lines
	return nil
}

// GetField ...
func (r *Role) GetField(key string) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(r)).FieldByName(key)
}

// Do ...
func (r *Role) Do(a string, filter []string) error {
	fmt.Printf("# Role: %+v\n", r.Name)
	if len(filter) == 0 {
		filter = defaultTasks
	}
	a = strings.Title(a)
	v := r.GetField(a)
	if !v.IsValid() {
		return fmt.Errorf("Could not get field %s: %s / %s", a, v, a)
	}
	before := v.Interface().([]string)
	if len(before) > 0 && dotfile.HasOne([]string{"exec"}, filter) {
		for _, c := range before {
			task := &dotfile.ExecTask{
				Cmd: c,
			}
			if err := task.Do(a); err != nil {
				return err
			}
		}
	}
	if r.Env != nil {
		for k, v := range r.Env {
			k = strings.ToTitle(k)
			// Set role environment
			if err := dotfile.SetEnv(k, v); err != nil {
				// fmt.Fprintf(os.Stderr, err)
				return err
			}
		}
	}
	if r.Pkg != nil && Options.Packages && dotfile.HasOne([]string{"package"}, filter) {
		for _, v := range r.Pkg {
			if len(v.OS) > 0 && !dotfile.HasOSType(v.OS...) {
				continue
			}
			if v.Action != "" && strings.ToLower(v.Action) != strings.ToLower(a) {
				continue
			}
			task := &dotfile.PkgTask{
				Name: v.Name,
				Sudo: Options.Sudo,
			}
			if err := task.Do(a); err != nil {
				return err
			}
		}
	}
	if r.Copy != nil && dotfile.HasOne([]string{"copy"}, filter) {
		for s, t := range r.Copy {
			task := &dotfile.CopyTask{
				Source: s,
				Target: t,
			}
			if err := task.Do(a); err != nil {
				return err
			}
		}
	}
	if r.Line != nil && dotfile.HasOne([]string{"line"}, filter) {
		for s, t := range r.Line {
			task := &dotfile.LineTask{
				File: s,
				Line: t,
			}
			if err := task.Do(a); err != nil {
				return err
			}
		}
	}
	if r.Link != nil && dotfile.HasOne([]string{"link"}, filter) {
		for s, t := range r.Link {
			task := &dotfile.LinkTask{
				Source: s,
				Target: t,
			}
			if err := task.Do(a); err != nil {
				return err
			}
		}
	}
	if r.Template != nil && dotfile.HasOne([]string{"template"}, filter) {
		for s, t := range r.Template {
			task := &dotfile.TemplateTask{
				Source: s,
				Target: t,
				Env:    r.Env,
			}
			if err := task.Do(a); err != nil {
				return err
			}
		}
	}
	if r.Env != nil {
		// TODO: Restore original environment
	}
	after := r.GetField("Post" + a).Interface().([]string)
	if len(after) > 0 && dotfile.HasOne([]string{"exec"}, filter) {
		for _, c := range after {
			task := &dotfile.ExecTask{
				Cmd: c,
			}
			if err := task.Do(a); err != nil {
				return err
			}
		}
	}
	return nil
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
