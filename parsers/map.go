package parsers

import (
	"fmt"
	"reflect"
)

// Map ...
type Map map[string]string // *MapItem // map[interface{}]interface{}

// // MapItem ...
// type MapItem struct {
// 	Name string
// 	Item
// }

// func (i *MapItem) String() string {
// 	return fmt.Sprintf("MapItem{%+v}", i.Name)
// }

// NewMap ...
func NewMap(i interface{}) *Map {
	m := &Map{}
	if i != nil {
		fmt.Println("TODO NewMap", i)
	}
	return m
}

// Len ...
// func (m *Map) Len() int {
// 	return len(*m)
// }

// Add ...
func (m *Map) Add(key, val interface{}) {
	k := key.(string)
	v := val.(string)
	// s, ok := val.(string)
	// var v *MapItem
	// if ok {
	// 	v = &MapItem{Name: s}
	// } else {
	// 	v = val.(*MapItem)
	// 	// TODO log.Fatal
	// }
	(*m)[k] = v
}

// // Merge ...
// func (m *Map) Merge(in map[string]string) *Map {
// 	(*m) = *(&Map{})
// 	for k, v := range in {
// 		// fmt.Println("MERGE k", k, "v", v)
// 		(*m)[k] = &MapItem{v, nil}
// 	}
// 	return m
// }

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
		// 		return fmt.Errorf("Unable to unmarshal %s into MapItem: %+v", T, w)
		// 	}
		// }
	case interface{}:
		switch v := val.(type) {
		case map[interface{}]interface{}:
			fmt.Println("TODO i{} -> map[interface{}]interface{}", v)
			break
		case []interface{}:
			for _, w := range v {
				switch x := w.(type) {
				case string:
					(*m)[x] = "" // w.(interface{})
					break
				// case interface{}:
				default:
					return fmt.Errorf("??? %+v", w)
				}
			}
			break
		// case interface{}:
		// 	fmt.Println("i{} -> interface{}", v)
		// 	break
		default:
			t := reflect.TypeOf(v)
			T := t.Elem()
			if t.Kind() == reflect.Map {
				T = reflect.MapOf(t.Key(), t.Elem())
			}
			return fmt.Errorf("Unable to unmarshal %s into MapItem: %+v", T, v)
		}
	default:
		t := reflect.TypeOf(val)
		T := t.Elem()
		if t.Kind() == reflect.Map {
			T = reflect.MapOf(t.Key(), t.Elem())
		}
		return fmt.Errorf("Unable to unmarshal %s into Map: %+v", T, val)
	}
	return nil
}

// GetStringMapString ...
func (m *Map) GetStringMapString() map[string]string {
	r := map[string]string{}
	for k, v := range *m {
		r[k] = v // .(string)
	}
	return r
}
