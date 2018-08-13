package tasks

import (
	"fmt"

	"github.com/LEI/dot/cli/config/types"
	"github.com/LEI/dot/system"
)

// Copy task
type Copy struct {
	Task
	Source, Target string
}

func (c *Copy) String() string {
	return fmt.Sprintf("copy[%s:%s]", c.Source, c.Target)
}

// Check copy task
func (c *Copy) Check() error {
	if c.Source == "" {
		return fmt.Errorf("copy: empty source")
	}
	err := system.CheckCopy(c.Source, c.Target)
	switch err {
	case system.ErrFileAlreadyExist:
		c.Done()
	default:
		return err
	}
	return nil
}

// Install copy task
func (c *Copy) Install() error {
	cmd := fmt.Sprintf("cp %s %s", c.Source, c.Target)
	if !c.ShouldInstall() {
		if Verbose > 0 {
			fmt.Fprintf(Stdout, "# %s\n", cmd)
		}
		return ErrSkip
	}
	fmt.Fprintf(Stdout, "$ %s\n", cmd)
	return system.Copy(c.Source, c.Target)
}

// Remove copy task
func (c *Copy) Remove() error {
	cmd := fmt.Sprintf("rm %s", c.Target)
	if !c.ShouldRemove() {
		if Verbose > 0 {
			fmt.Fprintf(Stdout, "# %s\n", cmd)
		}
		return ErrSkip
	}
	fmt.Fprintf(Stdout, "$ %s\n", cmd)
	return system.Remove(c.Target)
}

// Files task slice
type Files []*Copy

func (files *Files) String() string {
	// s := ""
	// for i, c := range *files {
	// 	s += fmt.Sprintf("%s", c)
	// 	if i > 0 {
	// 		s += "\n"
	// 	}
	// }
	// return s
	return fmt.Sprintf("%s", *files)
}

// Parse copy tasks
func (files *Files) Parse(i interface{}) error {
	cc := &Files{}
	m, err := types.NewMapPaths(i)
	if err != nil {
		return err
	}
	for k, v := range *m {
		c := &Copy{
			Source: k,
			Target: v,
		}
		// *cc = append(*cc, c)
		cc.Add(*c)
	}
	*files = *cc
	return nil
}

// Add a dir
func (files *Files) Add(c Copy) {
	*files = append(*files, &c)
}
