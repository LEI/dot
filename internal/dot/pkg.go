package dot

import (
	"fmt"
	"strings"
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
	return fmt.Sprintf("Install %s package %s %s", p.Type, p.Name, strings.Join(p.Args, " "))
}

// UndoString string
func (p *Pkg) UndoString() string {
	return fmt.Sprintf("Remove %s package %s %s", p.Type, p.Name, strings.Join(p.Args, " "))
}

// Status check task
func (p *Pkg) Status() error {
	// if hookExists(p.Target) {
	// 	return ErrAlreadyExist
	// }
	return nil
}

// Do task
func (p *Pkg) Do() error {
	return fmt.Errorf("not implemented")
}

// Undo task
func (p *Pkg) Undo() error {
	return fmt.Errorf("not implemented")
}
