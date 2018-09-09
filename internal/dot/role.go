package dot

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/LEI/dot/internal/conf"
	"github.com/LEI/dot/internal/git"
	"github.com/LEI/dot/internal/shell"
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
)

// TODO flag
var ignoredFilePatterns = []string{
	// "*.json",
	"*.md",
	"*.toml",
	"*.yml",
	"*.DS_Store",
	".git",
}

// RoleConfig struct
type RoleConfig struct {
	Role *Role // `mapstructure:",squash"`
}

// Role struct
type Role struct {
	Name string
	Path string
	URL  string   // Git repository URL
	Git  *url.URL // repo host and user
	// Scheme string // Remote type: git (default), https, ssh

	// Tasks []string

	// Ignore []string
	// Target string

	Disabled bool `mapstructure:",omitempty"`

	OS          []string
	Env         *Env
	Vars        *Vars
	IncludeVars []string

	Deps []string `mapstructure:"depends"`
	Pkgs []*Pkg   `mapstructure:"pkg"`

	Dirs   []*Dir   `mapstructure:"dir"`
	Files  []*Copy  `mapstructure:"copy"`
	Links  []*Link  `mapstructure:"link"`
	Tpls   []*Tpl   `mapstructure:"template"`
	Lines  []*Line  `mapstructure:"line"`
	Blocks []*Block `mapstructure:"block"`

	Install     []*Hook
	PostInstall []*Hook `mapstructure:"post_install"`
	Remove      []*Hook
	PostRemove  []*Hook `mapstructure:"post_remove"`

	// loadedRole bool
	configFile string
}

func (r *Role) String() (s string) {
	// return fmt.Sprintf("%s %s", r.Name, r.URL)
	prefix := "  "
	s += fmt.Sprintf("%s\n", r.Name)
	s += fmt.Sprintf("%sPath: %s\n", prefix, r.Path)
	s += fmt.Sprintf("%sURL: %s\n", prefix, r.URL)
	if len(r.OS) > 0 {
		s += fmt.Sprintf("%sOS: %s\n", prefix, r.OS)
	}
	if len(*r.Env) > 0 {
		s += fmt.Sprintf("%sEnv: %s\n", prefix, *r.Env)
	}
	if len(r.IncludeVars) > 0 {
		s += fmt.Sprintf("%sIncludeVars: %s\n", prefix, r.IncludeVars)
	}
	if len(*r.Vars) > 0 {
		s += fmt.Sprintf("%sVars: %s\n", prefix, *r.Vars)
	}
	if len(r.Deps) > 0 {
		s += fmt.Sprintf("%sDeps: %s\n", prefix, r.Deps)
	}
	if len(r.Pkgs) > 0 {
		s += fmt.Sprintf("%sPkgs: %s\n", prefix, r.Pkgs)
	}
	s += formatRole(prefix, r)
	return strings.TrimRight(s, "\n")
}

func formatRole(prefix string, r *Role) (s string) {
	if len(r.Dirs) > 0 {
		s += fmt.Sprintf("%sDirs:\n", prefix)
		// s += formatRoleTasks(prefix+prefix, r.Dirs)
		for _, d := range r.Dirs {
			s += formatTask(prefix+prefix, d)
		}
	}
	if len(r.Files) > 0 {
		s += fmt.Sprintf("%sFiles:\n", prefix)
		// s += formatRoleTasks(prefix+prefix, r.Files)
		for _, f := range r.Files {
			s += formatTask(prefix+prefix, f)
		}
	}
	if len(r.Links) > 0 {
		s += fmt.Sprintf("%sLinks:\n", prefix)
		// s += formatRoleTasks(prefix+prefix, r.Links)
		for _, l := range r.Links {
			s += formatTask(prefix+prefix, l)
		}
	}
	if len(r.Tpls) > 0 {
		s += fmt.Sprintf("%sTemplates:\n", prefix)
		// s += formatRoleTasks(prefix+prefix, r.Tpls)
		for _, t := range r.Tpls {
			s += formatTask(prefix+prefix, t)
		}
	}
	if len(r.Lines) > 0 { // r.Lines != nil
		s += fmt.Sprintf("%sLines:\n", prefix)
		// s += formatRoleTasks(prefix+prefix, r.Lines)
		for _, l := range r.Lines {
			s += formatTask(prefix+prefix, l)
		}
	}
	if len(r.Blocks) > 0 { // r.Blocks != nil
		s += fmt.Sprintf("%sBlocks:\n", prefix)
		// s += formatRoleTasks(prefix+prefix, r.Blocks)
		for _, b := range r.Blocks {
			s += formatTask(prefix+prefix, b)
		}
	}
	if len(r.Install) > 0 {
		s += fmt.Sprintf("%sInstall:\n", prefix)
		s += formatRoleHooks(prefix+prefix, r.Install)
	}
	if len(r.PostInstall) > 0 {
		s += fmt.Sprintf("%sPostInstall:\n", prefix)
		s += formatRoleHooks(prefix+prefix, r.PostInstall)
	}
	if len(r.Remove) > 0 {
		s += fmt.Sprintf("%sRemove:\n", prefix)
		s += formatRoleHooks(prefix+prefix, r.Remove)
	}
	if len(r.PostRemove) > 0 {
		s += fmt.Sprintf("%sPostRemove:\n", prefix)
		s += formatRoleHooks(prefix+prefix, r.PostRemove)
	}
	return s
}

