package dotfile

import ()

// PkgTask struct
type PkgTask struct {
	Name string
	Sudo bool
}

// Do ...
func (t *PkgTask) Do(a string) error {
	return do(t, a)
}

// Install copy
func (t *PkgTask) Install() error {
	sudo = t.Sudo
	return PacInstall(t.Name)
}

// Remove copy
func (t *PkgTask) Remove() error {
	sudo = t.Sudo
	return PacInstall(t.Name)
}
