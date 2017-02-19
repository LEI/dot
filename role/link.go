package role

import (
	"fmt"
	// "github.com/LEI/dot/fileutil"
	"os"
)

type Link struct {
	Pattern string
	// Target string
	Type string
}

func NewLink(pattern string) *Link {
	return &Link{Pattern: pattern}
}

func (l *Link) Link(target string) error {
	fmt.Println("DO LINK", l, target)
	// return fileutil.Link(l.Path, target)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (l *Link) Set(value string) {
	fmt.Println("Set", l, value)
	// switch val := value.(type) {
	// case string:
	// 	l.Path = val
	// 	// *l = append(*l, val)
	// default:
	// 	*l = val.(Link)
	// }
}

func (l *Link) String() string {
	str := l.Pattern
	if l.Type != "" {
		str += fmt.Sprintf("[%s]", l.Type)
	}
	return fmt.Sprintf("%s", str)
}

func (l *Link) Sync(target string) error {
	fmt.Printf("Sync: Link %s/%s\n", target, l)
	return nil
}

func (r *Role) Links() []*Link {
	p := r.Package
	if p == nil {
		p = &Package{}
	}
	// r.Config.UnmarshalKey("link", &p.Link)
	// r.Config.UnmarshalKey("links", &p.Links)
	// p.Links := make([]interface{}, 0)
	ln := r.Config.Get("link")
	if ln != nil {
		p.Links = append(p.Links, castAsLink(ln))
		p.Link = nil
	}
	lln := r.Config.Get("links")
	if lln != nil {
		for _, ln := range lln.([]interface{}) {
			p.Links = append(p.Links, castAsLink(ln))
		}
	}
	r.Config.Set("links", p.Links)
	r.Config.Set("link", p.Link)
	return p.Links
}

func castAsLink(value interface{}) *Link {
	var ln *Link
	switch v := value.(type) {
	case string:
		ln = &Link{Pattern: v}
	case map[string]interface{}:
		pattern, ok := v["pattern"].(string)
		if !ok {
			fatal(fmt.Errorf("'pattern' not found in %+v\n", v))
		}
		fileType, ok := v["type"].(string)
		if !ok {
			fileType = ""
		}
		ln = &Link{
			Pattern: pattern,
			Type: fileType,
		}
	case map[interface{}]interface{}:
		pattern, ok := v[interface{}("pattern")].(string)
		if !ok {
			fatal(fmt.Errorf("'pattern' not found in %+v\n", v))
		}
		fileType, ok := v[interface{}("type")].(string)
		if !ok {
			fileType = ""
		}
		ln = &Link{
			Pattern: pattern,
			Type: fileType,
		}
	default:
		fatal(fmt.Errorf("(%T) %s\n", v, v))
	}
	return ln
}

func fatal(msg interface{}) {
	fmt.Fprintf(os.Stderr, "Error while parsing link: %s", msg)
	os.Exit(64)
}

func fataln(msg interface{}) {
	fatal(fmt.Sprintf("%s\n", msg))
}
