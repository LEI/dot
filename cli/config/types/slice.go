package types

// https://github.com/golang/go/wiki/SliceTricks

import (
	"fmt"
)

// Slice task
type Slice []string

func (s *Slice) String() string {
	return fmt.Sprintf("%s", *s)
}

// Value slice
func (s *Slice) Value() []string {
	return *s
}

// Parse slice
func (s *Slice) Parse(i interface{}) error {
	ss, err := NewSlice(i)
	*s = *ss
	return err
}

// NewSlice ...
func NewSlice(i interface{}) (*Slice, error) {
	s := &Slice{}
	if i == nil {
		return s, nil
	}
	switch v := i.(type) {
	case string:
		*s = append(*s, v)
	case []string:
		*s = append(*s, v...)
	case []interface{}:
		for _, val := range v {
			*s = append(*s, val.(string))
		}
	default:
		return s, fmt.Errorf("unable to parse slice: %+v", v)
	}
	return s, nil
}

// SliceMap task
type SliceMap []interface{}

func (s *SliceMap) String() string {
	return fmt.Sprintf("%s", *s)
}

// Value slice map
func (s *SliceMap) Value() []interface{} {
	return *s
}

// // Value slice map
// func (s *SliceMap) Value() []string {
// 	ss := []string{}
// 	for _, m := range *s {
// 		ss = append(ss, m.(string))
// 	}
// 	return ss // *s
// }

// Parse slice map
func (s *SliceMap) Parse(i interface{}) error {
	ss, err := NewSliceMap(i)
	*s = *ss
	return err
}

// NewSliceMap ...
func NewSliceMap(i interface{}) (*SliceMap, error) {
	sm := &SliceMap{}
	if i == nil {
		return sm, nil
	}
	switch v := i.(type) {
	case string:
		*sm = append(*sm, v)
	case []string:
		// *sm = append(*sm, v...)
		for _, s := range v {
			*sm = append(*sm, s)
		}
	case []interface{}:
		for _, val := range v {
			switch m := val.(type) {
			case string:
				*sm = append(*sm, m)
			case interface{}:
				m, err := NewMap(m)
				if err != nil {
					return sm, err
				}
				*sm = append(*sm, m)
			default:
				return sm, fmt.Errorf("unable to parse slice map interface: %+v", v)
			}
		}
	default:
		return sm, fmt.Errorf("unable to parse slice map: %+v", v)
	}
	return sm, nil
}

// // NewSlice ...
// func NewSlice(i interface{}) ([]string, error) {
// 	s := []string{}
// 	if i == nil {
// 		return s, nil
// 	}
// 	switch v := i.(type) {
// 	case string:
// 		s = append(s, v)
// 	case []string:
// 		s = append(s, v...)
// 	case []interface{}:
// 		for _, val := range v {
// 			s = append(s, val.(string))
// 		}
// 	default:
// 		return s, fmt.Errorf("unable to parse slice: %+v", v)
// 	}
// 	return s, nil
// }
