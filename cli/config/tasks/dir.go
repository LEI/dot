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
	return fmt.Sprintf("dir %s", d.Path)
}

// Check dir task
func (d *Dir) Check() error {
	if d.Path == "" {
		return fmt.Errorf("dir: empty path")
	}
	err := system.CheckDirectory(d.Path)
	switch err {
	case nil:
	case system.ErrDirExist:
		d.toDo = true
	default:
		return err
	}
	return nil
}

// Install dir task
func (d *Dir) Install() error {
	if !d.DoInstall() {
		return ErrSkip
	}
	return system.CreateDirectory(d.Path)
}

// Remove dir task
func (d *Dir) Remove() error {
	if !d.DoRemove() {
		return ErrSkip
	}
	return system.RemoveDirectory(d.Path)
}

// Dirs task slice
type Dirs []*Dir

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

// Parse dir tasks
func (dirs *Dirs) Parse(i interface{}) error {
	dd := &Dirs{}
	s, err := types.NewSlice(i)
	if err != nil {
		return err
	}
	for _, v := range *s {
		d := &Dir{Path: v}
		*dd = append(*dd, d)
	}
	*dirs = *dd
	return nil
}

// Check dir tasks
func (dirs *Dirs) Check() error {
	for _, d := range *dirs {
		if err := d.Check(); err != nil {
			return err
		}
	}
	return nil
}

// Install dir tasks
func (dirs *Dirs) Install() error {
	for _, d := range *dirs {
		if err := d.Install(); err != nil && err != ErrSkip {
			return err
		}
	}
	return nil
}

// Remove dir tasks
func (dirs *Dirs) Remove() error {
	for _, d := range *dirs {
		if err := d.Remove(); err != nil && err != ErrSkip {
			return err
		}
	}
	return nil
}
