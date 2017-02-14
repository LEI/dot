package main

import (
	// "fmt"
	"os"
	"path/filepath"
)

func makeDirs(dst string, paths []string) error {
	for _, dir := range paths {
		dir = filepath.Join(dst, expand(dir))
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
		logSuccess.Printf("%s\n", dir)
	}
	return nil
}

func removeDirs(dst string, paths []string) error {
	for _, dir := range paths {
		dir = filepath.Join(dst, expand(dir))
		err := os.Remove(dir)
		if err != nil {
			return err
		}
		logSuccess.Printf("%s\n", dir)
	}
	return nil
}
