package role

import (
	"fmt"
)

type Dir struct {
	Path string
}

func (d *Dir) String() string {
	return fmt.Sprintf("%s", d.Path)
}

func (r *Role) Dirs() []*Dir {
	p := r.Package
	if p == nil {
		p = &Package{}
	}
	// r.Config.UnmarshalKey("dir", &p.Dir)
	// r.Config.UnmarshalKey("dirs", &p.Dirs)
	dir := r.Config.GetString("dir")
	if dir != "" {
		p.Dirs = append(p.Dirs, &Dir{Path: dir})
		p.Dir = nil
	}
	for _, d := range r.Config.GetStringSlice("dirs") {
		p.Dirs = append(p.Dirs, &Dir{Path: d})
	}
	r.Config.Set("dirs", p.Dirs)
	r.Config.Set("dir", p.Dir)
	return p.Dirs
}
