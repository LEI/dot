package role

import (
	"fmt"
	"github.com/LEI/dot/log"
	// "os"
	// "path/filepath"
)

type Link struct {
	Path string
	// Source string
	// Target string
	Type  string
}

func (l *Link) String() string {
	str := l.Path
	if l.Type != "" {
		str += fmt.Sprintf("[%s]", l.Type)
	}
	return fmt.Sprintf("%s", str)
}

func (r *Role) Links() []*Link {
	p := r.Package
	if p == nil {
		p = &Package{}
	}
	// r.Config.UnmarshalKey("link", &p.Link)
	// r.Config.UnmarshalKey("links", &p.Links)
	// p.Links := make([]interface{}, 0)
	l := r.Config.Get("link")
	if l != nil {
		p.Links = append(p.Links, castAsLink(l))
		p.Link = nil
	}
	links := r.Config.Get("links")
	if links != nil {
		for _, l := range links.([]interface{}) {
			p.Links = append(p.Links, castAsLink(l))
		}
	}
	r.Config.Set("links", p.Links)
	r.Config.Set("link", p.Link)
	return p.Links
}

func castAsLink(value interface{}) *Link {
	var l *Link
	switch v := value.(type) {
	case string:
		l = &Link{Path: v}
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
			Type:    fileType,
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
			Type:    fileType,
		}
	default:
		log.Fatal(fmt.Errorf("(%T) %s\n", v, v))
	}
	return l
}
