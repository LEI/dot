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
		SuccessLogger.Printf("%s\n", dir)
	}
	return nil
}
