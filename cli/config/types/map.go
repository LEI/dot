package types

import (
	"fmt"
	"path/filepath"
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
		s, t, err := parsePaths(v)
		if err != nil {
			return m, err
		}
		(*m)[s] = t
	case []string:
		for _, val := range v {
			s, t, err := parsePaths(val)
			if err != nil {
				return m, err
			}
			(*m)[s] = t
		}
	case []interface{}:
		for _, val := range v {
			s, t, err := parsePaths(val.(string))
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
			val, _ := t.(string)
			(*m)[s] = val
			// if t != nil {
			// 	str = t.(string)
			// } else {
			// 	str = ""
			// }
			// (*m)[s] = str
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

// // NewMap parse
// func NewMap(i interface{}) (map[string]string, error) {
// 	m := map[string]string{}
// 	if i == nil {
// 		return m, nil
// 	}
// 	switch v := i.(type) {
// 	case string:
// 		s, t, err := parsePaths(v)
// 		if err != nil {
// 			return m, err
// 		}
// 		m[s] = t
// 	case []string:
// 		for _, val := range v {
// 			s, t, err := parsePaths(val)
// 			if err != nil {
// 				return m, err
// 			}
// 			m[s] = t
// 		}
// 	case []interface{}:
// 		for _, val := range v {
// 			s, t, err := parsePaths(val.(string))
// 			if err != nil {
// 				return m, err
// 			}
// 			m[s] = t
// 		}
// 	case map[string]string:
// 		for s, t := range v {
// 			m[s] = t
// 		}
// 	case map[string]interface{}:
// 		for s, t := range v {
// 			m[s] = t.(string)
// 		}
// 	case map[interface{}]interface{}:
// 		for s, t := range v {
// 			m[s.(string)] = t.(string)
// 		}
// 	default:
// 		return m, fmt.Errorf("unable to parse map: %+v", v)
// 	}
// 	return m, nil
// }

// Parse src:dst paths
func parsePaths(p string) (src, dst string, err error) {
	parts := filepath.SplitList(p)
	switch len(parts) {
	case 1:
		src = p
	case 2:
		src = parts[0]
		dst = parts[1]
	default:
		return src, dst, fmt.Errorf("unhandled path spec: %s", src)
	}
	return src, dst, nil
	// src = s
	// if strings.Contains(s, ":") {
	// 	parts := strings.Split(s, ":")
	// 	if len(parts) != 2 {
	// 		return src, dst, fmt.Errorf("unable to parse dest spec: %s", s)
	// 	}
	// 	src = parts[0]
	// 	dst = parts[1]
	// }
	// // if dst == "" && isDir(src) {
	// // 	dst = PathHead(src)
	// // }
	// return src, dst, nil
}