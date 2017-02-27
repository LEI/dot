package role

import (
	"fmt"
)

type Line struct {
	Path string
	Line string
}

func (l *Line) String() string {
	return fmt.Sprintf("%s`%s`", l.Path, l.Line)
}

func (r *Role) GetLines() []*Line {
	if r.Package == nil {
		r.Package = &Package{}
	}
	r.Config.UnmarshalKey("line", &r.Package.Line)
	r.Config.UnmarshalKey("lines", &r.Package.Lines)
	if r.Package.Lines == nil {
		r.Package.Lines = make([]*Line, 0)
	}
	if r.Package.Line != nil {
		r.Package.Lines = append(r.Package.Lines, r.Package.Line) // .(map[string]interface{})
		r.Package.Line = nil
	}
	r.Config.Set("lines", r.Package.Lines)
	r.Config.Set("line", r.Package.Line)
	return r.Package.Lines
}
