package role

import (
	"fmt"
)

type Dir struct {
	Path string
}

func (d *Dir) Get() string {
	return fmt.Sprintf("%s", d)
}

func (d *Dir) Set(value interface{}) {
	switch val := value.(type) {
	case string:
		d.Path = val
		// *d = append(*d, val)
	// case []string:
	// 	*d = make(Dir, 0)
	// 	for _, v := range val {
	// 		*d = append(*d, v)
	// 	}
	// }
	default:
		*d = val.(Dir)
	}
}

func (d *Dir) Sync() error {
	fmt.Println("Sync", d)
	return nil
}
