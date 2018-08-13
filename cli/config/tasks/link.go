package tasks

import (
	"fmt"

	"github.com/LEI/dot/cli/config/types"
	"github.com/LEI/dot/system"
)

// Link task
type Link struct {
	Task
	Source, Target string
	// backup bool
	// overwrite bool
}

func (l *Link) String() string {
	return fmt.Sprintf("link[%s:%s]", l.Source, l.Target)
}

// Check link task
func (l *Link) Check() error {
	if l.Source == "" {
		return fmt.Errorf("link: empty source")
	}
	// if l.Target == "" {
	// 	return fmt.Errorf("link: missing target")
	// }
	err := system.CheckSymlink(l.Source, l.Target)
	switch err {
	case system.ErrLinkAlreadyExist:
		l.Done()
	default:
		return err
	}
	return nil
}

// Install link task
func (l *Link) Install() error {
	cmd := fmt.Sprintf("ln -s %s %s", l.Source, l.Target)
	if !l.ShouldInstall() {
		if Verbose > 0 {
			fmt.Fprintf(Stdout, "# %s\n", cmd)
		}
		return ErrSkip
	}
	fmt.Fprintf(Stdout, "$ %s\n", cmd)
	return system.Symlink(l.Source, l.Target)
}

// Remove link task
func (l *Link) Remove() error {
	cmd := fmt.Sprintf("rm %s", l.Target)
	if !l.ShouldRemove() {
		if Verbose > 0 {
			fmt.Fprintf(Stdout, "# %s\n", cmd)
		}
		return ErrSkip
	}
	fmt.Fprintf(Stdout, "$ %s\n", cmd)
	return system.Unlink(l.Target)
}

// Links task slice
type Links []*Link

func (links *Links) String() string {
	// s := ""
	// for i, l := range *links {
	// 	s += fmt.Sprintf("%s", l)
	// 	if i > 0 {
	// 		s += "\n"
	// 	}
	// }
	// return s
	return fmt.Sprintf("%s", *links)
}

// Parse link tasks
func (links *Links) Parse(i interface{}) error {
	ll := &Links{}
	m, err := types.NewMapPaths(i)
	if err != nil {
		return err
	}
	for k, v := range *m {
		l := &Link{
			Source: k,
			Target: v,
		}
		// *ll = append(*ll, l)
		ll.Add(*l)
	}
	*links = *ll
	return nil
}

// Add a dir
func (links *Links) Add(l Link) {
	*links = append(*links, &l)
}