func formatRoleHooks(prefix string, hooks []*Hook) (s string) {
	for _, h := range hooks {
		s += fmt.Sprintf("%s→ %s\n", prefix, h.String())
	}
	return s
}

// func formatRoleTasks(prefix string, is []interface{}) (s string) {
// 	for _, i := range is {
// 		s += formatTask(prefix, i.(Tasker))
// 	}
// 	return s
// }

// func formatTask(prefix string, t Tasker) (s string) {
// 	return fmt.Sprintf("%s→ %s\n", prefix, t.String())
// }

func formatTask(prefix string, i interface{}) (s string) {
	return fmt.Sprintf("%s→ %s\n", prefix, i)
}

// func formatTasks(prefix string, i interface{}) string {
// 	s := fmt.Sprintf("%s%s\n", prefix, i)
// 	return s
// }

// ShouldRun check
func (r *Role) ShouldRun() bool {
	return !r.Disabled
}

// Sync role repository
// TODO update role URL?
func (r *Role) Sync() (string, error) {
	// u, err := url.Parse(r.URL)
	// if err != nil {
	// 	return err
	// }
	out := ""
	repo, err := git.NewRepo(r.Git, r.URL, r.Path)
	if err != nil {
		return out, err
	}
	if dirExists(r.Path) {
		// fmt.Fprintf(dotCli.Out(), "Checking %s...\n", name)
		status, err := repo.Status()
		out += status
		if err != nil {
			return out, err
		}
		pull, err := repo.Pull()
		out += pull
		if err != nil {
			return out, err
		}
	} else {
		// fmt.Fprintf(dotCli.Out(), "Cloning %s into %s...\n", name, repo.Dir)
		clone, err := repo.Clone()
		out += clone
		if err != nil {
			return out, err
		}
	}
	return out, nil
}

// GetConfigFile path
func (r *Role) GetConfigFile() string {
	return r.configFile
}

// SetConfigFile path
func (r *Role) SetConfigFile(name string) *Role {
	if !filepath.IsAbs(name) {
		name = filepath.Join(r.Path, name)
	}
	r.configFile = name
	return r
}

// Load role config
func (r *Role) Load() error {
	if r.configFile == "" {
		return fmt.Errorf("role %s: empty config file path", r.Name)
	}
	role, err := roleLoadConfig(r.configFile)
	if err != nil {
		return fmt.Errorf("role %s: %s", r.Name, err)
	}
	// fmt.Printf("Using role config file: %s\n", r.configFile)
	// fmt.Printf("Merging role config:\n%+v\nwith original struct:\n%+v\n", role, r)
	if err = mergo.Merge(r, role); err != nil {
		return err
	}
	if r.Name == "" {
		return fmt.Errorf("empty role name: %+v", r)
	}
	return nil
}

// roleLoadConfig from a file path
func roleLoadConfig(path string) (*Role, error) {
	rc := &RoleConfig{}
	data, err := conf.ReadFile(path)
	if err != nil {
		return rc.Role, err
	}
	dc := &mapstructure.DecoderConfig{
		DecodeHook:       roleDecodeHook,
		ErrorUnused:      DecodeErrorUnused,
		WeaklyTypedInput: DecodeWeaklyTypedInput,
		Result:           &rc,
	}
	decoder, err := mapstructure.NewDecoder(dc)
	if err != nil {
		return rc.Role, err
	}
	err = decoder.Decode(data)
	if err != nil {
		return rc.Role, err
	}
	return rc.Role, nil
}

