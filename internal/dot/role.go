package dot

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

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

// RoleConfig struct
type RoleConfig struct {
	Role Role // `mapstructure:",squash"`
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
	Pkgs Pkgs     `mapstructure:"pkg"`

	Dirs      Dirs      `mapstructure:"dir"`
	Files     Files     `mapstructure:"copy"`
	Links     Links     `mapstructure:"link"`
	Templates Templates `mapstructure:"template"`
	Lines     Lines     `mapstructure:"line"`

	Install     Hooks
	PostInstall Hooks `mapstructure:"post_install"`
	Remove      Hooks
	PostRemove  Hooks `mapstructure:"post_remove"`

	// Ignore []string
	// Target string

	// synced bool
}

func (r *Role) String() string {
	// return fmt.Sprintf("%s %s", r.Name, r.URL)
	s := ""
	s += fmt.Sprintf("%s\n", r.Name)
	s += fmt.Sprintf("  Path: %s\n", r.Path)
	s += fmt.Sprintf("  URL: %s\n", r.URL)

	if r.Deps != nil {
		s += fmt.Sprintf("  Deps: %s\n", r.Deps)
	}
	if r.Env != nil {
		s += fmt.Sprintf("  Env: %s\n", r.Env)
	}
	if r.OS != nil {
		s += fmt.Sprintf("  OS: %s\n", r.OS)
	}
	if r.Pkgs != nil {
		s += fmt.Sprintf("  Pkgs: %s\n", r.Pkgs)
	}
	if t := tasksPrefix("  ", r); t != "" {
		s += t
	}
	if r.Install != nil {
		s += fmt.Sprintf("  Install: %s\n", r.Install)
	}
	if r.PostInstall != nil {
		s += fmt.Sprintf("  PostInstall: %s\n", r.PostInstall)
	}
	if r.Remove != nil {
		s += fmt.Sprintf("  Remove: %s\n", r.Remove)
	}
	if r.PostRemove != nil {
		s += fmt.Sprintf("  PostRemove: %s\n", r.PostRemove)
	}
	return strings.TrimRight(s, "\n")
}

func tasksPrefix(prefix string, r *Role) string {
	s := ""
	if r.Dirs != nil {
		s += fmt.Sprintf("%sDirs: %s\n", prefix, r.Dirs)
	}
	if r.Files != nil {
		s += fmt.Sprintf("%sFiles: %s\n", prefix, r.Files)
	}
	if r.Lines != nil {
		s += fmt.Sprintf("%sLines: %s\n", prefix, r.Lines)
	}
	if r.Links != nil {
		s += fmt.Sprintf("%sLinks: %s\n", prefix, r.Links)
	}
	if r.Templates != nil {
		s += fmt.Sprintf("%sTemplates: %s\n", prefix, r.Templates)
	}
	return s
}

// LoadConfig ...
func (r *Role) LoadConfig() error {
	cfgPath := filepath.Join(r.Path, ".dot.yml")
	role, err := LoadRole(cfgPath)
	if err != nil {
		return fmt.Errorf("%s: %s", r.Name, err)
	}
	// fmt.Printf("MERGE %+v\n", role.Env)
	return mergo.Merge(r, role)
}

// Parse all role tasks
func (r *Role) Parse(target string) error {
	if r.Path == "" {
		r.Path = filepath.Join(os.ExpandEnv("$HOME"), ".dot", r.Name)
	}
	// fmt.Println("prepare", r.Name)
	if err := r.ParseDirs(target); err != nil {
		return err
	}
	if err := r.ParseFiles(target); err != nil {
		return err
	}
	if err := r.ParseLinks(target); err != nil {
		return err
	}
	if err := r.ParseTemplates(target); err != nil {
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
		if err := d.Prepare(target); err != nil {
			return err
		}
	}
	return nil
}

// ParseFiles tasks
func (r *Role) ParseFiles(target string) error {
	files := Files{}
	for _, c := range r.Files {
		c.Source = os.ExpandEnv(c.Source)
		c.Target = os.ExpandEnv(c.Target)
		if err := c.Prepare(target); err != nil {
			return err
		}
		paths, err := preparePaths(target, c.Source, c.Target)
		if err != nil {
			return err
		}
		for k, v := range paths {
			// cc := c
			// cc.Source = k
			// cc.Target = v
			c.Source = k
			c.Target = v
			files = append(files, c)
		}
	}
	r.Files = files
	return nil
}

// ParseLinks tasks
func (r *Role) ParseLinks(target string) error {
	links := Links{}
	for _, l := range r.Links {
		l.Source = os.ExpandEnv(l.Source)
		l.Target = os.ExpandEnv(l.Target)
		if err := l.Prepare(target); err != nil {
			return err
		}
		paths, err := preparePaths(target, l.Source, l.Target)
		if err != nil {
			return err
		}
		for k, v := range paths {
			l.Source = k
			l.Target = v
			links = append(links, l)
		}
	}
	r.Links = links
	return nil
}

