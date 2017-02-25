package role

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path"
	"strings"
	"text/template"
)

var (
	PathSep = string(os.PathSeparator)
)

type Meta struct {
	Source, Target string
	Roles          []*Role
}

func (m *Meta) String() string {
	return fmt.Sprintf("%s -> %s = %+v", m.Source, m.Target, m.Roles)
}

type Role struct {
	Name, Origin   string
	Source, Target string
	Os             []string
	Config         *viper.Viper
	Package        *Package // `mapstructure:",squash"`
}

// func (r *Role) New(v interface{}) *Role {
// 	return &Role{}
// }

// func (r *Role) String() string {
// 	return fmt.Sprintf("%s (%s) [%s -> %s] Dir: %s, Dirs: %v, Link: %s, Links: %v, Lines: %+v",
// 		r.Name, r.Origin, r.Source, r.Target, r.Dir, r.Dirs, r.Link, r.Links, r.Lines)
// }

// func (r *Role) Origin() string {
// 	origin := r.Source
// 	return fmt.Sprintf("%s", origin)
// }

func Roles() *Meta {
	return &Meta{Roles:make([]*Role, 0)}
}

func (m *Meta) ParseRoles() error {
	for i, r := range m.Roles {
		r, err := r.Init(m.Source, m.Target)
		if err != nil {
			return err
		}
		m.Roles[i] = r
	}
	return nil
}

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
