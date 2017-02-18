package role

import (
	"fmt"
	"github.com/LEI/dot/fileutil"
	"path/filepath"
	"os"
)

type File struct { // []string
	Path string
	Type string
	Stat *os.FileInfo
}

func (f *File) Glob() ([]*File, error) {
	var files []*File
	paths, err := filepath.Glob(f.Path)
	// fmt.Printf("Find: %s -> %+v\n", link.Path, paths)
	if err != nil {
		return files, err
	}
	for _, src := range paths {
		// f, err := &File{Path: src, Type: fileType} // NewFile(src)
		f, err := NewFile(&File{Path: src, Type: f.Type})
		if err != nil {
			return files, err
		}
		files = append(files, f)
	}
	return files, nil
}

func NewFile(value interface{}) (*File, error) {
	var file *File
	switch val := value.(type) {
	case string:
		if val == "" {
			return nil, fmt.Errorf("Empty File path\n")
		}
		file = &File{Type: "", Path: val}
	case map[string]interface{}:
		file = &File{Type: val["type"].(string), Path: val["path"].(string)}
	case *File:
		file = val
	default:
		// file = val
		return nil, fmt.Errorf("Unknown type %T for %+v\n", val, val)
	}
	file.Path = os.ExpandEnv(file.Path)
	return file, nil
}

func NewFileBase(value interface{}, base string) (*File, error) {
	file, err := NewFile(value)
	if err != nil {
		return file, err
	}
	file.Path = filepath.Join(base, file.Path)
	return file, nil
}

func (f *File) Base() string {
	return filepath.Base(f.Path)
}

func (f *File) Sync() error {
	fmt.Println("Sync", f.Path)
	return nil
}

func (f *File) Link(target string) error {
	return fileutil.Link(f.Path, target)
	// if err != nil {
	// 	return err
	// }

}

func (f *File) String() string {
	return fmt.Sprintf("%s", f.Path)
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
