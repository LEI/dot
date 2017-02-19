package fileutil

import (
	"fmt"
	"os"
)

type File interface {
	init() error
	run() error
}

// type File struct {
// 	// Name string
// 	Path string
// 	Stat *os.FileInfo
// }

type Dir struct {
	Path string
}

type Link struct {
	Path   string
	Target string
}

func (l *Link) init() {
	fmt.Println("init", l)
}

func (l *Link) run() {
	fmt.Println("run", l)
}

type Line struct {
	Path string
	Line string
}

func NewLink(src string, dst string) string {
	return &Link{Path: src, Target: dst}
}

// func (f *File) Link(dst string) error {
// 	fmt.Printf("Link %s -> %s", f.path, dst)
// 	return nil
// }

func (f *File) Stat() (*os.FileInfo, error) {
	fi, err := os.Stat(f.Path)
	if err != nil {
		return fi, err
	}
	f.Stat = fi
	return fi, nil
}

func (f *File) String() string {
	return fmt.Sprintf("%+v", f)
}

func Run(f *File) error {
	err := f.init()
	if err != nil {
		return err
	}
	err := f.run()
	if err != nil {
		return err
	}
}
