package main

import (
	"fmt"
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
		fmt.Printf("%s %s\n", OkSymbol, dir)
	}
	return nil
}
