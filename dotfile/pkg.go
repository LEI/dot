package dotfile

import (
	"fmt"
)

// PkgTask struct
type PkgTask struct {
	Name string
}

// Do ...
func (t *PkgTask) Do(a string) error {
	return do(t, a)
}

// Install copy
func (t *PkgTask) Install() error {
	fmt.Printf("TODO install package %s\n", t.Name)
	return nil
}

// Remove copy
func (t *PkgTask) Remove() error {
	fmt.Printf("TODO remove package %s\n", t.Name)
	return nil
}