// https://github.com/ernesto-jimenez/gogen/tree/master/cmd/gounmarshalmap
func roleDecodeHook(f reflect.Type, t reflect.Type, i interface{}) (interface{}, error) {
	switch val := i.(type) {
	case string:
		switch t {
		case reflect.TypeOf((*Env)(nil)),
			reflect.TypeOf((Env)(nil)):
			i = NewEnv(val)
		case reflect.TypeOf((*Hook)(nil)):
			i = NewHook(val)
		case reflect.TypeOf((*Pkg)(nil)):
			i = NewPkg(val) // &Pkg{Name: []string{val}}
		case reflect.TypeOf((*Dir)(nil)):
			i = NewDir(val) // Path
		case reflect.TypeOf((*Link)(nil)):
			i = NewLink(val)
		case reflect.TypeOf((*Tpl)(nil)):
			i = NewTpl(val)
		// case reflect.TypeOf((*Line)(nil)):
		// 	i = NewLine(val)
		// case reflect.TypeOf((*Block)(nil)):
		// 	i = NewBlock(val)
		default:
			// 	fmt.Println("FROM", f)
			// 	fmt.Printf("%+v (%T)\n", val, val)
			// 	fmt.Println("TO", t)
		}
	case []string:
		switch t {
		case reflect.TypeOf((*Env)(nil)):
			i = NewEnv(val)
		}
	case map[interface{}]interface{}:
		switch t {
		case reflect.TypeOf(([]*Line)(nil)):
			i = decodeLines(val)
		case reflect.TypeOf(([]*Block)(nil)):
			i = decodeBlocks(val)
		}
	}
	// default:
	// 	fmt.Printf("%+v (%T)\n", val, val)
	return i, nil
}

// Transform map[i{}]i{} to []*Line
func decodeLines(in map[interface{}]interface{}) []*Line {
	lines := []*Line{}
	for k, v := range in {
		lines = append(lines, &Line{
			Target: k.(string),
			Data:   v.(string),
		})
	}
	return lines
}

// Transform map[i{}]i{} to []*Block
func decodeBlocks(in map[interface{}]interface{}) []*Block {
	blocks := []*Block{}
	for k, v := range in {
		blocks = append(blocks, &Block{
			Target: k.(string),
			Data:   v.(string),
		})
	}
	return blocks
}

// Parse all role tasks
func (r *Role) Parse(target string) error {
	if err := r.ParseEnv(); err != nil {
		return err
	}
	if err := r.ParseVars(); err != nil {
		return err
	}
	// if err := r.ParsePkgs(target); err != nil {
	// 	return err
	// }
	if err := r.ParseDirs(target); err != nil {
		return err
	}
	if err := r.ParseFiles(target); err != nil {
		return err
	}
	if err := r.ParseLinks(target); err != nil {
		return err
	}
	if err := r.ParseTpls(target); err != nil {
		return err
	}
	if err := r.ParseLines(target); err != nil {
		return err
	}
	if err := r.ParseBlocks(target); err != nil {
		return err
	}
	if err := r.ParseHooks(target); err != nil {
		return err
	}
	return nil
}

// ParseEnv role
func (r *Role) ParseEnv() error {
	if r.Env == nil {
		r.Env = &Env{}
	}
	e, err := parseEnviron(r.Env)
	if err != nil {
		return err
	}
	r.Env = e
	return nil
}

// ParseVars role
func (r *Role) ParseVars() error {
	if r.Vars == nil {
		r.Vars = &Vars{}
	}
	vars, err := parseVars(r.Env, r.Vars, r.IncludeVars...)
	if err != nil {
		return err
	}
	r.Vars = vars
	return nil
}

