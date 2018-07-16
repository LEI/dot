package parsers

import (
	"os"
	"fmt"
	"reflect"
)

// // Item ...
// type Item interface {
// 	// New (string) *Item
// 	String() string
// }

// // SliceItem ...
// type SliceItem struct {
// 	Name string
// 	Item
// }

// func (i *SliceItem) String() string {
// 	return fmt.Sprintf("SliceItem{%+v}", i.Name)
// }

// Slice ...
type Slice []string

// NewSlice ...
func NewSlice(i interface{}) *Slice {
	s := &Slice{}
	if i == nil {
		return s
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
				fmt.Printf("Unable to unmarshal %s into SliceItem: %+v\n", t, v)
			} else {
				fmt.Printf("Unable to assert SliceItem: %+v\n", v)
				os.Exit(1)
			}
		}
	} else {
		fmt.Printf("Unable to assert Slice: %+v\n", i)
		os.Exit(1)
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
	return s
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
	switch val := i.(type) {
	case []string:
		// 	*s = append(*s, &item)
		for _, v := range val {
			// item := Item(v)
			fmt.Println("=", v)
			// *s = append(*s, &item)
		}
		break
	// case map[string]interface{}:
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
