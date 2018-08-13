package tasks

import (
	"fmt"

	"github.com/LEI/dot/cli/config/types"
	// "github.com/LEI/dot/system"
)

// Template task
type Template struct {
	Task
	Source, Target string
	Env            map[string]string
	Vars           map[string]interface{}
	// backup bool
	// overwrite bool
}

func (l *Template) String() string {
	return fmt.Sprintf("template[%s:%s]", l.Source, l.Target)
}

// Check template task
func (l *Template) Check() error {
	if l.Source == "" {
		return fmt.Errorf("template: empty source")
	}
	return nil
}

// Install template task
func (l *Template) Install() error {
	cmd := fmt.Sprintf("gotpl %s %s", l.Source, l.Target)
	if !l.ShouldInstall() {
		if Verbose > 0 {
			fmt.Fprintf(Stdout, "# %s\n", cmd)
		}
		return ErrSkip
	}
	fmt.Fprintf(Stdout, "$ %s\n", cmd)
	return nil // system.Template(l.Source, l.Target)
}

// Remove template task
func (l *Template) Remove() error {
	cmd := fmt.Sprintf("rm %s", l.Target)
	if !l.ShouldRemove() {
		if Verbose > 0 {
			fmt.Fprintf(Stdout, "# %s\n", cmd)
		}
		return ErrSkip
	}
	fmt.Fprintf(Stdout, "$ %s\n", cmd)
	return nil // system.Remove(l.Target)
}

// Templates task slice
type Templates []*Template

func (templates *Templates) String() string {
	// s := ""
	// for i, l := range *templates {
	// 	s += fmt.Sprintf("%s", l)
	// 	if i > 0 {
	// 		s += "\n"
	// 	}
	// }
	// return s
	return fmt.Sprintf("%s", *templates)
}

// Parse template tasks
func (templates *Templates) Parse(i interface{}) error {
	ll := &Templates{}
	m, err := types.NewMapPaths(i)
	if err != nil {
		return err
	}
	for k, v := range *m {
		l := &Template{
			Source: k,
			Target: v,
		}
		// *ll = append(*ll, l)
		ll.Add(*l)
	}
	*templates = *ll
	return nil
}

// Add a dir
func (templates *Templates) Add(l Template) {
	*templates = append(*templates, &l)
}
