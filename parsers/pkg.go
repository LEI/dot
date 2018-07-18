package parsers

import (
	"fmt"
	"reflect"
)

// Pkg ...
type Pkg struct {
	Name   string
	OS     Slice
	Action string // install, remove
}

// func (p *Pkg) String() string {
// 	return fmt.Sprintf("%s %s%+v", p.Name, p.Action, p.OS)
// }

// Packages ...
type Packages []*Pkg

// Add ...
func (p *Packages) Add(i interface{}) error {
	pkg := &Pkg{}
	if i == nil {
		return fmt.Errorf("Trying to add nil to pkgs: %+v", p)
	}
	if val, ok := i.(string); ok {
		pkg.Name = val
	} else if val, ok := i.(Pkg); ok {
		*pkg = val
	} else if val, ok := i.(map[interface{}]interface{}); ok {
		// Get name
		name, ok := val["name"].(string)
		if !ok {
			return fmt.Errorf("Missing pkg name: %+v", val)
		}
		pkg.Name = name
		pkg.OS = *NewSlice(val["os"])
		// pkg.Action = NewSlice(val["action"])
		action, ok := val["action"].(string)
		if ok {
			pkg.Action = action
		}
		// } else if val, ok := i.(*Pkg); ok {
		// 	pkg = val
		// } else if val, ok := i.([]string); ok {
		// 	fmt.Println("MS", val)
		// } else if val, ok := i.([]interface{}); ok {
		// 	// pkg.OS = *NewSlice(val["os"])
		// 	fmt.Println("IS", val)
		// } else if val, ok := i.(map[string]string); ok {
		// 	fmt.Println("MSS", val, i)
		// } else if val, ok := i.(map[string]interface{}); ok {
		// 	fmt.Println("MSI", val, i)
		// } else if val, ok := i.(interface{}); ok {
		// 	fmt.Println("II", val, i)
	} else {
		return fmt.Errorf("Unable to assert Pkg: %+v", i)
	}
	*p = append(*p, pkg)
	return nil
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
			if err := p.Add(v); err != nil {
				return err
			}
		}
	case []interface{}:
		for _, v := range val {
			if err := p.Add(v); err != nil {
				return err
			}
		}
	default:
		t := reflect.TypeOf(val)
		T := t.Elem()
		if t.Kind() == reflect.Map {
			T = reflect.MapOf(t.Key(), t.Elem())
		}
		return fmt.Errorf("Unable to unmarshal packages (%s) into struct: %+v", T, val)
	}
	return nil
}
