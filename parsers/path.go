package parsers

import (
	"fmt"
	// "os"
	// "path/filepath"
	"reflect"
	// "github.com/LEI/dot/dotfile"
)

// Paths ...
type Paths map[string]string

func (p *Paths) String() string {
	s := ""
	for k, v := range *p {
		s += fmt.Sprintf("%s:%s", k, v)
	}
	return s
}

// Add ...
func (p *Paths) Add(i interface{}) error {
	if i == nil {
		return fmt.Errorf("Trying to add nil to paths: %+v", p)
	}
	var src, dst string
	if val, ok := i.(string); ok {
		src = val
	} else if val, ok := i.(struct{ Source, Target string }); ok {
		src = val.Source
		dst = val.Target
	} else if val, ok := i.(map[interface{}]interface{}); ok {
		// Get name
		s, ok := val["source"].(string)
		if !ok {
			return fmt.Errorf("Missing path source: %+v", val)
		}
		src = s
		t, _ := val["target"].(string)
		if ok {
			dst = t
		}
	}
	// src, dst := dotfile.SplitPath(s)
	// // if p.Dir != "" {
	// // 	src = filepath.Join(p.Dir, src)
	// // }
	// // if p.Dst != "" {
	// // 	dst = filepath.Join(p.Dst, dst)
	// // }
	// src = dotfile.ExpandEnv(src)
	// dst = dotfile.ExpandEnv(dst)
	(*p)[src] = dst
	return nil
}

// func addPaths (p *Paths, v string, target string) error {
// 	// v = filepath.Join(source, v)
// 	paths, err := ParsePath(v, target)
// 	if err != nil {
// 		return err
// 	}
// 	for s, t := range paths {
// 		(*p)[s] = t
// 	}
// 	return nil
// }

// UnmarshalYAML ...
func (p *Paths) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Avoid assignment to entry in nil map
	// FIXME: invalid memory address or nil pointer dereference
	if *p == nil {
		*p = make(Paths)
	}
	var i interface{}
	if err := unmarshal(&i); err != nil {
		return err
	}
	// paths := map[string]string{}
	switch val := i.(type) {
	case []string:
		for _, v := range val {
			// (*p)[v] = v
			p.Add(v)
		}
	// case interface{}:
	// 	s := val.(string)
	// 	(*p)[s] = s
	case []interface{}:
		for _, v := range val {
			p.Add(v)
		}
	case map[string]string:
		// p = i.(*Paths)
		for k, v := range val {
			if k != "" {
				fmt.Printf("Unmarshal: ignore key '%s'\n", k)
			}
			p.Add(v)
		}
	case map[interface{}]interface{}:
		for k, v := range val {
			if k.(string) != "" {
				fmt.Printf("Unmarshal: ignore key '%s'\n", k.(string))
			}
			// (*p)[v.(string)] = v.(string)
			p.Add(v.(string))
		}
	default:
		t := reflect.TypeOf(val)
		T := t.Elem()
		if t.Kind() == reflect.Map {
			T = reflect.MapOf(t.Key(), t.Elem())
		}
		return fmt.Errorf("Unable to unmarshal %s into struct: %+v", T, val)
	}
	return nil
}
