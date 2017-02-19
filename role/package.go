package role

import (
	// "fmt"
)

// type Package map[string]interface{}
type Package struct {
	Dir string
	Dirs []string
	Link interface{}
	Links []interface{}
	Line []interface{}
	Lines []interface{}
	Template interface{}
}

// func (p *Package) String() string {
// 	return fmt.Sprintf("%+v", p)
// }

func (p *Package) GetDirs() []string {
	if p.Dir != "" {
		p.Dirs = append(p.Dirs, p.Dir)
		p.Dir = ""
	}
	return p.Dirs
}

func (p *Package) GetLinks() []interface{} {
	if p.Link != nil {
		p.Links = append(p.Links, p.Link)
		p.Link = nil
	}
	return p.Links
}

func (p *Package) GetLines() []interface{} {
	if p.Line != nil {
		p.Lines = append(p.Lines, p.Line)
		p.Line = nil
	}
	return p.Lines
}
