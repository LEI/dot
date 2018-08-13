package tasks

import (
	"fmt"

	"github.com/LEI/dot/cli/config/types"
	// "github.com/LEI/dot/system"
	"github.com/mitchellh/mapstructure"
)

// Package task
type Package struct {
	Task
	Name string
	Opts []string
	// latest bool
}

func (l *Package) String() string {
	return fmt.Sprintf("package[%s]", l.Name)
}

// Check package task
func (l *Package) Check() error {
	if l.Name == "" {
		return fmt.Errorf("package: empty name")
	}
	return nil
}

// Install package task
func (l *Package) Install() error {
	cmd := fmt.Sprintf("pkg %s %s", l.Name, l.Opts)
	if !l.ShouldInstall() {
		if Verbose > 0 {
			fmt.Fprintf(Stdout, "# %s\n", cmd)
		}
		return ErrSkip
	}
	fmt.Fprintf(Stdout, "$ %s\n", cmd)
	return nil // system.InstallPackage(l.Name, l.Opts)
}

// Remove package task
func (l *Package) Remove() error {
	cmd := fmt.Sprintf("rm %s", l.Name)
	if !l.ShouldRemove() {
		if Verbose > 0 {
			fmt.Fprintf(Stdout, "# %s\n", cmd)
		}
		return ErrSkip
	}
	fmt.Fprintf(Stdout, "$ %s\n", cmd)
	return nil // system.RemovePackage(l.Name, l.Opts)
}

// Packages task slice
type Packages []*Package

func (packages *Packages) String() string {
	// s := ""
	// for i, l := range *packages {
	// 	s += fmt.Sprintf("%s", l)
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