// ParseHooks tasks
func (r *Role) ParseHooks(target string) error {
	for _, h := range r.Install {
		if err := parseRoleHook(r.Env, h); err != nil {
			return fmt.Errorf("%s: %s", r.Name+" install hook", err)
		}
	}
	for _, h := range r.PostInstall {
		if err := parseRoleHook(r.Env, h); err != nil {
			return fmt.Errorf("%s: %s", r.Name+" post_install hook", err)
		}
	}
	for _, h := range r.Remove {
		if err := parseRoleHook(r.Env, h); err != nil {
			return fmt.Errorf("%s: %s", r.Name+" remove hook", err)
		}
	}
	for _, h := range r.PostRemove {
		if err := parseRoleHook(r.Env, h); err != nil {
			return fmt.Errorf("%s: %s", r.Name+" post_remove hook", err)
		}
	}
	return nil
}

// Hook environment variables are not expanded now to allow
// command substitution to be done at runtime
func parseRoleHook(e *Env, h *Hook) error {
	if h == nil || h.Command == "" && (h.URL == "" || h.Dest == "") {
		return fmt.Errorf("empty command")
	}
	if h.Env == nil {
		h.Env = &Env{}
	}
	// Merge given environment (global role config)
	for k, v := range *e {
		if _, ok := (*h.Env)[k]; !ok {
			(*h.Env)[k] = v
		}
	}
	// if h.Command != "" {
	// 	h.Command = os.Expand(h.Command, func(s string) string {
	// 		return (*h.Env)[s]
	// 	})
	// }
	if h.URL != "" {
		h.URL = os.Expand(h.URL, func(s string) string {
			return (*h.Env)[s]
		})
	}
	if h.Dest != "" {
		h.Dest = os.Expand(h.Dest, func(s string) string {
			return (*h.Env)[s]
		})
	}
	return nil
}

// ParseDirs tasks
func (r *Role) ParseDirs(target string) error {
	for _, d := range r.Dirs {
		d.Path = os.ExpandEnv(d.Path)
		if !filepath.IsAbs(d.Path) {
			d.Path = filepath.Join(target, d.Path)
		}
		// if err := d.Prepare(); err != nil {
		// 	return err
		// }
	}
	return nil
}

// ParseFiles tasks
func (r *Role) ParseFiles(target string) error {
	files := []*Copy{}
	for _, c := range r.Files {
		c.Source = os.ExpandEnv(c.Source)
		c.Target = os.ExpandEnv(c.Target)
		if c.Target == "" {
			src, dst, err := parsePaths(c.Source)
			if err != nil {
				return err
			}
			c.Source = src
			c.Target = dst
		}
		if !filepath.IsAbs(c.Source) {
			c.Source = filepath.Join(r.Path, c.Source)
		}
		if !filepath.IsAbs(c.Target) {
			c.Target = filepath.Join(target, c.Target)
		}
		paths, err := preparePaths(target, c.Source, c.Target)
		if err != nil {
			return err
		}
		for k, v := range paths {
			cc := *c
			cc.Source = k
			cc.Target = v
			// if err := cc.Prepare(); err != nil {
			// 	return err
			// }
			files = append(files, &cc)
		}
	}
	r.Files = files
	return nil
}

// ParseLinks tasks
func (r *Role) ParseLinks(target string) error {
	links := []*Link{}
	for _, l := range r.Links {
		l.Source = os.ExpandEnv(l.Source)
		l.Target = os.ExpandEnv(l.Target)
		if l.Target == "" {
			src, dst, err := parsePaths(l.Source)
			if err != nil {
				return err
			}
			l.Source = src
			l.Target = dst
		}
		if !filepath.IsAbs(l.Source) {
			l.Source = filepath.Join(r.Path, l.Source)
		}
		if !filepath.IsAbs(l.Target) {
			l.Target = filepath.Join(target, l.Target)
		}
		paths, err := preparePaths(target, l.Source, l.Target)
		if err != nil {
			return err
		}
		for k, v := range paths {
			ll := *l
			ll.Source = k
			ll.Target = v
			// if err := ll.Prepare(); err != nil {
			// 	return err
			// }
			links = append(links, &ll)
		}
	}
	r.Links = links
	return nil
}

