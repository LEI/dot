// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	// "strings"

	"github.com/spf13/cobra"
)

var (
	// Source, Target         string
	// OnlyDirs       bool
	defaultDirMode os.FileMode = 0755
)

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Symlink",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cloneOrPull(Directory)
		if err != nil {
			return err
		}
		return parseArgs("link", args, func(source, target string) error {
			err := linkParse(source, target, Directory)
			if err != nil {
				return err
			}
			return nil
		})
	},
}

func init() {
	RootCmd.AddCommand(linkCmd)

	// linkCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory")
	// linkCmd.Flags().StringVarP(&Target, "target", "t", "", "Target `path`, directory or file")
	// linkCmd.Flags().BoolVarP(&OnlyDirs, "only-dirs", "", OnlyDirs, "Ignore files")
}

func linkParse(source, target, dir string) error {
	paths, err := filepath.Glob(source)
	if err != nil {
		return err // false
	}
	for _, p := range paths {
		_, f := path.Split(p)
		t := path.Join(target, f)
		err := linkFile(p, t)
		if err != nil {
			return err
		}
	}
	return nil // true
}

func linkFile(source, target string) error {
	real, err := readLink(target)
	if err != nil && os.IsExist(err) {
		return err
	}
	if real == source { // Symlink already exists
		fmt.Printf("ln -s %s %s\n", source, target)
		return nil
	}
	if real != "" {
		fmt.Fprintf(os.Stderr, "[WARN] target %s is a link to %s\n", target, real)
		// return nil
	}
	fi, err := os.Stat(target)
	if err != nil && os.IsExist(err) {
		return err
	}
	if fi != nil {
		fmt.Fprintf(os.Stderr, "[WARN] target %s exists\n", target)
		// os.Exit(1)
		return nil
	}
	fmt.Printf("ln -s %s %s\n", source, target)
	err = os.Symlink(source, target)
	if err != nil {
		return err
	}
	return nil
}

func isSymlink(fi os.FileInfo) bool {
	return fi != nil && fi.Mode()&os.ModeSymlink != 0
}

func readLink(path string) (string, error) {
	fi, err := os.Lstat(path)
	if err != nil { // os.IsExist(err)
		// if os.IsNotExist(err) {
		// return path, nil
		// }
		return "", err
	}
	if !isSymlink(fi) {
		return "", nil
	}
	real, err := os.Readlink(path)
	return real, err
}
