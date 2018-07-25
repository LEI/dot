package parsers

import (
	"fmt"
	// "os"
	"reflect"
)

// Slice ...
type Slice []string

// NewSlice ...
func NewSlice(i interface{}) (*Slice, error) {
	s := &Slice{}
	if i == nil {
		return s, nil
	}
	if val, ok := i.(string); ok {
		*s = append(*s, val)
	} else if val, ok := i.([]string); ok {
		for _, v := range val {
			*s = append(*s, v)
		}
	} else if val, ok := i.([]interface{}); ok {
		for _, v := range val {
			if item, ok := v.(string); ok {
				*s = append(*s, item)
			} else if v != nil {
				t := reflect.TypeOf(v)
				// T := t.Elem()
				// if t.Kind() == reflect.Map {
				// 	T = reflect.MapOf(t.Key(), t.Elem())
				// }
				return s, fmt.Errorf("unable to unmarshal %s into Slice element: %+v", t, v)
			} else {
				return s, fmt.Errorf("unable to assert Slice element: %+v", v)
			}
		}
	} else {
		return s, fmt.Errorf("unable to assert Slice: %+v", i)
	}
	// switch val := i.(type) {
	// case nil:
	// case interface{}:
	// 	v, ok := val.(Slice)
	// 	if !ok {
	// 		fmt.Printf("Unable to assert Item: %+v\n", val)
	// 		os.Exit(1)
	// 	}
	// 	fmt.Println("=======", v)
	// 	// for _, v := range v {
	// 	// 	vi, ok = v.(*Item)
	// 	// 	if !ok {
	// 	// 		fmt.Printf("Unable to assert []string: %+v\n", val)
	// 	// 		os.Exit(1)
	// 	// 	}
	// 	// 	*s = append(*s, vi)
	// 	// }
	// 	// v, ok := val.([]interface{})
	// 	// if !ok {
	// 	// 	fmt.Printf("Unable to assert []string: %+v\n", val)
	// 	// 	os.Exit(1)
	// 	// }
	// 	// for _, v := range v {
	// 	// 	*s = append(*s, v.(*Item))
	// 	// }
	// case []string:
	// 	for _, v := range val {
	// 		fmt.Println("NewSlice 2", v)
	// 	}
	// default:
	// 	// t := reflect.TypeOf(val)
	// 	// T := t.Elem()
	// 	// if t.Kind() == reflect.Map {
	// 	// 	T = reflect.MapOf(t.Key(), t.Elem())
	// 	// }
	// 	// fmt.Printf("Unable to unmarshal %s into Slice: %+v\n", T, val)
	// 	fmt.Printf("Unable to unmarshal ???: %+v\n", val)
	// 	os.Exit(1)
	// }
	return s, nil
}

// Value return a new iterator
func (s *Slice) Value() []string {
	r := []string{}
	for _, v := range *s {
		r = append(r, v)
	}
	return r
}

// UnmarshalYAML ...
func (s *Slice) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if *s == nil {
		*s = make(Slice, 0)
	}
	var i interface{} // *Item
	if err := unmarshal(&i); err != nil {
		return err
	}
	fmt.Println("NewSlice UnmarshalYAML", i)
	s, err := NewSlice(i)
	if err != nil {
		return err
	}
	/*
	switch val := i.(type) {
	case []string:
		// 	*s = append(*s, &item)
		for _, v := range val {
			// item := Item(v)
			fmt.Println("=", v)
			// *s = append(*s, &item)
		}
	// case map[string]interface{}:
	default:
		t := reflect.TypeOf(val)
		T := t.Elem()
		if t.Kind() == reflect.Map {
			T = reflect.MapOf(t.Key(), t.Elem())
		}
		return fmt.Errorf("unable to unmarshal %s into struct: %+v", T, val)
	}
	*/
	return nil
}
