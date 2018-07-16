package parsers

import (
	"fmt"
	"reflect"
	// "github.com/LEI/dot/dotfile"
)

// Tpl ...
type Tpl struct {
	Source, Target string
	Data           interface{}
}

// Templates ...
type Templates []*Tpl

// func (p *Templates) String() string {
// 	s := ""
// 	for _, v := range *p {
// 		s+= fmt.Sprintf("%+v", v)
// 	}
// 	return s
// }

// Add ...
func (t *Templates) Add(i interface{}) error {
	tpl := &Tpl{}
	if i == nil {
		return fmt.Errorf("Trying to add nil to tmpls: %+v", t)
	}
	if val, ok := i.(string); ok {
		tpl.Source = val
	} else if val, ok := i.(Tpl); ok {
		*tpl = val
	} else if val, ok := i.(map[interface{}]interface{}); ok {
		// Get source
		src, ok := val["source"].(string)
		if !ok {
			return fmt.Errorf("Missing tpl source: %+v", val)
		}
		tpl.Source = src
		dst, ok := val["target"].(string)
		if ok {
			tpl.Target = dst
		}
		data, ok := val["vars"].(interface{})
		if ok {
			tpl.Data = data
		}
		// } else if val, ok := i.(*Tpl); ok {
		// 	tpl = val
		// } else if val, ok := i.([]string); ok {
		// 	fmt.Println("MS", val)
		// } else if val, ok := i.([]interface{}); ok {
		// 	// tpl.OS = *NewSlice(val["os"])
		// 	fmt.Println("IS", val)
		// } else if val, ok := i.(map[string]string); ok {
		// 	fmt.Println("MSS", val, i)
		// } else if val, ok := i.(map[string]interface{}); ok {
		// 	fmt.Println("MSI", val, i)
		// } else if val, ok := i.(interface{}); ok {
		// 	fmt.Println("II", val, i)
	} else {
		return fmt.Errorf("Unable to assert Tpl: %+v", i)
	}
	*t = append(*t, tpl)
	return nil
}

// UnmarshalYAML ...
func (t *Templates) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var i interface{}
	if err := unmarshal(&i); err != nil {
		return err
	}
	switch val := i.(type) {
	case map[string]string:
		fmt.Println("Unmarshal tpl -> map[string]string", val)
		for k, v := range val {
			if k != "" {
				return fmt.Errorf("Unexpected key: %s", k)
			}
			if err := t.Add(v); err != nil {
				return err
			}
		}
		break
	case []interface{}:
		for _, v := range val {
			if err := t.Add(v); err != nil {
				return err
			}
		}
		break
	default:
		t := reflect.TypeOf(val)
		T := t.Elem()
		if t.Kind() == reflect.Map {
			T = reflect.MapOf(t.Key(), t.Elem())
		}
		return fmt.Errorf("Unable to unmarshal templates (%s) into struct: %+v", T, val)
	}
	return nil
}
