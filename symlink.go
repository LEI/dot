package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

func linkFiles(source string, target string, globs []interface{}) error {
	if source == "" {
		fmt.Printf("%s -> %s ~ %+v\n", source, target, globs)
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
			fmt.Printf("%s: unknown type '%s'\n", path, reflect.TypeOf(path))
		}
	}
	return nil
}

func findLinks(source string, target string, options *Link) error {
	paths, _ := filepath.Glob(filepath.Join(source, expand(options.Path)))
	if Verbose > 0 {
		fmt.Printf("LINK %s \t-> %s\nOPTIONS %+v\nPATHS %+v\n", source, target, options, paths)
	}
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

func checkLinkOptions(src string, dst string, options *Link) (bool, error) {
	// name := strings.Replace(src, source+PathSeparator, "", 1)
	for _, pattern := range IgnoreFiles {
		re := regexp.MustCompile(pattern)
		if re.FindStringIndex(filepath.Base(src)) != nil {
			return false, nil
		}
	}
	if options.Type != "" {
		fi, err := os.Stat(src)
		if err != nil {
			return false, err
		}
		switch options.Type {
		case "f", "file":
			if fi.IsDir() {
				return false, nil
			}
		case "d", "directory":
			if fi.IsDir() != true {
				return false, nil
			}
		}
	}
	return true, nil
}

func linkFile(src string, dst string, options *Link) error {
	if shouldLink, err := checkLinkOptions(src, dst, options); shouldLink == false {
		return err
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
