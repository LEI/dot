package cmd

import (
	"fmt"
	"os"
	"reflect"
)

// Paths ...
type Paths map[string]string

// Add ...
func (p *Paths) Add(s string) error {
	src, dst := SplitPath(s)
	// if p.Dir != "" {
	// 	src = filepath.Join(p.Dir, src)
	// }
	// if p.Dst != "" {
	// 	dst = filepath.Join(p.Dst, dst)
	// }
	src = os.ExpandEnv(src)
	dst = os.ExpandEnv(dst)
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
		break
	// case interface{}:
	// 	s := val.(string)
	// 	(*p)[s] = s
	case []interface{}:
		for _, v := range val {
			// (*p)[v.(string)] = v.(string)
			p.Add(v.(string))
		}
		break
	case map[string]string:
		// p = i.(*Paths)
		for k, v := range val {
			if k != "" {
				fmt.Printf("Unmarshal: ignore key '%s'\n", k)
			}
			p.Add(v)
		}
		break
	case map[interface{}]interface{}:
		for k, v := range val {
			if k.(string) != "" {
				fmt.Printf("Unmarshal: ignore key '%s'\n", k.(string))
			}
			// (*p)[v.(string)] = v.(string)
			p.Add(v.(string))
		}
		break
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
