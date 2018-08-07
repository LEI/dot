package tasks

import (
	"fmt"
)

// Deps task
type Deps []string
// struct {
// 	Vars map[string]interface{}
// 	value []string
// }

// Parse data
func (d *Deps) Parse(i interface{}) error {
	if i == nil {
		return nil
	}
	switch v := i.(type) {
	case string:
		*d = append(*d, v)
	case []string:
		*d = append(*d, v...)
	case []interface{}:
		for _, s := range v {
			*d = append(*d, s.(string))
		}
	default:
		return fmt.Errorf("unable to parse role deps: %+v", v)
	}
	return nil
}
