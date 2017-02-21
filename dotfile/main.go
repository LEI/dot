package dotfile

import (
	// "fmt"
	"os"
	"path/filepath"
)

func Base(path string) string {
	return filepath.Base(path)
}

func Glob(pattern string) ([]string, error) {
	paths, err := filepath.Glob(pattern)
	return paths, err
}

func List(pattern string, callback ...func(*File) bool) ([]*File, error) {
	files := []*File{}
	paths, err := Glob(pattern)
	if err != nil {
		return files, err
	}
GLOB:
	for _, p := range paths {
		f := &File{path: p}
		_, err := f.Stat()
		if err != nil {
			return files, err
		}
		for _, cb := range callback {
			if !cb(f) {
				continue GLOB
			}
		}
		files = append(files, f)
	}
	return files, nil
}

func Match(name string, patterns ...string) (bool, error) {
	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, name)
		if err != nil {
			return matched, err
		}
		if matched {
			return true, nil
		}
	}

	return false, nil
}

func IsSymlink(fi os.FileInfo) bool {
	return fi != nil && fi.Mode()&os.ModeSymlink != 0
}

// func Symlink(f *File, dst string) error {
// 	return f.Symlink(dst)
// }

// switch ft {
// case "directory", "dir", "d":
// 	if !fi.IsDir() {
// 		continue GLOB
// 	}
// case "file", "f":
// 	if fi.IsDir() {
// 		continue GLOB
// 	}
// }
