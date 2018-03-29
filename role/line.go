package role

import (
	"fmt"
	"github.com/LEI/dot/log"
)

type Line struct {
	Path string
	Line string
}

func (l *Line) String() string {
	return fmt.Sprintf("%s`%s`", l.Path, l.Line)
}

// func (r *Role) GetLines() []*Line {
// 	if r.Package == nil {
// 		r.Package = &Package{}
// 	}
// 	if r.Config.IsSet("line") {
// 		r.Config.UnmarshalKey("line", &r.Package.Line)
// 	}
// 	if r.Config.IsSet("lines") {
// 		p := map[string]interface{}{}
// 		r.Config.UnmarshalKey("lines", &p)
// 		r.Package.Lines = p
// 	}
// 	if r.Package.Lines == nil {
// 		r.Package.Lines = make([]*Line, 0)
// 	}
// 	if r.Package.Line != nil {
// 		r.Package.Lines = append(r.Package.Lines, r.Package.Line) // .(map[string]interface{})
// 		r.Package.Line = nil
// 	}
// 	log.Printf("pkg %+v\n", r.Config)
// 	log.Printf("LINES %+v\n", r.Package.Lines)
// 	log.Printf("LINE %+v\n", r.Package.Line)
// 	r.Config.Set("lines", r.Package.Lines)
// 	r.Config.Set("line", r.Package.Line)
// 	return r.Package.Lines
// }

func (r *Role) GetLines() []*Line {
	if r.Package == nil {
		r.Package = &Package{}
	}
	// r.Config.UnmarshalKey("line", &r.Package.Line)
	// r.Config.UnmarshalKey("lines", &r.Package.Lines)
	if r.Package.Lines == nil {
		r.Package.Lines = make([]*Line, 0)
	}
	if r.Config.IsSet("line") {
		ln := r.Config.Get("line")
		if ln != nil { // ! reflect.ValueOf(ln).IsNil()
			r.Package.Lines = append(r.Package.Lines, castAsLine(ln, ""))
			r.Package.Line = nil
		}
	}
	if r.Config.IsSet("lines") {
		lines := r.Config.Get("lines")
		if lines != nil { // ! reflect.ValueOf(lines).IsNil()
			for k, ln := range lines.(map[string]interface{}) { // .([]*Line)
				r.Package.Lines = append(r.Package.Lines, castAsLine(ln, k))
			}
		}
	}
	r.Config.Set("lines", r.Package.Lines)
	r.Config.Set("line", r.Package.Line)
	return r.Package.Lines
}

func castAsLine(value interface{}, key string) *Line {
	var l *Line
	switch v := value.(type) {
	case string:
		l = &Line{Path: key, Line: v}
	case map[string]interface{}:
		p, ok := v["path"].(string)
		if !ok {
			log.Fatal(fmt.Errorf("'path' not found in %+v\n", v))
		}
		str, ok := v["line"].(string)
		if !ok {
			log.Fatal(fmt.Errorf("'line' not found in %+v\n", v))
		}
		l = &Line{
			Path: p,
			Line: str,
		}
	case map[interface{}]interface{}:
		p, ok := v[interface{}("path")].(string)
		if !ok {
			log.Fatal(fmt.Errorf("'path' not found in %+v\n", v))
		}
		str, ok := v[interface{}("line")].(string)
		if !ok {
			log.Fatal(fmt.Errorf("'line' not found in %+v\n", v))
		}
		l = &Line{
			Path: p,
			Line: str,
		}
	default:
		log.Fatal(fmt.Errorf("Unhandled type (%T) %s for %v\n", v, v, value))
	}
	return l
}
