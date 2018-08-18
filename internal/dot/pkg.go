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
	bin, opts, err := pkg.Init(p.Manager, "install", p.Name, p.Args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return ""
	}
	return fmt.Sprintf("%s %s", bin, strings.Join(opts, " "))
}

// UndoString string
func (p *Pkg) UndoString() string {
	bin, opts, err := pkg.Init(p.Manager, "remove", p.Name, p.Args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return ""
	}
	return fmt.Sprintf("%s %s", bin, strings.Join(opts, " "))
}

// Status check task
func (p *Pkg) Status() error {
	ok, err := pkg.Has(p.Manager, p.Name, p.Args...)
	if err != nil {
		return err
	}
	if ok {
		return ErrAlreadyExist
	}
	return nil
}

// Do task
func (p *Pkg) Do() error {
	// if err := c.Status(); err != nil {
	// 	if err == ErrAlreadyExist {
	// 		return nil
	// 	}
	// 	return err
	// }
	// TODO OS, If, Action
	err := pkg.Install(p.Manager, p.Name, p.Args...)
	switch err {
	// case nil:
	case pkg.ErrExist:
		return ErrAlreadyExist
	default:
		return err
	}
}

// Undo task
func (p *Pkg) Undo() error {
	// if err := h.Status(); err != nil {
	// 	if err != ErrAlreadyExist {
	// 		return err
	// 	}
	// }
	// TODO OS, If, Action
	err := pkg.Remove(p.Manager, p.Name, p.Args...)
	switch err {
	// case nil:
	case pkg.ErrExist:
		return ErrAlreadyExist
	default:
		return err
	}
}
