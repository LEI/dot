package parsers

import (
	"fmt"
	"reflect"
)

// Pkg ...
type Pkg struct {
	Name   string
	Args   []string
	OS     Slice
	Action string // install, remove
}

// func (p *Pkg) String() string {
// 	return fmt.Sprintf("%s %s%+v", p.Name, p.Action, p.OS)
// }

// NewPkg ...
func NewPkg(i interface{}) (*Pkg, error) {
	pkg := &Pkg{}
	if i == nil {
		return pkg, fmt.Errorf("trying to add nil pkg: %+v", i)
	}
	if val, ok := i.(string); ok {
		pkg.Name = val
	} else if val, ok := i.(Pkg); ok {
		*pkg = val
	} else if val, ok := i.(map[interface{}]interface{}); ok {
		// Get name
		pkgName, ok := val["name"].(string)
		if !ok {
			return pkg, fmt.Errorf("missing pkg name: %+v", val)
		}
		pkg.Name = pkgName
		pkgArgs, ok := val["args"].([]string)
		if ok {
			pkg.Args = pkgArgs
		}
		pkgOS, err := NewSlice(val["os"])
		if err != nil {
			return pkg, err
		}
		pkg.OS = *pkgOS
		// pkg.Action = NewSlice(val["action"])
		pkgAction, ok := val["action"].(string)
		if ok {
			pkg.Action = pkgAction
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
		return pkg, fmt.Errorf("unable to assert Pkg: %+v", i)
	}
	return pkg, nil
}

// Packages ...
type Packages []*Pkg

// Add ...
func (p *Packages) Add(i interface{}) error {
	pkg, err := NewPkg(i)
	if err != nil {
		return err
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
		return fmt.Errorf("unable to unmarshal packages (%s) into struct: %+v", T, val)
	}
	return nil
}
