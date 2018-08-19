package dot

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/LEI/dot/internal/git"
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
)

// TODO flag
var ignoredFilePatterns = []string{
	"*.json",
	"*.md",
	"*.yml",
	".git",
}

var taskListFields = []string{
	// "Pkgs",
	"Dirs",
	"Files",
	"Links",
	"Tpls",
	"Lines",
	// "Install",
	// "PostInstall",
	// "Remove",
	// "PostRemove",
}

// RoleConfig struct
type RoleConfig struct {
	Role *Role // `mapstructure:",squash"`
}

// Role struct
type Role struct {
	Name string
	Path string
	URL  string
	// Tasks []string

	OS  []string
	Env map[string]string
	// Vars  types.Map
	// IncludeVars types.IncludeMap

	Deps []string `mapstructure:"dependencies"`
	Pkgs []*Pkg   `mapstructure:"pkg"`

	Dirs  []*Dir      `mapstructure:"dir"`
	Files []*Copy     `mapstructure:"copy"`
	Links []*Link     `mapstructure:"link"`
	Tpls  []*Template `mapstructure:"template"`
	Lines []*Line     `mapstructure:"line"`

	Install     []*Hook
	PostInstall []*Hook `mapstructure:"post_install"`
	Remove      []*Hook
	PostRemove  []*Hook `mapstructure:"post_remove"`

	// Ignore []string
	// Target string

	// synced bool
	configFile string
}

// NewRole from a config file path
func NewRole(path string) (*Role, error) {
	rc := &RoleConfig{}
	data, err := ReadConfigFile(path)
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
	return rc.Role, err
}

func roleDecodeHook(f reflect.Type, t reflect.Type, i interface{}) (interface{}, error) {
	switch val := i.(type) {
	case string:
		switch t {
		case reflect.TypeOf((*Hook)(nil)):
			i = &Hook{Command: val}
		case reflect.TypeOf((*Dir)(nil)):
			i = &Dir{Path: val}
		case reflect.TypeOf((*Pkg)(nil)):
			i = &Pkg{Name: val}
		case reflect.TypeOf((*Link)(nil)):
			i = &Link{Source: val}
		case reflect.TypeOf((*Template)(nil)):
			i = &Template{Source: val}
		case reflect.TypeOf((*Line)(nil)):
			i = &Template{Source: val}
		}
	case map[interface{}]interface{}:
		switch t {
		case reflect.TypeOf(([]*Line)(nil)):
			i = decodeLines(val)
		}
	}
	return i, nil
}

// Transform map[string]interface{} to []*Line
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

func (r *Role) String() string {
	// return fmt.Sprintf("%s %s", r.Name, r.URL)
	s := ""
	s += fmt.Sprintf("%s\n", r.Name)
	s += fmt.Sprintf("  Path: %s\n", r.Path)
	s += fmt.Sprintf("  URL: %s\n", r.URL)
	if len(r.Deps) > 0 {
		s += fmt.Sprintf("  Deps: %s\n", r.Deps)
	}
	if len(r.Env) > 0 {
		s += fmt.Sprintf("  Env: %s\n", r.Env)
	}
	if len(r.OS) > 0 {
		s += fmt.Sprintf("  OS: %s\n", r.OS)
	}
	if len(r.Pkgs) > 0 {
		s += fmt.Sprintf("  Pkgs: %s\n", r.Pkgs)
	}
	s += formatRoleTasks("  ", r)
	if len(r.Install) > 0 {
		s += fmt.Sprintf("  Install: %s\n", r.Install)
	}
	if len(r.PostInstall) > 0 {
		s += fmt.Sprintf("  PostInstall: %s\n", r.PostInstall)
	}
	if len(r.Remove) > 0 {
		s += fmt.Sprintf("  Remove: %s\n", r.Remove)
	}
	if len(r.PostRemove) > 0 {
		s += fmt.Sprintf("  PostRemove: %s\n", r.PostRemove)
	}
	return strings.TrimRight(s, "\n")
}

func formatRoleTasks(prefix string, r *Role) string {
	s := ""
	if len(r.Dirs) > 0 {
		s += fmt.Sprintf("%sDirs:\n", prefix)
		s += formatTasks(prefix+"  ", r.Dirs)
	}
	if len(r.Files) > 0 {
		s += fmt.Sprintf("%sFiles:\n", prefix)
		s += formatTasks(prefix+"  ", r.Files)
	}
	if len(r.Links) > 0 {
		s += fmt.Sprintf("%sLinks:\n", prefix)
		s += formatTasks(prefix+"  ", r.Links)
	}
	if len(r.Tpls) > 0 {
		s += fmt.Sprintf("%sTemplates:\n", prefix)
		s += formatTasks(prefix+"  ", r.Tpls)
	}
	if r.Lines != nil {
		s += fmt.Sprintf("%sLines:\n", prefix)
		s += formatTasks(prefix+"  ", r.Lines)
	}
	return s
}

