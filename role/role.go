package role

import (
	"os"
)

type Role struct {
	Name  string
	Path  string
	Files []File
}

type File struct {
	Path     string
	FileInfo *os.FileInfo
}

// func (role *Role) New(name string) Role {
// 	return *Role{Name: name}
// }

// func (role *Role) String() string {
// 	return *Role.Name
// }