// ParseTemplates tasks
func (r *Role) ParseTemplates(target string) error {
	templates := Templates{}
	for _, t := range r.Templates {
		t.Source = os.ExpandEnv(t.Source)
		t.Target = os.ExpandEnv(t.Target)
		if err := t.Prepare(target); err != nil {
			return err
		}
		paths, err := preparePaths(target, t.Source, t.Target)
		if err != nil {
			return err
		}
		for k, v := range paths {
			t.Source = k
			t.Target = v
			templates = append(templates, t)
		}
	}
	r.Templates = templates
	return nil
}

// ParseLines tasks
func (r *Role) ParseLines(target string) error {
	for _, l := range r.Lines {
		l.Target = os.ExpandEnv(l.Target)
		if err := l.Prepare(target); err != nil {
			return err
		}
	}
	return nil
}

func preparePaths(target, src, dst string) (map[string]string, error) {
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

// Check magix chars recognized by Match
func hasMeta(path string) bool {
	magicChars := `*?[`
	if runtime.GOOS == "windows" {
		magicChars = `*?[\`
	}
	return strings.ContainsAny(path, magicChars)
}

// NewRole ...
func NewRole() *Role {
	return &Role{}
}

// LoadRole ...
func LoadRole(path string) (Role, error) {
	rc := &RoleConfig{}
	data, err := Read(path)
	if err != nil {
		return rc.Role, err
	}
	decoderConfig := &mapstructure.DecoderConfig{
		// DecodeHook:       weaklyTypedHook,
		DecodeHook:       roleDecodeHook,
		ErrorUnused:      true,
		WeaklyTypedInput: true,
		Result:           &rc,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return rc.Role, err
	}
	if err := decoder.Decode(data); err != nil {
		return rc.Role, err
	}
	// if err := mapstructure.WeakDecode(data, &rc); err != nil {
	// 	return rc.Role, err
	// }
	return rc.Role, nil
}

func roleDecodeHook(f reflect.Type, t reflect.Type, i interface{}) (interface{}, error) {
	// input := i.(map[string]interface{})
	// f == reflect.TypeOf("")
	// fmt.Printf("\n---\nroleDecodeHook %s -> %s\n%+v\n---\n", f, t, i)
	switch val := i.(type) {
	case string:
		switch t {
		case reflect.TypeOf((*Hook)(nil)):
			i = &Hook{Command: val}
		// case reflect.TypeOf((*Dirs)(nil)):
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
		default:
			// fmt.Println("roleDecodeHook string", t, "/", val)
			// fmt.Println("FALLBACK2", t)
		}
	case map[interface{}]interface{}:
		// case map[string]interface{}:
		switch t {
		case reflect.TypeOf((Lines)(nil)):
			lines := Lines{}
			for k, v := range val {
				lines = append(lines, &Line{
					Target: k.(string),
					Data:   v.(string),
				})
			}
			i = lines
		}
	// case Line:
	// 	fmt.Println("LIIIIIIIIIIIINE")
	// case *Line:
	// 	fmt.Println("*LIIIIIIIIIIIINE")
	// case Lines:
	// 	fmt.Println("Lines", reflect.TypeOf(i), f, "======>", t, reflect.TypeOf(val))
	// 	// i = map[string]string{} // val
	// case *Lines:
	// 	fmt.Println("*Lines", reflect.TypeOf(i), f, "======>", t, reflect.TypeOf(val))
	default:
	}
	return i, nil
}

// func weaklyTypedHook(
// 	f reflect.Kind,
// 	t reflect.Kind,
// 	data interface{}) (interface{}, error) {
// 	dataVal := reflect.ValueOf(data)
// 	switch t {
// 	case reflect.String:
// 		switch f {
// 		case reflect.Bool:
// 			if dataVal.Bool() {
// 				return "1", nil
// 			}
// 			return "0", nil
// 		case reflect.Float32:
// 			return strconv.FormatFloat(dataVal.Float(), 'f', -1, 64), nil
// 		case reflect.Int:
// 			return strconv.FormatInt(dataVal.Int(), 10), nil
// 		case reflect.Slice:
// 			dataType := dataVal.Type()
// 			elemKind := dataType.Elem().Kind()
// 			if elemKind == reflect.Uint8 {
// 				return string(dataVal.Interface().([]uint8)), nil
// 			}
// 		case reflect.Uint:
// 			return strconv.FormatUint(dataVal.Uint(), 10), nil
// 		}
// 	}

// 	return data, nil
// }
