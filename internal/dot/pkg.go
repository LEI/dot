package dot

import (
	"fmt"
	"os"

	"github.com/LEI/dot/internal/pkg"
	"github.com/LEI/dot/internal/shell"
)

// Pkg task
type Pkg struct {
	Task    `mapstructure:",squash"` // Action, If, OS
	Name    []string
	Args    []string
	Manager string `mapstructure:"type"`
}

// NewPkg task
func NewPkg(s string) *Pkg {
	return &Pkg{Name: []string{s}}
}

func (p *Pkg) String() string {
	// switch Action {
	// case "install":
	// case "remove":
	// }
	m, err := pkg.NewPm(p.Manager)
	if err != nil || Action == "" {
		a := []string{}
		if p.Manager != "" {
			a = append(a, p.Manager)
		}
		a = append(a, p.Name...)
		a = append(a, p.Args...)
		return fmt.Sprint(shell.FormatArgs(a))
	}
	bin, opts, err := m.Build(Action, p.Name, p.Args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err pkg do: %s\n", err)
		return ""
	}
	return fmt.Sprintf("%s %s", bin, shell.FormatArgs(opts))
}

// Status check task
func (p *Pkg) Status() error {
	// TODO: upgrade flag (this will prevent any pkg.Pm.Install
	// function to execute even if it can upgrade the package
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
