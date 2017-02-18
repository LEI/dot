package role

import (
	"fmt"
)

type Directory struct {
	*File
}

func (d *Directory) Sync(target string) error {
	fmt.Println("Sync", d)
	return nil
}

func (d *Directory) String() string {
	return fmt.Sprintf("%s", d.Path)
}

// func (d *Dir) Set(value interface{}) {
// 	switch val := value.(type) {
// 	case string:
// 		d.Path = val
// 		// *d = append(*d, val)
// 	// case []string:
// 	// 	*d = make(Dir, 0)
// 	// 	for _, v := range val {
// 	// 		*d = append(*d, v)
// 	// 	}
// 	// }
// 	default:
// 		*d = val.(Dir)
// 	}
// }