func formatTasks(prefix string, i interface{}) string {
	s := fmt.Sprintf("%s%s\n", prefix, i)
	// s := ""
	// tasks := i.([]Tasker)
	// for _, t := range tasks {
	// 	s += fmt.Sprintf("%s%s\n", prefix, t)
	// }
	return s
}

// Sync role repository
func (r *Role) Sync() error {
	repo, err := git.NewRepo(r.Path, r.URL)
	if err != nil {
		return err
	}
	if dirExists(r.Path) {
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

// LoadConfig file
func (r *Role) LoadConfig() error {
	if r.configFile == "" {
		return fmt.Errorf("role %s: empty config file path", r.Name)
	}
	role, err := NewRole(r.configFile)
	if err != nil {
		return fmt.Errorf("role %s: %s", r.Name, err)
	}
	// fmt.Printf("Merging role config %+v with original %+v\n", role, r)
	return mergo.Merge(r, role)
}

// Parse all role tasks
func (r *Role) Parse(target string) error {
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
	return nil
}

// ParseDirs tasks
func (r *Role) ParseDirs(target string) error {
	for _, d := range r.Dirs {
		d.Path = os.ExpandEnv(d.Path)
		if !filepath.IsAbs(d.Path) {
			d.Path = filepath.Join(target, d.Path)
		}
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
			links = append(links, &ll)
		}
	}
	r.Links = links
	return nil
}

// ParseTpls tasks
func (r *Role) ParseTpls(target string) error {
	templates := []*Template{}
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
		paths, err := preparePaths(target, t.Source, t.Target)
		if err != nil {
			return err
		}
		for k, v := range paths {
			tt := *t
			tt.Source = k
			tt.Target = v
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
		return src, dst, fmt.Errorf("unhandled path spec: %s", src)
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
	//fmt.Println("+", src, dst)
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

// StatusDirs ...
func (r *Role) StatusDirs() error {
	return checkTasks(r.taskDirs())
}

// StatusFiles ...
func (r *Role) StatusFiles() error {
	return checkTasks(r.taskFiles())
}

// StatusLinks ...
func (r *Role) StatusLinks() error {
	return checkTasks(r.taskLinks())
}

// StatusTpls ...
func (r *Role) StatusTpls() error {
	return checkTasks(r.taskTpls())
}

// StatusLines ...
func (r *Role) StatusLines() error {
	return checkTasks(r.taskLines())
}

func checkTasks(s []Tasker) error {
	ok := 0 // make([]bool, len(r.Tpls))
	for _, t := range s {
		err := t.Status()
		switch err {
		case nil:
		case ErrAlreadyExist:
			ok++ // [i] = true
		default:
			return err
		}
	}
	if ok == len(s) {
		return ErrAlreadyExist
	}
	return nil
}

// Ok returns true if already installed
func (r *Role) Ok() bool {
	err := r.Status()
	ok := IsOk(err)
	if err != nil && !ok {
		fmt.Fprintf(os.Stderr, "warning while checking %s role status: %s\n", r.Name, err)
	}
	return ok // err == nil || err == ErrAlreadyExist
}

// Status of role tasks
func (r *Role) Status() error {
	// err != nil || err != ErrAlreadyExist
	// if err != nil && !IsOk(err) {
	// 	return err
	// } else if err == nil {
	// 	return nil
	// }
	if err := r.StatusDirs(); !IsOk(err) {
		return err
	}
	if err := r.StatusFiles(); !IsOk(err) {
		return err
	}
	if err := r.StatusLinks(); !IsOk(err) {
		return err
	}
	if err := r.StatusTpls(); !IsOk(err) {
		return err
	}
	if err := r.StatusLines(); !IsOk(err) {
		return err
	}
	return ErrAlreadyExist
}

// taskDirs ...
func (r *Role) taskDirs() []Tasker {
	s := make([]Tasker, len(r.Lines))
	for i := range r.Lines {
		s[i] = r.Lines[i]
	}
	return s
}

// taskFiles ...
func (r *Role) taskFiles() []Tasker {
	s := make([]Tasker, len(r.Lines))
	for i := range r.Lines {
		s[i] = r.Lines[i]
	}
	return s
}

// taskLinks ...
func (r *Role) taskLinks() []Tasker {
	s := make([]Tasker, len(r.Lines))
	for i := range r.Lines {
		s[i] = r.Lines[i]
	}
	return s
}

// taskTpls ...
func (r *Role) taskTpls() []Tasker {
	s := make([]Tasker, len(r.Lines))
	for i := range r.Lines {
		s[i] = r.Lines[i]
	}
	return s
}

// taskLines ...
func (r *Role) taskLines() []Tasker {
	s := make([]Tasker, len(r.Lines))
	for i := range r.Lines {
		s[i] = r.Lines[i]
	}
	return s
}
