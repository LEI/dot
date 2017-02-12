package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func linkFiles(source string, target string, globs []interface{}) error {
	if source == "" {
		return fmt.Errorf("Cannot link empty source")
	}
	if _, err := os.Stat(source); err != nil {
		return err
	}
	// var filePaths []string
	for _, glob := range globs {
		switch path := glob.(type) {
		case string:
			err := findLinks(source, target, &Link{Type: "", Path: path})
			if err != nil {
				return err
			}
		case map[string]interface{}:
			err := findLinks(source, target, &Link{
				Type: path["type"].(string),
				Path: path["path"].(string),
			})
			// err := findLinks(source, target, path.(Link))
			if err != nil {
				return err
			}
		default:
			fmt.Printf("%s: unknown type '%s'", path, reflect.TypeOf(path))
		}
	}
	return nil
}

func findLinks(source string, target string, options *Link) error {
	paths, _ := filepath.Glob(filepath.Join(source, expand(options.Path)))
	// filePaths = append(filePaths, paths...)
	for _, src := range paths {
		// fmt.Printf("%+v\n", src)
		dst := strings.Replace(src, source, target, 1)
		err := linkFile(src, dst, options)
		if err != nil {
			return err
		}
	}
	return nil
}

func linkFile(src string, dst string, options *Link) error {
	// name := strings.Replace(src, source+PathSeparator, "", 1)
	if options.Type != "" {
		fi, err := os.Stat(src)
		if err != nil {
			return err
		}
		switch options.Type {
		case "f", "file":
			if fi.IsDir() {
				return nil
			}
		case "d", "directory":
			if fi.IsDir() != true {
				return nil
			}
		}
	}
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
			fmt.Printf("%s %s == %s\n", OkSymbol, src, dst)
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
	fmt.Printf("%s %s -> %s\n", OkSymbol, src, dst)
	return nil
}
