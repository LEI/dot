package tasks

import (
	"fmt"
)

// Map task
type Map map[string]string

// Parse map
func (m *Map) Parse(i interface{}) error {
	mm, err := NewMap(i)
	if err != nil {
		return err
	}
	*m = *mm
	return nil
}

// NewMap parse
func NewMap(i interface{}) (*Map, error) {
	m := &Map{}
	if i == nil {
		return m, nil
	}
	switch v := i.(type) {
	case string:
		s, t, err := parseDest(v)
		if err != nil {
			return m, err
		}
		(*m)[s] = t
	case []string:
		for _, val := range v {
			s, t, err := parseDest(val)
			if err != nil {
				return m, err
			}
			(*m)[s] = t
		}
	case []interface{}:
		for _, val := range v {
			s, t, err := parseDest(val.(string))
			if err != nil {
				return m, err
			}
			(*m)[s] = t
		}
	case map[string]string:
		for s, t := range v {
			(*m)[s] = t
		}
	case map[string]interface{}:
		for s, t := range v {
			(*m)[s] = t.(string)
		}
	case map[interface{}]interface{}:
		for s, t := range v {
			(*m)[s.(string)] = t.(string)
		}
	default:
		return m, fmt.Errorf("unable to parse map: %+v", v)
	}
	return m, nil
}
