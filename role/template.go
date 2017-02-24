package role

import (
	"fmt"
)

type Template struct {
	Path string
	Vars map[string]string
}

func (t *Template) String() string {
	return fmt.Sprintf("%s", t.Path)
}

func (r *Role) Templates() []*Template {
	p := r.Package
	if p == nil {
		p = &Package{}
	}
	r.Config.UnmarshalKey("template", &p.Template)
	r.Config.UnmarshalKey("templates", &p.Templates)
	if p.Template != nil {
		p.Templates = append(p.Templates, p.Template) // .(map[string]interface{})
		p.Template = nil
	}
	r.Config.Set("templates", p.Templates)
	r.Config.Set("template", p.Template)
	return p.Templates
}
