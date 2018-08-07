package tasks

import (
	"fmt"

	"github.com/LEI/dot/system"
)

// Link task
type Link struct {
	Map
	// Task
	Source, Target string
	// backup bool
}

func (l *Link) String() string {
	return fmt.Sprintf("%s -> %s", l.Source, l.Target)
}

// // Parse slice
// func (l *Link) Parse(i interface{}) error {
// 	m, err := NewMap(i)
// 	// *l = *m
// 	return err
// }

// Check task
func (l *Link) Check() error {
	if l.Source == "" {
		return fmt.Errorf("link: empty source")
	}
	// if l.Target == "" {
	// 	return fmt.Errorf("link: missing target")
	// }
	return system.CheckSymlink(l.Source, l.Target)
}

// Execute task
func (l *Link) Execute() error {
	return system.Symlink(l.Source, l.Target)
}

// Links list
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

// Parse links
func (links *Links) Parse(i interface{}) error {
	newLinks := &Links{}
	m, err := NewMap(i)
	if err != nil {
		return err
	}
	for k, v := range *m {
		*newLinks = append(*newLinks, &Link{
			Source: k,
			Target: v,
		})
	}
	*links = *newLinks
	return nil
}

// Check links task
func (links *Links) Check() error {
	// // cli.Errors
	// fmt.Println("link", *links)
	for _, l := range *links {
		if err := l.Check(); err != nil {
			return err
		}
	}
	return nil
}

// Execute links task
func (links *Links) Execute() error {
	for _, l := range *links {
		if err := l.Execute(); err != nil {
			return err
		}
	}
	return nil
}