// ParseTpls tasks
func (r *Role) ParseTpls(target string) error {
	templates := []*Tpl{}
	for _, t := range r.Tpls {
		t.Source = os.ExpandEnv(t.Source)
		t.Target = os.ExpandEnv(t.Target)
		if t.Target == "" {
			src, dst, err := parsePaths(t.Source)
			if err != nil {
				return err
			}
			t.Source = src
			t.Target = dst
		}
		if !filepath.IsAbs(t.Source) {
			t.Source = filepath.Join(r.Path, t.Source)
		}
		if !filepath.IsAbs(t.Target) {
			t.Target = filepath.Join(target, t.Target)
		}
		// Merge task env with role env
		if t.Env == nil {
			t.Env = &Env{}
		}
		// for k, v := range r.Env {
		// 	if _, ok := t.Env[k]; !ok {
		// 		t.Env[k] = v
		// 	}
		// }
		// Merge task vars with role vars
		if t.Vars == nil {
			t.Vars = &Vars{}
		}
		// for k, v := range r.Vars {
		// 	_, ok := t.Vars[k]
		// 	if !ok {
		// 		t.Vars[k] = v
		// 	}
		// }
		// Glob templates
		paths, err := preparePaths(target, t.Source, t.Target)
		if err != nil {
			return err
		}
		for k, v := range paths {
			tt := *t
			tt.Source = k
			tt.Target = v
			if err := tt.Prepare(); err != nil {
				return err
			}
			for k, v := range *r.Env {
				if _, ok := (*tt.Env)[k]; !ok {
					(*tt.Env)[k] = v
				}
			}
			for k, v := range *r.Vars {
				_, ok := (*tt.Vars)[k]
				if !ok {
					(*tt.Vars)[k] = v
				}
			}
			templates = append(templates, &tt)
		}
	}
	r.Tpls = templates
	return nil
}

// ParseLines tasks
func (r *Role) ParseLines(target string) error {
	for _, l := range r.Lines {
		l.Target = os.ExpandEnv(l.Target)
		if !filepath.IsAbs(l.Target) {
			l.Target = filepath.Join(target, l.Target)
		}
		// if err := l.Prepare(); err != nil {
		// 	return err
		// }
	}
	return nil
}

// ParseBlocks tasks
func (r *Role) ParseBlocks(target string) error {
	for _, b := range r.Blocks {
		b.Target = os.ExpandEnv(b.Target)
		if !filepath.IsAbs(b.Target) {
			b.Target = filepath.Join(target, b.Target)
		}
		// if err := b.Prepare(); err != nil {
		// 	return err
		// }
	}
	return nil
}

// Parse src:dst paths
func parsePaths(p string) (src, dst string, err error) {
	parts := filepath.SplitList(p)
	switch len(parts) {
	case 1:
		src = p
	case 2:
		src = parts[0]
		dst = parts[1]
	default:
		// unhandled path spec
		return src, dst, fmt.Errorf(
			"too many paths (%d): %s",
			len(parts),
			src,
		)
	}
	return src, dst, nil
	// src = s
	// if strings.Contains(s, ":") {
	// 	parts := strings.Split(s, ":")
	// 	if len(parts) != 2 {
	// 		return src, dst, fmt.Errorf("unable to parse dest spec: %s", s)
	// 	}
	// 	src = parts[0]
	// 	dst = parts[1]
	// }
	// // if dst == "" && isDir(src) {
	// // 	dst = PathHead(src)
	// // }
	// return src, dst, nil
}

