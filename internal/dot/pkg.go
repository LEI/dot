package dot

import (
	"fmt"
	"strings"

	"github.com/LEI/dot/internal/pkg"
)

// Pkg task
type Pkg struct {
	Task
	Name   string
	Args   []string
	OS     []string
	Action string
	Type   string
}

func (p *Pkg) String() string {
	// return fmt.Sprintf("%s %s %s %s %s", p.Type, p.Action, p.Name, p.Args, p.OS)
	return p.Name // fmt.Sprintf("%s %s", p.Name, p.Args)
}

// DoString string
func (p *Pkg) DoString() string {
	return fmt.Sprintf("install %s %s", p.Name, strings.Join(p.Args, " "))
}

// UndoString string
func (p *Pkg) UndoString() string {
	return fmt.Sprintf("remove %s %s", p.Name, strings.Join(p.Args, " "))
}

// Status check task
func (p *Pkg) Status() error {
	// if hookExists(p.Target) {
	// 	return ErrAlreadyExist
	// }
	// return nil
	return ErrAlreadyExist
}

// Do task
func (p *Pkg) Do() error {
	// TODO OS, Action
	return pkg.Install(p.Type, p.Name, p.Args...)
}

// Undo task
func (p *Pkg) Undo() error {
	// TODO OS, Action
	return pkg.Remove(p.Type, p.Name, p.Args...)
}
