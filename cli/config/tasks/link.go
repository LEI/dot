package tasks

import (
	"fmt"
)

// Link task
type Link struct {
	Source, Target string
	// backup bool
}

func (l *Link) String() string {
	return fmt.Sprintf("%s -> %s", l.Source, l.Target)
}

// Links list
type Links []*Link

// Parse data
func (l *Links) Parse(i interface{}) error {
	if i == nil {
		return nil
	}
	switch v := i.(type) {
	case string:
		s, t, err := parseDest(v)
		if err != nil {
			return err
		}
		*l = append(*l, &Link{Source: s, Target: t})
	case []string:
		for _, val := range v {
			s, t, err := parseDest(val)
			if err != nil {
				return err
			}
			*l = append(*l, &Link{Source: s, Target: t})
		}
	case []interface{}:
		for _, val := range v {
			s, t, err := parseDest(val.(string))
			if err != nil {
				return err
			}
			*l = append(*l, &Link{Source: s, Target: t})
		}
	case map[string]string:
		for s, t := range v {
			*l = append(*l, &Link{Source: s, Target: t})
		}
	case map[string]interface{}:
		for s, t := range v {
			*l = append(*l, &Link{Source: s, Target: t.(string)})
		}
	case map[interface{}]interface{}:
		for s, t := range v {
			*l = append(*l, &Link{Source: s.(string), Target: t.(string)})
		}
	default:
		return fmt.Errorf("unable to parse role links: %+v", v)
	}
	return nil
}
