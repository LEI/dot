package dotfile

import (
	"fmt"
)

// PkgTask struct
type PkgTask struct {
	Name string
	Sudo bool
	Task
}

// Status package
func (t *PkgTask) Status() bool {
	return true
}

// Do package task
func (t *PkgTask) Do(a string) (string, error) {
	return do(t, a)
}

// List package
func (t *PkgTask) List() (string, error) {
	str := fmt.Sprintf("Pkg: %s", t.Name)
	// t.Sudo
	return str, nil
}

// Install package
func (t *PkgTask) Install() (string, error) {
	sudo = t.Sudo
	return PacInstall(t.Name)
}

// Remove package
func (t *PkgTask) Remove() (string, error) {
	sudo = t.Sudo
	return PacInstall(t.Name)
}
