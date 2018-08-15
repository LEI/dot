package tasks

import (
	"fmt"

	"github.com/LEI/dot/cli/config/types"
	"github.com/LEI/dot/system"
)

// Line task
type Line struct {
	Task
	File, Line string
}

func (l *Line) String() string {
	return fmt.Sprintf("line[%s:%s]", l.File, l.Line)
}

// Check line task
func (l *Line) Check() error {
	if l.File == "" {
		return fmt.Errorf("line: missing file")
	}
	// if l.Line == "" {
	// 	return fmt.Errorf("line: missing line")
	// }
	err := system.CheckLine(l.File, l.Line)
	switch err {
	case system.ErrLineAlreadyExist:
		l.Done()
	default:
		return err
	}
	return nil
}

// Install line task
func (l *Line) Install() error {
	cmd := fmt.Sprintf("echo '%s' >> %s", l.Line, l.File)
	if !l.ShouldInstall() {
		if Verbose > 0 {
			fmt.Fprintf(Stdout, "# %s\n", cmd)
		}
		return ErrSkip
	}
	fmt.Fprintf(Stdout, "$ %s\n", cmd)
	return system.CreateLine(l.File, l.Line)
}

// Remove line task
func (l *Line) Remove() error {
	cmd := fmt.Sprintf("sed -i '#^%s$#d' %s", l.Line, l.File)
	if !l.ShouldRemove() {
		if Verbose > 0 {
			fmt.Fprintf(Stdout, "# %s\n", cmd)
		}
		return ErrSkip
	}
	fmt.Fprintf(Stdout, "$ %s\n", cmd)
	return system.RemoveLine(l.File, l.Line)
}

// Lines task slice
type Lines []*Line

func (lines *Lines) String() string {
	return fmt.Sprintf("%s", *lines)
}

// Parse line tasks
func (lines *Lines) Parse(i interface{}) error {
	ll := &Lines{}
	m, err := types.NewMap(i)
	if err != nil {
		return err
	}
	for k, v := range *m {
		l := &Line{
			File: k,
			Line: v.(string),
		}
		// *ll = append(*ll, l)
		ll.Add(*l)
	}
	*lines = *ll
	return nil
}

// Add a line
func (lines *Lines) Add(l Line) {
	*lines = append(*lines, &l)
}
