package role

import (
	"fmt"
	// "github.com/LEI/dot/fileutil"
	// "path/filepath"
	// "os"
)

type Glob struct {
	Dir string
	Pattern string
	Files []*File
}

// func (g *Glob) Set(value interface{}) *Glob {
// 	return value
// }

func (g *Glob) String() string {
	return fmt.Sprintf("%+v", g)
}

func (g *Glob) Search() error {
	fmt.Println(g)
	return nil
}
