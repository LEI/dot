package role

import (
	"fmt"
	"github.com/LEI/dot/fileutil"
	"path/filepath"
	"os"
)

type Role struct {
	Name string
	Origin string
	Source string
	Target string
	handlers []Handler
}

func (r *Role) Sync() error {
	for _, h := range r.handlers {
		err := h.Sync(r.Target)
		if err != nil {
			return err
		}
	}
	return nil
}

type Handler interface {
	// Name() string
	// Set(interface{})
	Sync(target string) error
	// Stat() (*os.FileInfo, error)
	String() string
}

type File struct { // []string
	Path string
	Type string
	info os.FileInfo
}

type Link struct {
	*File
	Target string
}

func NewFile(value interface{}) (*File, error) {
	var file *File
	switch val := value.(type) {
	case string:
		file = &File{Type: "", Path: val}
	case map[string]interface{}:
		file = &File{Type: val["type"].(string), Path: val["path"].(string)}
	case *File:
		file = val
	default:
		// file = val
		return file, fmt.Errorf("Unknown type %T for %+v\n", val, val)
	}
	if file.Path == "" {
		return file, fmt.Errorf("Empty File path\n")
	}
	// if filepath.IsAbs(file.Path) {
	// 	fmt.Printf("%s: file path is not absolute\n", file.Path)
	// }
	file.Path = os.ExpandEnv(file.Path)
	return file, nil
}

func (f *File) Init() error {

	return nil
}

func (f *File) NameMatches(patterns []string) (bool, error) {
	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, f.Base())
		if err != nil {
			return matched, err
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
}

func (f *File) Stat() (os.FileInfo, error) {
	fi, err := os.Stat(f.Path)
	// (*f.info) = fi
	f.info = fi
	if err != nil {
		return fi, nil
	}
	return fi, nil
}

func (f *File) IsDir() bool {
	if f.info == nil {
		fi, err := f.Stat()
		if err != nil || fi == nil {
			fmt.Printf("%s: %s", f, err)
			return false
		}
	}
	return f.info.IsDir()
}

func (f *File) Base() string {
	return filepath.Base(f.Path)
}

func (f *File) Sync(target string) error {
	fmt.Println("Sync File", f, target)
	return nil
}

func (f *File) String() string {
	return fmt.Sprintf("%s", f.Path)
}

func (f *File) GlobAsLink() ([]*Link, error) {
	var result []*Link
	paths, err := filepath.Glob(f.Path)
	// fmt.Printf("Find: %s -> %+v\n", link.Path, paths)
	if err != nil {
		return result, err
	}
	for _, src := range paths {
		// f, err := &File{Path: src, Type: fileType} // NewFile(src)
		link, err := NewLink(&File{Path: src, Type: f.Type})
		if err != nil {
			return result, err
		}
		switch link.Type {
		case "directory":
			if !link.IsDir() {
				continue
			}
		case "file":
			if link.IsDir() {
				continue
			}
		}
		result = append(result, link)
	}
	return result, nil
}

// func (f *File) Set(value interface{}) {
// 	fmt.Println("Set", f, value)
// 	switch vaf := value.(type) {
// 	case string:
// 		l.Path = val
// 		// *f = append(*l, val)
// 	default:
// 		*f = val.(File)
// 	}
// }

// func NewLink(file *File, target string) error
func NewLink(value interface{}) (*Link, error) {
	file, err := NewFile(value)
	if err != nil {
		return nil, err
	}
	link := &Link{File: file}
	return link, nil
}

func (l *Link) Sync(target string) error {
	fmt.Println("Sync Link", l, target)
	return fileutil.Link(l.Path, target)
	// if err != nil {
	// 	return err
	// }
}

func (l *Link) String() string {
	return fmt.Sprintf("%s", l.Path)
}
