package dotfile

// PkgTask struct
type PkgTask struct {
	Name string
	Sudo bool
	Task
}

// Do ...
func (t *PkgTask) Do(a string) (string, error) {
	return do(t, a)
}

// Install copy
func (t *PkgTask) Install() (string, error) {
	sudo = t.Sudo
	return PacInstall(t.Name)
}

// Remove copy
func (t *PkgTask) Remove() (string, error) {
	sudo = t.Sudo
	return PacInstall(t.Name)
}
