package role

import (
	"fmt"
	"os"
	"path"
	"strings"
)

var (
	PathSep = string(os.PathSeparator)
)

type Meta struct {
	Source, Target string
	Roles []*Role
}

func (m *Meta) String() string {
	return fmt.Sprintf("%s -> %s = %+v", m.Source, m.Target, m.Roles)
}

type Role struct {
	Name, Origin string
	Source, Target string
	Os []string
	Dir string
	Dirs []string
	Link interface{}
	Links []interface{}
	Lines map[string]string
}

// func (r *Role) New(v interface{}) *Role {
// 	return &Role{}
// }

// func (r *Role) String() string {
// 	return fmt.Sprintf("%s: %s -> %s", r.Name, r.Source, r.Target)
// }

// func (r *Role) Origin() string {
// 	origin := r.Source
// 	return fmt.Sprintf("%s", origin)
// }

func (r *Role) Init(source, target string) error {
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
		return fmt.Errorf("Invalid role: %+v\n", r)
	}
	if r.Source == "" {
		r.Source = source
	}
	if r.Target == "" {
		r.Target = target
	}
	return nil
}

func (r *Role) IsOs(types []string) bool {
	if len(r.Os) == 0 || len(types) == 0 {
		return true
	}
	for _, o := range r.Os {
		for _, s := range types {
			if o == s {
				return true
			}
		}
	}
	return false
}

func (r *Role) GetDirs() []string {
	if r.Dir != "" {
		r.Dirs = append(r.Dirs, r.Dir)
		r.Dir = ""
	}
	return r.Dirs
}
