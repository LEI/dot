package role

import (
	// "bytes"
	"fmt"
	// "os"
	// "strings"
	// "text/template"
)

type Template struct {
	Path string
	// Vars map[string]string
	// Tmpl template.Template
}

func (t *Template) String() string {
	return fmt.Sprintf("%s", t.Path)
}

func (r *Role) GetTemplates() []*Template {
	if r.Package == nil {
		r.Package = &Package{}
	}
	// r.Config.UnmarshalKey("template", &r.Package.Template)
	// r.Config.UnmarshalKey("templates", &r.Package.Templates)
	// if r.Package.Template != nil {
	// 	r.Package.Templates = append(r.Package.Templates, r.Package.Template) // .(map[string]interface{})
	// 	r.Package.Template = nil
	// }
	tpl := r.Config.GetString("template")
	if tpl != "" {
		r.Package.Templates = append(r.Package.Templates, &Template{Path: tpl})
		r.Package.Template = nil
	}
	for _, t := range r.Config.GetStringSlice("templates") {
		r.Package.Templates = append(r.Package.Templates, &Template{Path: t})
	}
	r.Config.Set("templates", r.Package.Templates)
	r.Config.Set("template", r.Package.Template)
	return r.Package.Templates
}
