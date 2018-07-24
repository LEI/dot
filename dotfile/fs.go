package dotfile

import (
	"fmt"
	// "io"
	"os"

	// "golang.org/x/tools/godoc/vfs"

	"github.com/absfs/osfs"
)

//"github.com/blang/vfs"
var fs osfs.FileSystem

func init() {
	fs, err := osfs.NewFS()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if Verbose > 42 {
		fmt.Println(fs)
	}
}

// // FileSystem ...
// type FileSystem interface {
// 	Open(name string) (File, error)
// 	Stat(name string) (os.FileInfo, error)
// 	Create(name string) (File, error)
// 	Remove(name string) error
// 	IsExist(err error) bool
// 	IsNotExist(err error) bool
// }

// // File ...
// type File interface {
// 	Name() string
// 	Sync() error
// 	Truncate(int64) error
// 	Stat() (os.FileInfo, error)
// 	Readdir(count int) ([]os.FileInfo, error)
// 	io.Reader
// 	io.ReaderAt
// 	io.Writer
// 	io.Seeker
// 	io.Closer
// }

// // osFS implements FileSystem using the local disk.
// type osFS struct{
// 	// FileSystem
// 	// Stdout, Stderr *os.File
// }

// func (osFS) Open(name string) (File, error)        { return os.Open(name) }
// func (osFS) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }
// func (osFS) Create(name string) (File, error) { return os.Create(name) }
// func (osFS) Remove(name string) error { return os.Remove(name) }
// func (osFS) IsExist(err error) bool { return os.IsExist(err) }
// func (osFS) IsNotExist(err error) bool { return os.IsNotExist(err) }
