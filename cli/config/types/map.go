package types

import (
	"fmt"
	"path/filepath"
)

// Map task
type Map map[string]interface{}

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
func NewMap(i interface{}, fields ...string) (*Map, error) {
	keyField := ""
	if len(fields) > 0 {
		keyField = fields[0]
	}
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
			switch V := val.(type) {
			case string:
				s, t, err := parsePaths(V)
				if err != nil {
					return m, err
				}
				(*m)[s] = t
			case map[interface{}]interface{}:
				key, ok := V[keyField].(string)
				if !ok {
					return m, fmt.Errorf("invalid map key (%s): %s", keyField, V)
				}
				// fmt.Println(V, keyField)
				(*m)[key] = V
			default:
				return m, fmt.Errorf("invalid map element: %s", V)
			}
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
			sv := s.(string)
			switch tv := t.(type) {
			case string:
				(*m)[sv] = tv
			case interface{}:
				(*m)[sv] = tv
			case []interface{}:
				(*m)[sv] = tv
			default:
				fmt.Println("??????", sv, tv)
			}
			// (*m)[s.(string)] = t.(string)
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

// MapPaths task
type MapPaths map[string]string

// Parse map string
func (m *MapPaths) Parse(i interface{}) error {
	mm, err := NewMapPaths(i)
	if err != nil {
		return err
	}
	*m = *mm
	return nil
}

// NewMapPaths parse
func NewMapPaths(i interface{}) (*MapPaths, error) {
	m := &MapPaths{}
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
			sv, ok := s.(string)
			if !ok {
				return m, fmt.Errorf("unable to parse key string %s in: %+v", sv, v)
			}
			tv, ok := t.(string)
			if !ok {
				return m, fmt.Errorf("unable to parse value string %s in: %+v", tv, v)
			}
			// switch tv := t.(type) {
			// case string:
			// 	(*m)[sv] = tv
			// // case interface{}:
			// default:
			// 	fmt.Println(sv, tv)
			// }
			// (*m)[s.(string)] = t.(string)
		}
	default:
		return m, fmt.Errorf("unable to parse map string: %+v", v)
	}
	return m, nil
}

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

// // IncludeMap ...
// type IncludeMap string

// func (m *IncludeMap) Parse(i interface{}) error {
// 	file := i.(string)
// 	if strings.HasPrefix(file, "~/") {
// 		file = filepath.Join(os.Getenv("HOME"), file[2:])
// 	}
// 	bytes, err := ioutil.ReadFile(file)
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			return nil
// 		}
// 		return err
// 	}
// 	if err := yaml.Unmarshal(bytes, &vars); err != nil {
// 		return err
// 	}
// 	*m
// 	return nil
// 	// mm, err := NewMap(i)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// *m = *mm
// 	// return nil
// }
