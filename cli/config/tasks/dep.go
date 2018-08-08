package tasks

import (
	// "fmt"

	"github.com/LEI/dot/cli/config/types"
)

// Deps task
type Deps struct {
	types.Slice
	// Vars map[string]interface{}
	// value []string
}

// Parse dep
// func (d *Deps) Parse(i interface{}) error {
// 	// if i == nil {
// 	// 	return nil
// 	// }
// 	// switch v := i.(type) {
// 	// case string:
// 	// 	*d = append(*d, v)
// 	// case []string:
// 	// 	*d = append(*d, v...)
// 	// case []interface{}:
// 	// 	for _, s := range v {
// 	// 		*d = append(*d, s.(string))
// 	// 	}
// 	// default:
// 	// 	return fmt.Errorf("unable to parse role deps: %+v", v)
// 	// }
// 	return d.Parse(i)
// }
