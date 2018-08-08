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
	return fmt.Sprintf("link %s -> %s", l.Source, l.Target)
}

// Check link task
func (l *Link) Check() error {
	if l.Source == "" {
		return fmt.Errorf("link: empty source")
	}
	// if l.Target == "" {
	// 	return fmt.Errorf("link: missing target")
	// }
	// fmt.Printf("Checking %+v\n", l)
	err := system.CheckSymlink(l.Source, l.Target)
	switch err {
	case nil:
	case system.ErrLinkExist:
		l.toDo = true
	default:
		return err
	}
	return nil
}

// Install link task
func (l *Link) Install() error {
	if !l.DoInstall() {
		return ErrSkip
	}
	return system.Symlink(l.Source, l.Target)
}

// Remove link task
func (l *Link) Remove() error {
	if !l.DoRemove() {
		return ErrSkip
	}
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
	m, err := types.NewMap(i)
	if err != nil {
		return err
	}
	for k, v := range *m {
		*ll = append(*ll, &Link{
			Source: k,
			Target: v,
		})
	}
	*links = *ll
	return nil
}

// Check link tasks
func (links *Links) Check() error {
	// // cli.Errors
	// fmt.Println("link", *links)
	// if *links == nil {
	// 	return nil
	// }
	for _, l := range *links {
		if err := l.Check(); err != nil {
			return err
		}
	}
	return nil
}

// Install link tasks
func (links *Links) Install() error {
	for _, l := range *links {
		if err := l.Install(); err != nil && err != ErrSkip {
			return err
		}
	}
	return nil
}

// Remove link tasks
func (links *Links) Remove() error {
	for _, l := range *links {
		if err := l.Remove(); err != nil && err != ErrSkip {
			return err
		}
	}
	return nil
}
