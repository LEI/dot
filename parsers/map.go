package parsers

import (
	"fmt"
	"reflect"
)

// Map ...
type Map map[string]string // map[interface{}]interface{}

// Add ...
func (m *Map) Add(key, val interface{}) error {
	k := key.(string)
	v := val.(string)
	(*m)[k] = v
	return nil
}

// UnmarshalYAML ...
func (m *Map) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if *m == nil {
		*m = make(Map)
	}
	var i interface{}
	if err := unmarshal(&i); err != nil {
		return err
	}
	switch val := i.(type) {
	case string:
		(*m)[val] = ""
	case []string:
		for _, v := range val {
			(*m)[v] = ""
		}
	case map[string]string:
		// for _, v := range val {
		// 	(*m)[v] = ""
		// }
		fmt.Println("TODO m[s]s ->", val)
	// case map[string]interface{}:
	case map[interface{}]interface{}:
		fmt.Println("TODO m[i{}]i{} ->", val)
		// for k, v := range val {
		// 	fmt.Println("k",k,"v",v)
		// 	switch w := v.(type) {
		// 	default:
		// 		t := reflect.TypeOf(w)
		// 		T := t.Elem()
		// 		if t.Kind() == reflect.Map {
		// 			T = reflect.MapOf(t.Key(), t.Elem())
		// 		}
		// 		return fmt.Errorf("unable to unmarshal %s into Map element: %+v", T, w)
		// 	}
		// }
	case interface{}:
		switch v := val.(type) {
		case map[interface{}]interface{}:
			fmt.Println("TODO i{} -> map[interface{}]interface{}", v)
		case []interface{}:
			for _, w := range v {
				switch x := w.(type) {
				case string:
					(*m)[x] = "" // w.(interface{})
				// case interface{}:
				default:
					return fmt.Errorf("unable to unmarshal: expected string, found %+v", w)
				}
			}
		// case interface{}:
		// 	fmt.Println("i{} -> interface{}", v)
		default:
			t := reflect.TypeOf(v)
			T := t.Elem()
			if t.Kind() == reflect.Map {
				T = reflect.MapOf(t.Key(), t.Elem())
			}
			return fmt.Errorf("unable to unmarshal %s into Map element: %+v", T, v)
		}
	default:
		t := reflect.TypeOf(val)
		T := t.Elem()
		if t.Kind() == reflect.Map {
			T = reflect.MapOf(t.Key(), t.Elem())
		}
		return fmt.Errorf("unable to unmarshal %s into Map: %+v", T, val)
	}
	return nil
}

// Value ...
func (m *Map) Value() map[string]string {
	r := map[string]string{}
	for k, v := range *m {
		r[k] = v // .(string)
	}
	return r
}
