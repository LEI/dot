package role

import (
	"bytes"
	"fmt"
	// "github.com/LEI/dot/config"
	"github.com/LEI/dot/git"
	"github.com/spf13/viper"
	"os"
	"path"
	"strings"
	"text/template"
)

var (
	PathSep = string(os.PathSeparator)
)

type Role struct {
	Name, Origin   string
	Source, Target string
	Os             []string
	Package        *Package // `mapstructure:",squash"`
	Config         *viper.Viper
	Repo           *git.Repository
}

func (r *Role) Title() string {
	return strings.Title(r.Name)
}

// func (r *Role) String() string {
// 	return fmt.Sprintf("%s (%s) [%s -> %s] Dir: %s, Dirs: %v, Link: %s, Links: %v, Lines: %+v",
// 		r.Name, r.Origin, r.Source, r.Target, r.Dir, r.Dirs, r.Link, r.Links, r.Lines)
// }

// func (r *Role) Origin() string {
// 	origin := r.Source
// 	return fmt.Sprintf("%s", origin)
// }

func (r *Role) Init(source, target string) (*Role, error) {
	// r = &Role{Package: &Package{}}
	switch {
	case r.Name == "" && r.Origin != "":
		r.Name = r.Origin
		if strings.Contains(r.Name, PathSep) {
			r.Name = path.Base(r.Name)
		}
	case r.Name != "" && r.Origin == "":
		r.Origin = r.Name
	}
	if r.Name == "" || r.Origin == "" {
		return r, fmt.Errorf("Invalid role: %+v\n", r)
	}
	if r.Source == "" {
		r.Source = source
	}
	if r.Target == "" {
		r.Target = target
	}
	return r, nil
}

// name=user/repo
// user/repo
func ParseOrigin(str string) (name string, dir string, url string, err error) {
	var nameSep = "="
	name = str
	if strings.HasPrefix(str, PathSep) {
		dir = str
		name = path.Dir(dir)
		fi, err := os.Stat(str)
		if err != nil && os.IsExist(err) {
			return name, dir, url, err
			// fmt.Fprintln(os.Stderr, err)
			// os.Exit(1)
		}
		if err != nil || fi == nil {
			return name, dir, url, err
		}
	} else if strings.Contains(str, PathSep) {
		if strings.Contains(str, nameSep) {
			parts := strings.Split(str, nameSep)
			if len(parts) != 2 {
				err = fmt.Errorf("Invalid spec: '%s'", str)
				return
			}
			name = parts[0]
			url = parts[1]
		} else {
			// name = path.Base(str)
			url = str
		}
	} else {
		err = fmt.Errorf("Unknown git origin: '%s'", str)
	}
	return
}

func (r *Role) GetEnv() (map[string]string, error) {
	r.Config.UnmarshalKey("env", &r.Package.Env)
	env := make(map[string]string, 0)
	for k, v := range r.Package.Env {
		k = strings.ToTitle(k)
		if v == "" {
			val, ok := os.LookupEnv(k)
			if !ok {
				fmt.Printf("Warn: LookupEnv failed for '%s'", k)
			}
			v = val
		} // v = os.ExpandEnv(v)
		templ, err := template.New(k).Option("missingkey=zero").Parse(v)
		if err != nil {
			return env, err
		}
		buf := &bytes.Buffer{}
		err = templ.Execute(buf, Env())
		if err != nil {
			return env, err
		}
		v = buf.String()
		env[k] = v
	}
	return env, nil
}

func Env() map[string]string {
	env := make(map[string]string, 0)
	for _, i := range os.Environ() {
		sep := strings.Index(i, "=")
		env[i[0:sep]] = i[sep+1:]
	}
	return env
}

func (r *Role) IsOs(types []string) bool {
	if len(r.Os) == 0 { // || len(types) == 0
		return true
	}
	return hasOne(r.Os, types)
}

func hasOne(in []string, list []string) bool {
	for _, a := range in {
		for _, b := range list {
			if b == a {
				return true
			}
		}
	}
	return false
}
