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

func (r *Role) GetDirs() []*Dir {
	// r.Config.UnmarshalKey("dir", &r.Package.Dir)
	// r.Config.UnmarshalKey("dirs", &r.Package.Dirs)
	if r.Package == nil {
		r.Package = &Package{}
	}
	if r.Package.Dirs == nil {
		r.Package.Dirs = make([]*Dir, 0)
	}
	dir := r.Config.GetString("dir")
	if dir != "" {
		r.Package.Dirs = append(r.Package.Dirs, &Dir{Path: dir})
		r.Package.Dir = nil
	}
	for _, d := range r.Config.GetStringSlice("dirs") {
		r.Package.Dirs = append(r.Package.Dirs, &Dir{Path: d})
	}
	r.Config.Set("dirs", r.Package.Dirs)
	r.Config.Set("dir", r.Package.Dir)
	return r.Package.Dirs
}
