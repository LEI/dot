package tasks

import (
	"fmt"
)

// Link task
type Link struct {
	Task
	Source, Target string
	// backup bool
}

func (l *Link) String() string {
	return fmt.Sprintf("%s -> %s", l.Source, l.Target)
}

// Check task
func (l *Link) Check() error {
	return nil
}

// Execute task
func (l *Link) Execute() error {
	fmt.Printf("-> %+v\n", l)
	return nil
}

// Links list
type Links []*Link

func (links *Links) String() string {
	s := ""
	for i, l := range *links {
		s += fmt.Sprintf("%s", l)
		if i > 0 {
			s += "\n"
		}
	}
	return s
}

// Parse data
func (links *Links) Parse(i interface{}) error {
	if i == nil {
		return nil
	}
	switch v := i.(type) {
	case string:
		s, t, err := parseDest(v)
		if err != nil {
			return err
		}
		*links = append(*links, &Link{Source: s, Target: t})
	case []string:
		for _, val := range v {
			s, t, err := parseDest(val)
			if err != nil {
				return err
			}
			*links = append(*links, &Link{Source: s, Target: t})
		}
	case []interface{}:
		for _, val := range v {
			s, t, err := parseDest(val.(string))
			if err != nil {
				return err
			}
			*links = append(*links, &Link{Source: s, Target: t})
		}
	case map[string]string:
		for s, t := range v {
			*links = append(*links, &Link{Source: s, Target: t})
		}
	case map[string]interface{}:
		for s, t := range v {
			*links = append(*links, &Link{Source: s, Target: t.(string)})
		}
	case map[interface{}]interface{}:
		for s, t := range v {
			*links = append(*links, &Link{Source: s.(string), Target: t.(string)})
		}
	default:
		return fmt.Errorf("unable to parse role links: %+v", v)
	}
	return nil
}

// Check links task
func (links *Links) Check() error {
	// TODO cli.Errors
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
