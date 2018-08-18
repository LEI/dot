package dot

import (
	"fmt"
	"os"
	"strings"

	"github.com/LEI/dot/internal/pkg"
)

// Pkg task
type Pkg struct {
	Task
	Name    string
	Args    []string
	OS      []string
	Action  string
	Manager string `mapstructure:"type"`
}

func (p *Pkg) String() string {
	// return fmt.Sprintf("%s %s %s %s %s", p.Manager, p.Action, p.Name, p.Args, p.OS)
	return p.Name // fmt.Sprintf("%s %s", p.Name, p.Args)
}

// DoString string
func (p *Pkg) DoString() string {
	bin, opts, err := pkg.Build("install", p.Manager, p.Name, p.Args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return ""
	}
	return fmt.Sprintf("%s %s", bin, strings.Join(opts, " "))
}

// UndoString string
func (p *Pkg) UndoString() string {
	bin, opts, err := pkg.Build("remove", p.Manager, p.Name, p.Args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return ""
	}
	return fmt.Sprintf("%s %s", bin, strings.Join(opts, " "))
}

// Status check task
func (p *Pkg) Status() error {
	// if hookExists(p.Target) {
	// 	return ErrAlreadyExist
	// }
	return nil // ErrAlreadyExist
}

// Do task
func (p *Pkg) Do() error {
	// TODO OS, If, Action
	return nil // pkg.Install(p.Manager, p.Name, p.Args...)
}

// Undo task
func (p *Pkg) Undo() error {
	// TODO OS, If, Action
	return nil // pkg.Remove(p.Manager, p.Name, p.Args...)
}
