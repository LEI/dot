package dot

import "fmt"

// Templates task list
type Templates []*Template

// Template task
type Template struct {
	Source string
	Target string
	Env    map[string]string
	Vars   map[string]interface{}
}

func (t *Template) String() string {
	return fmt.Sprintf("tpl %s %s", t.Source, t.Target)
}
