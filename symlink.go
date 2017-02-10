package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func linkFiles(source string, dest string, globs []interface{}) error {
	if source == "" {
		return fmt.Errorf("Cannot link empty source")
	}
	if _, err := os.Stat(source); err != nil {
		return err
	}
	// var filePaths []string
	for _, glob := range globs {
		switch l := glob.(type) {
		case string:
			paths, _ := filepath.Glob(filepath.Join(source, expand(l)))
			// filePaths = append(filePaths, paths...)
			for _, src := range paths {
				// fmt.Printf("%+v\n", src)
				dst := strings.Replace(src, source, dest, 1)
				err := linkFile(src, dst)
				if err != nil {
					return err
				}
			}
		default:
			fmt.Println("Unhandled type for", l)
		}
	}
	return nil
}

func linkFile(src string, dst string) error {
	name := strings.Replace(src, source+PathSeparator, "", 1)
	fi, err := os.Lstat(dst)
	if err != nil && os.IsExist(err) {
		return err
	}
	if fi != nil && (fi.Mode()&os.ModeSymlink != 0) {
		link, err := os.Readlink(dst)
		if err != nil {
			return err
		}
		if link == src {
			fmt.Printf("%s %s == %s\n", OkSymbol, name, dst)
			return nil
		}
		msg := dst + " is an existing symlink to " + link + ", replace it with " + src + "?"
		if ok := confirm(msg); ok {
			err := os.Remove(dst)
			if err != nil {
				return err
			}
		}
		// return nil
	} else if fi != nil {
		msg := dst + " is an existing file, move it to " + dst + ".backup and replace it with " + src + "?"
		if ok := confirm(msg); ok {
			err := os.Rename(dst, dst+".backup")
			if err != nil {
				return err
			}
		}
	}
	err = os.Symlink(src, dst)
	if err != nil {
		return err
	}
	fmt.Printf("%s %s -> %s\n", OkSymbol, name, dst)
	return nil
}
