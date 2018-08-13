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