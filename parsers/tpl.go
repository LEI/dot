package parsers

import (
	"fmt"
	// "reflect"

	// "github.com/LEI/dot/dotfile"
)

// Tpl ...
type Tpl struct {
	Name  string
	Paths *Paths
}

// Templates ...
type Templates []*Tpl

func (p *Templates) String() string {
	s := ""
	for _, v := range *p {
		s+= fmt.Sprintf("%+v", v)
	}
	return s
}

// UnmarshalYAML ...
func (p *Templates) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// TODO c.f. paths
	return nil
}
