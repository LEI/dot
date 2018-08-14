package tasks

import (
	"fmt"

	"github.com/LEI/dot/cli/config/types"
	// "github.com/LEI/dot/system"
	"github.com/LEI/dot/pkg/pkg"
	// "github.com/LEI/dot/pkg/ostype"
	"github.com/mitchellh/mapstructure"
)

// Package task
type Package struct {
	Task
	Type   string // package manager
	Name   string
	Opts   types.Slice // todo args -> opts
	types.HasOS `mapstructure:",squash"`
	// sudo   bool
	// latest bool
}

func (p *Package) String() string {
	return fmt.Sprintf("package[%s]", p.Name)
}

// Check package task
func (p *Package) Check() error {
	if p.Name == "" {
		return fmt.Errorf("package: empty name")
	}
	if !p.CheckOS() { // len(p.OS) > 0 && !ostype.Has(p.OS...) {
		return fmt.Errorf("package %s: only for %s", p.Name, p.OS)
	}
	return nil
}

// Install package task
func (p *Package) Install() error {
	if !p.CheckOS() { // len(p.OS) > 0 && !ostype.Has(p.OS...) {
		return ErrSkip // fmt.Errorf("package %s: only for %s", p.Name, p.OS)
	}
	// str := fmt.Sprintf("pkg install %s %s", p.Name, p.Opts)
	if !p.ShouldInstall() {
		// if Verbose > 0 {
		// 	fmt.Fprintf(Stdout, "# %s\n", str)
		// }
		return ErrSkip
	}
	// fmt.Fprintf(Stdout, "$ %s\n", str)
	return pkg.Install(p.Type, p.Name, p.Opts...)
}

// Remove package task
func (p *Package) Remove() error {
	if !p.CheckOS() { // len(p.OS) > 0 && !ostype.Has(p.OS...) {
		return ErrSkip // fmt.Errorf("package %s: only for %s", p.Name, p.OS)
	}
	// str := fmt.Sprintf("pkg remove %s %s", p.Name, p.Opts)
	if !p.ShouldRemove() {
		// if Verbose > 0 {
		// 	fmt.Fprintf(Stdout, "# %s\n", str)
		// }
		return ErrSkip
	}
	// fmt.Fprintf(Stdout, "$ %s\n", str)
	return pkg.Remove(p.Type, p.Name, p.Opts...)
}

// Packages task slice
type Packages []*Package

func (packages *Packages) String() string {
	// s := ""
	// for i, p := range *packages {
	// 	s += fmt.Sprintf("%s", p)
	// 	if i > 0 {
	// 		s += "\n"
	// 	}
	// }
	// return s
	return fmt.Sprintf("%s", *packages)
}

// Parse package tasks
func (packages *Packages) Parse(i interface{}) error {
	pp := &Packages{}
	s, err := types.NewSliceMap(i)
	if err != nil {
		return err
	}
	for _, v := range *s {
		p := &Package{}
		switch val := v.(type) {
		case string:
			// p = &Exec{Command: val}
			p.Name = val
		case *types.Map:
			mapstructure.Decode(val, &p)
		case interface{}:
			p = val.(*Package)
		default:
			return fmt.Errorf("invalid pkg map: %+v", v)
		}
		// fmt.Printf("COMMAND [%+v] = [%+v]\n", v, *p)
		pp.Add(*p)
	}
	*packages = *pp
	return nil
}

// Add a dir
func (packages *Packages) Add(p Package) {
	*packages = append(*packages, &p)
}