func preparePaths(target, src, dst string) (map[string]string, error) {
	ret := map[string]string{}
	//*links = append(*links, l)
	if hasMeta(src) { // strings.Contains(src, "*")
		// fmt.Println("*** GLOB", src, dst)
		glob, err := filepath.Glob(src)
		if err != nil {
			return ret, err
		}
	GLOB:
		for _, s := range glob {
			// Extract source file name
			_, n := filepath.Split(s)
			for _, i := range ignoredFilePatterns {
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
	// fmt.Printf("+ %q %q\n", src, dst)
	_, name := filepath.Split(src)
	if name == "" {
		return "", fmt.Errorf("no source file name for src / dst: %s / %s", src, dst)
	}
	if !filepath.IsAbs(dst) {
		dst = filepath.Join(target, dst)
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

// Check magic chars recognized by Match
func hasMeta(path string) bool {
	magicChars := `*?[`
	if runtime.GOOS == "windows" {
		magicChars = `*?[\`
	}
	return strings.ContainsAny(path, magicChars)
}

// Ok returns true if the role already installed.
func (r *Role) Ok() bool {
	err := r.Status()
	exists := IsExist(err)
	if err != nil && !exists {
		fmt.Fprintf(os.Stderr, "%s: %s\n", r.Name, err)
	}
	return exists // err == nil || err == ErrExist
}

// Status reports the state of all tasks of the role,
// failing at the first encountered.
func (r *Role) Status() error {
	// err != nil || err != ErrExist
	// if err != nil && !IsExist(err) {
	// 	return err
	// } else if err == nil {
	// 	return nil
	// }
	if err := r.StatusHooks(); err != nil { // !IsExist(err)
		return err
	}
	/* // Skip packages check to speed up the listing
	if err := r.StatusPkgs(); !IsExist(err) {
		return err
	} */
	if err := r.StatusDirs(); !IsExist(err) {
		return err
	}
	if err := r.StatusFiles(); !IsExist(err) {
		return err
	}
	if err := r.StatusLinks(); !IsExist(err) {
		return err
	}
	if err := r.StatusTpls(); !IsExist(err) {
		return err
	}
	if err := r.StatusLines(); !IsExist(err) {
		return err
	}
	if err := r.StatusBlocks(); !IsExist(err) {
		return err
	}
	return ErrExist
}

// StatusPkgs ...
func (r *Role) StatusPkgs() error {
	c := 0
	for _, t := range r.Pkgs {
		if err := checkTask(t); err != nil {
			return err
		}
		c++
	}
	if c == len(r.Pkgs) {
		return ErrExist
	}
	return nil
}

// StatusDirs ...
func (r *Role) StatusDirs() error {
	c := 0
	for _, t := range r.Dirs {
		if err := checkTask(t); err != nil {
			return err
		}
		c++
	}
	if c == len(r.Dirs) {
		return ErrExist
	}
	return nil
}

// StatusFiles ...
func (r *Role) StatusFiles() error {
	c := 0
	for _, t := range r.Files {
		if err := checkTask(t); err != nil {
			return err
		}
		c++
	}
	if c == len(r.Files) {
		return ErrExist
	}
	return nil
}

// StatusLinks ...
func (r *Role) StatusLinks() error {
	c := 0
	for _, t := range r.Links {
		if err := checkTask(t); err != nil {
			return err
		}
		c++
	}
	if c == len(r.Links) {
		return ErrExist
	}
	return nil
}

// StatusTpls ...
func (r *Role) StatusTpls() error {
	c := 0
	for _, t := range r.Tpls {
		if err := checkTask(t); err != nil {
			return err
		}
		c++
	}
	if c == len(r.Tpls) {
		return ErrExist
	}
	return nil
}

// StatusLines ...
func (r *Role) StatusLines() error {
	c := 0
	for _, t := range r.Lines {
		if err := checkTask(t); err != nil {
			return err
		}
		c++
	}
	if c == len(r.Lines) {
		return ErrExist
	}
	return nil
}

// StatusBlocks ...
func (r *Role) StatusBlocks() error {
	c := 0
	for _, t := range r.Blocks {
		if err := checkTask(t); err != nil {
			return err
		}
		c++
	}
	if c == len(r.Blocks) {
		return ErrExist
	}
	return nil
}

// StatusHooks ...
func (r *Role) StatusHooks() error {
	for _, t := range r.Install {
		if err := checkTask(t); err != nil {
			return err
		}
	}
	for _, t := range r.PostInstall {
		if err := checkTask(t); err != nil {
			return err
		}
	}
	for _, t := range r.Remove {
		if err := checkTask(t); err != nil {
			return err
		}
	}
	for _, t := range r.PostRemove {
		if err := checkTask(t); err != nil {
			return err
		}
	}
	return nil
}

// Check if a task is already executed
func checkTask(t Tasker) error {
	if err := t.Status(); err != nil {
		if !IsExist(err) {
			return err
		}
	}
	// terr, ok := err.(*OpError)
	// if ok {
	// 	err = terr.Err
	// }
	// switch err {
	// case nil:
	// case ErrExist:
	// 	c++ // [i] = true
	// default:
	// 	return err
	// }
	return nil
}

// Init role before install or remove
func (r *Role) Init() error {
	for k, v := range *r.Env {
		// fmt.Printf("$ export %s=%q\n", k, v)
		fmt.Printf("$ %s=%s\n", k, shell.FormatArgs([]string{v}))
	}
	return nil
}
