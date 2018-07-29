package dotfile

import (
	"fmt"
	"strings"
)

// PkgTask struct
type PkgTask struct {
	Name string
	Opts []string
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
	args := strings.Split(t.Name, " ")
	for _, a := range t.Opts {
		args = append(args, "--" + a)
	}
	return PacInstall(args...)
}

// Remove package
func (t *PkgTask) Remove() (string, error) {
	sudo = t.Sudo
	args := strings.Split(t.Name, " ")
	return PacRemove(args...)
}
