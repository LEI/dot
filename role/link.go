package role

import (
	"fmt"
	"os"
	"path/filepath"
)

var IgnoreNames = []string{".git", ".*\\.md"}

type Link struct {
	Pattern string
	// Source string
	// Target string
	Type string
	Files []string // map[string]*os.FileInfo
}

func NewLink(pattern string) *Link {
	return &Link{Pattern: pattern}
}

func (l *Link) String() string {
	str := l.Pattern
	if l.Type != "" {
		str += fmt.Sprintf("[%s]", l.Type)
	}
	return fmt.Sprintf("%s", str)
}

func (l *Link) GlobFiles(source string) ([]string, error) {
	glob := filepath.Join(source, l.Pattern)
	paths, err := filepath.Glob(glob)
	if err != nil {
		return paths, err
	}
	GLOB:
	for _, file := range paths {
		base := filepath.Base(file)
		for _, pattern := range IgnoreNames {
			ignore, err := filepath.Match(pattern, base)
			if err != nil {
				return paths, err
			}
			if ignore {
				fmt.Printf("# ignore: %s\n", file)
				continue GLOB
			}
		}
		fi, err := os.Stat(file)
		if err != nil {
			return paths, err
		}
		switch {
		case l.Type == "directory" && !fi.IsDir(),
			l.Type == "file" && fi.IsDir():
			fmt.Printf("# ignore: %s (not a %s)\n", file, l.Type)
			continue // GLOB
		}
		l.Files = append(l.Files, file)
		// l.Files[file] = fi
	}

	return l.Files, nil
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
		l = &Link{Pattern: v}
	case map[string]interface{}:
		pattern, ok := v["pattern"].(string)
		if !ok {
			fatal(fmt.Errorf("'pattern' not found in %+v\n", v))
		}
		fileType, ok := v["type"].(string)
		if !ok {
			fileType = ""
		}
		l = &Link{
			Pattern: pattern,
			Type:    fileType,
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
		l = &Link{
			Pattern: pattern,
			Type:    fileType,
		}
	default:
		fatal(fmt.Errorf("(%T) %s\n", v, v))
	}
	return l
}
