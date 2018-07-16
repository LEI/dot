package parsers

import (
	"fmt"
	"os"
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

// NewPkg ...
func NewPkg(i interface{}) *Pkg {
	p := &Pkg{}
	if i == nil {
		return p
	}
	if val, ok := i.(string); ok {
		p.Name = val
	} else if val, ok := i.(Pkg); ok {
		*p = val
	} else if val, ok := i.(map[interface{}]interface{}); ok {
		// Get name
		name, ok := val["name"].(string)
		if !ok {
			fmt.Printf("Missing pkg name: %+v\n", val)
			os.Exit(1)
		}
		p.Name = name
		p.OS = *NewSlice(val["os"])
		// p.Action = NewSlice(val["action"])
		action, ok := val["action"].(string)
		if ok {
			p.Action = action
		}
	// } else if val, ok := i.(*Pkg); ok {
	// 	p = val
	// } else if val, ok := i.([]string); ok {
	// 	fmt.Println("MS", val)
	// } else if val, ok := i.([]interface{}); ok {
	// 	// p.OS = *NewSlice(val["os"])
	// 	fmt.Println("IS", val)
	// } else if val, ok := i.(map[string]string); ok {
	// 	fmt.Println("MSS", val, i)
	// } else if val, ok := i.(map[string]interface{}); ok {
	// 	fmt.Println("MSI", val, i)
	// } else if val, ok := i.(interface{}); ok {
	// 	fmt.Println("II", val, i)
	} else {
		fmt.Printf("Unable to assert Pkg: %+v\n", i)
		os.Exit(1)
	}
	return p
}

// Packages ...
type Packages []*Pkg

// UnmarshalYAML ...
func (p *Packages) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var i interface{}
	if err := unmarshal(&i); err != nil {
		return err
	}
	switch val := i.(type) {
	case []string:
		for _, v := range val {
			*p = append(*p, NewPkg(v)) // &Pkg{Name: v}
		}
		break
	case []interface{}:
		for _, v := range val {
			*p = append(*p, NewPkg(v))
		}
		break
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
