package dot

import (
	"fmt"
	"os"

	"github.com/LEI/dot/internal/cli"
	"github.com/LEI/dot/internal/pkg"
)

// Pkg task
type Pkg struct {
	Task    `mapstructure:",squash"` // Action, If, OS
	Name    []string
	Args    []string
	Manager string `mapstructure:"type"`
}

func (p *Pkg) String() string {
	// return fmt.Sprintf("%s %s %s %s %s", p.Manager, p.Action, p.Name, p.Args, p.OS)
	// fmt.Sprintf("%s %s", p.Name, p.Args)
	return cli.FormatArgs(append(p.Name, p.Args...))
}

// Type task name
func (p *Pkg) Type() string {
	return "pkg"
}

// DoString string
func (p *Pkg) DoString() string {
	m, err := pkg.NewPm(p.Manager)
	if err != nil {
		return "<none>"
	}
	bin, opts, err := m.Build("install", p.Name, p.Args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err pkg do: %s\n", err)
		return ""
	}
	return fmt.Sprintf("%s %s", bin, cli.FormatArgs(opts))
}

// UndoString string
func (p *Pkg) UndoString() string {
	m, err := pkg.NewPm(p.Manager)
	if err != nil {
		return "<none>"
	}
	bin, opts, err := m.Build("remove", p.Name, p.Args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err pkg undo: %s\n", err)
		return ""
	}
	return fmt.Sprintf("%s %s", bin, cli.FormatArgs(opts))
}

// Status check task
func (p *Pkg) Status() error {
	exists, err := pkg.Has(p.Manager, p.Name, p.Args...)
	if err != nil {
		return err
	}
	if exists {
		return ErrExist
	}
	return nil
}

// Do task
func (p *Pkg) Do() error {
	if err := p.Status(); err != nil {
		switch err {
		case ErrExist, ErrSkip:
			return nil
		default:
			return err
		}
	}
	// fmt.Println("pkg.Install", p.Manager, p.Name, p.Args)
	err := pkg.Install(p.Manager, p.Name, p.Args...)
	// fmt.Println(p.Name, "=======>", err)
	switch err {
	// case nil:
	case pkg.ErrExist:
		return ErrExist
	default:
		return err
	}
}

// Undo task
func (p *Pkg) Undo() error {
	if err := p.Status(); err != nil {
		switch err {
		case ErrExist:
			// continue
		case ErrSkip:
			return nil
		default:
			return err
		}
	}
	err := pkg.Remove(p.Manager, p.Name, p.Args...)
	switch err {
	// case nil:
	case pkg.ErrExist:
		return ErrExist
	default:
		return err
	}
}
