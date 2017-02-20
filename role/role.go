package role

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path"
	"strings"
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

func (r *Role) New(source, target string) (*Role, error) {
	if r == nil {
		r = &Role{}
	}
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
