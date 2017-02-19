package role

import (
	"fmt"
)

type Dir struct {
	Path string
}

// func (d *Dir) Get() string {
// 	return fmt.Sprintf("%s", d)
// }

// func (d *Dir) Set(value interface{}) {
// 	switch val := value.(type) {
// 	case string:
// 		d.Path = val
// 		// *d = append(*d, val)
// 	// case []string:
// 	// 	*d = make(Dir, 0)
// 	// 	for _, v := range val {
// 	// 		*d = append(*d, v)
// 	// 	}
// 	// }
// 	default:
// 		*d = val.(Dir)
// 	}
// }

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
