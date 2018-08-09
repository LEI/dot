package tasks

import (
	"fmt"

	"github.com/LEI/dot/cli/config/types"
	"github.com/LEI/dot/system"
)

// Dir task
type Dir struct {
	Task
	Path string
}

func (d *Dir) String() string {
	return fmt.Sprintf("dir[%s]", d.Path)
}

// Check dir task
func (d *Dir) Check() error {
	if d.Path == "" {
		return fmt.Errorf("dir: empty path")
	}
	err := system.CheckDir(d.Path)
	switch err {
	case system.ErrDirExist:
		d.toDo = true
	default:
		return err
	}
	return nil
}

// Install dir task
func (d *Dir) Install() error {
	cmd := fmt.Sprintf("mkdir -p %s", d.Path)
	if !d.DoInstall() {
		if Verbose {
			fmt.Fprintf(Stdout, "# %s\n", cmd)
		}
		return ErrSkip
	}
	fmt.Fprintf(Stdout, "$ %s\n", cmd)
	return system.CreateDir(d.Path)
}

// Remove dir task
func (d *Dir) Remove() error {
	cmd := fmt.Sprintf("rmdir %s", d.Path)
	if !d.DoRemove() {
		if Verbose {
			fmt.Fprintf(Stdout, "# %s\n", cmd)
		}
		return ErrSkip
	}
	empty, err := system.IsEmptyDir(d.Path)
	if err != nil {
		return err
	}
	if !empty {
		// return fmt.Errorf("remove %s: directory not empty")
		if Verbose {
			fmt.Fprintf(Stdout, "# %s\n", cmd)
		}
		return ErrSkip
	}
	fmt.Fprintf(Stdout, "$ %s\n", cmd)
	return system.RemoveDir(d.Path) // RmDirIfEmpty
}

// Dirs task slice
type Dirs []*Dir

// // Value return the underlying slice
// func (dirs *Dirs) Value() []*Dir {
// 	return *dirs // (*dirs).value
// }

func (dirs *Dirs) String() string {
	// s := ""
	// for i, l := range *dirs {
	// 	s += fmt.Sprintf("%s", l)
	// 	if i > 0 {
	// 		s += "\n"
	// 	}
	// }
	// return s
	return fmt.Sprintf("%s", *dirs)
}

// Add a dir
func (dirs *Dirs) Add(d *Dir) {
	// (*dirs).value = append(dirs.Value(), d)
	*dirs = append(*dirs, d)
}

// Parse dir tasks
func (dirs *Dirs) Parse(i interface{}) error {
	dd := &Dirs{}
	s, err := types.NewSlice(i)
	if err != nil {
		return err
	}
	for _, v := range *s {
		d := &Dir{Path: v}
		// *dd = append(*dd, d)
		// (*dd).value = append((*dd).value, d)
		dd.Add(d)
	}
	*dirs = *dd
	return nil
}
