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

func (t *Pkg) String() string {
	// switch Action {
	// case "install":
	// case "remove":
	// }
	m, err := pkg.NewPm(t.Manager)
	if err != nil || Action == "" {
		s := []string{} // m.Bin
		if t.Manager != "" {
			s = append(s, t.Manager)
		}
		s = append(s, t.Name...)
		s = append(s, t.Args...)
		return fmt.Sprint(shell.FormatArgs(s))
	}
	opts, err := m.BuildOptions(Action, t.Name, t.Args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pkg: m.BuildOptions errored with: %s\n", err)
		return "<err>"
	}
	opts = m.ParseOpts(opts)
	return fmt.Sprintf("%s %s", m.Bin, shell.FormatArgs(opts))
}

// Init task
func (t *Pkg) Init() error {
	switch Action {
	case "install":
		// Update or upgrade package manager before installation
		if err := pkg.Init(t.Manager); err != nil {
			return err
		}
	case "remove":
	}
	return nil
}

// Status check task
func (t *Pkg) Status() error {
	// TODO: upgrade flag (this will prevent any pkg.Pm.Install
	// function to execute even if it can upgrade the package
	exists, err := pkg.Has(t.Manager, t.Name, t.Args...)
	if err != nil {
		return err
	}
	if exists {
		return ErrExist
	}
	return nil
}

// Do task
func (t *Pkg) Do() error {
	if err := t.Status(); err != nil {
		switch err {
		case ErrExist, ErrSkip:
			return nil
		default:
			return err
		}
	}
	// fmt.Println("pkg.Install", t.Manager, t.Name, t.Args)
	err := pkg.Install(t.Manager, t.Name, t.Args...)
	// fmt.Println(t.Name, "=======>", err)
	switch err {
	// case nil:
	case pkg.ErrExist:
		return ErrExist
	default:
		return err
	}
}

// Undo task
func (t *Pkg) Undo() error {
	if err := t.Status(); err != nil {
		switch err {
		case ErrExist:
			// continue
		case ErrSkip:
			return nil
		default:
			return err
		}
	}
	err := pkg.Remove(t.Manager, t.Name, t.Args...)
	switch err {
	// case nil:
	case pkg.ErrExist:
		return ErrExist
	default:
		return err
	}
}
