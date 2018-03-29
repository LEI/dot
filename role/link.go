package role

import (
	"fmt"
	// "reflect"
	"github.com/LEI/dot/log"
	// "os"
	// "path/filepath"
)

type Link struct {
	Path string
	// Source string
	// Target string
	Type string
}

func (l *Link) String() string {
	str := l.Path
	if l.Type != "" {
		str += fmt.Sprintf("[%s]", l.Type)
	}
	return fmt.Sprintf("%s", str)
}

func (r *Role) GetLinks() []*Link {
	if r.Package == nil {
		r.Package = &Package{}
	}
	// r.Config.UnmarshalKey("link", &r.Package.Link)
	// r.Config.UnmarshalKey("links", &r.Package.Links)
	if r.Package.Links == nil {
		r.Package.Links = make([]*Link, 0)
	}
	if r.Config.IsSet("link") {
		ln := r.Config.Get("link")
		if ln != nil { // ! reflect.ValueOf(ln).IsNil()
			r.Package.Links = append(r.Package.Links, castAsLink(ln))
			r.Package.Link = nil
		}
	}
	if r.Config.IsSet("links") {
		links := r.Config.Get("links")
		if links != nil { // ! reflect.ValueOf(links).IsNil()
			for _, ln := range links.([]interface{}) { // .([]*Link)
				r.Package.Links = append(r.Package.Links, castAsLink(ln))
			}
		}
	}
	r.Config.Set("links", r.Package.Links)
	r.Config.Set("link", r.Package.Link)
	return r.Package.Links
}

func castAsLink(value interface{}) *Link {
	var l *Link
	switch v := value.(type) {
	case string:
		l = &Link{Path: v}
	// case *Link:
	// 	log.Fatal(fmt.Errorf("??? (%T) %s for %v\n", v, v, value))
	// 	l = v
	// case []interface{}:
	// 	log.Fatal(fmt.Errorf("??? (%T) %s for %v\n", v, v, value))
	case map[string]interface{}:
		p, ok := v["path"].(string)
		if !ok {
			log.Fatal(fmt.Errorf("'path' not found in %+v\n", v))
		}
		fileType, ok := v["type"].(string)
		if !ok {
			fileType = ""
		}
		l = &Link{
			Path: p,
			Type: fileType,
		}
	case map[interface{}]interface{}:
		p, ok := v[interface{}("path")].(string)
		if !ok {
			log.Fatal(fmt.Errorf("'path' not found in %+v\n", v))
		}
		fileType, ok := v[interface{}("type")].(string)
		if !ok {
			fileType = ""
		}
		l = &Link{
			Path: p,
			Type: fileType,
		}
	default:
		log.Fatal(fmt.Errorf("Unhandled type (%T) %s for %v\n", v, v, value))
	}
	return l
}
