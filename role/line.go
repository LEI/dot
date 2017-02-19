package role

import (
	"fmt"
)

type Line struct {
	File string
	Line string
}

func (l *Line) String() string {
	return fmt.Sprintf("%s`%s`", l.File, l.Line)
}

func (r *Role) Lines() []*Line {
	p := r.Package
	if p == nil {
		p = &Package{}
	}
	r.Config.UnmarshalKey("line", &p.Line)
	r.Config.UnmarshalKey("lines", &p.Lines)
	if p.Line != nil {
		p.Lines = append(p.Lines, p.Line) // .(map[string]interface{})
		p.Line = nil
	}
	r.Config.Set("lines", p.Lines)
	r.Config.Set("line", p.Line)
	return p.Lines
}
