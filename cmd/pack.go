package cmd

import (
	"fmt"
	// "os"
	"reflect"
)

// Packages ...
type Packages []*Pkg

// Pkg ...
type Pkg struct {
	Name string
	OS   []string
}

// UnmarshalYAML ...
func (p *Packages) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var i interface{}
	if err := unmarshal(&i); err != nil {
		return err
	}
	switch val := i.(type) {
	case []string:
		for _, v := range val {
			*p = append(*p, &Pkg{Name: v})
		}
		break
	case []interface{}:
		for _, v := range val {
			pkg := &Pkg{}
			switch V := v.(type) {
			case string:
				pkg.Name = V
				break
			case interface{}:
				switch m := V.(type) {
				case map[interface{}]interface{}:
					name, ok := m["name"].(string)
					if !ok {
						return fmt.Errorf("Missing pkg name: %+v", m)
					}
					pkg.Name = name
					switch n := m["os"].(type) {
					case []string:
						pkg.OS = n
						break
					case interface{}:
						// n = n.(interface{})
						// // FIXME: interface {} is []interface {}, not []string
						for _, o := range n.([]interface{}) {
							pkg.OS = append(pkg.OS, o.(string))
						}
						break
					case []interface{}:
						for _, o := range n {
							pkg.OS = append(pkg.OS, o.(string))
						}
						break
					default:
						t := reflect.TypeOf(val)
						T := t.Elem()
						if t.Kind() == reflect.Map {
							T = reflect.MapOf(t.Key(), t.Elem())
						}
						return fmt.Errorf("Unable to unmarshal pkg os: %+v (%s)", m, T)
					}
					// m, ok := w["os"].([]string)
					// if ok {
					// 	pkg.OS = m
					// } else {
					// 	if n, ok := w["os"].([]interface{}); ok {
					// 		for _, o := range n {
					// 			pkg.OS = append(pkg.OS, o.(string))
					// 		}
					// 	} else if w != nil {
					// 		return fmt.Errorf("Invalid pkg os list: %+v", w)
					// 	}
					// }
					break
				default:
					return fmt.Errorf("Unable to unmarshal pkg: %+v", m)
				}
				break
			// case map[string]string:
			// 	fmt.Println("s!!!!!!!!!!", V)
			// 	break
			// case map[interface{}]interface{}:
			// 	fmt.Println("i!!!!!!!!!!", V)
			// 	break
			default:
				return fmt.Errorf("Unable to unmarshal into struct: %+v", V)
			}
			*p = append(*p, pkg)
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
