package dotfile

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	DryRun bool // = true
	// AlreadyLinked = fmt.Errorf("Already linked")
	Skip = fmt.Errorf("Skip")
)

type FileHandler func(string, os.FileInfo) error

func Base(path string) string {
	return filepath.Base(path)
}

func Glob(pattern string) ([]string, error) {
	paths, err := filepath.Glob(pattern)
	return paths, err
}

func BaseMatch(name string, patterns ...string) (bool, error) {
	matched, err := Match(Base(name), patterns...)
	return matched, err
}

func Match(name string, patterns ...string) (bool, error) {
	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, name)
		if err != nil || matched {
			return matched, err
		}
	}
	return false, nil
}

func List(pattern string, handlers ...FileHandler) ([]string, error) {
	result := []string{}
	paths, err := Glob(pattern)
	if err != nil {
		return paths, err
	}
GLOB:
	for _, p := range paths {
		fi, err := os.Stat(p)
		if err != nil {
			return paths, err
		}
		for _, f := range handlers {
			err := f(p, fi)
			if err != nil {
				if err == Skip {
					continue GLOB
				}
				return paths, err
			}
		}
		// result[p] = fi
		result = append(result, p)
	}
	return result, nil
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
