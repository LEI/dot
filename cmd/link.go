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
	// Source, Target string
	Directory      string
	defaultDirMode os.FileMode = 0755
)

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Symlink",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return linkCommand(args, Directory)
	},
}

func init() {
	RootCmd.AddCommand(linkCmd)

	linkCmd.Flags().StringVarP(&Directory, "dir", "d", Directory, "Repository path")
	// linkCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory")
	// linkCmd.Flags().StringVarP(&Target, "target", "t", "", "Target `path`, directory or file")
}

func linkCommand(in []string, dir string) error {
	for _, arg := range in {
		err := parseArg(arg, dir, func(source, target string) error {
			err := linkGlob(source, target)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func linkGlob(source, target string) error {
	paths, err := filepath.Glob(source)
	if err != nil {
		return err // false
	}
	for _, s := range paths {
		_, f := path.Split(s)
		t := path.Join(target, f)
		changed, err := linkFile(s, t)
		if err != nil {
			return err
		}
		prefix := "# "
		if changed {
			prefix = ""
		}
		fmt.Printf("%sln -s %s %s\n", prefix, s, t)
	}
	return nil // true
}

func linkFile(source, target string) (bool, error) {
	real, err := readLink(target)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if real == source { // Symlink already exists
		return false, nil
	}
	if real != "" {
		// fmt.Fprintf(os.Stderr, "# %s is a link to %s, not %s\n", target, real, source)
		// os.Exit(1)
		return false, fmt.Errorf("# %s is a link to %s, not to %s\n", target, real, source)
	}
	fi, err := os.Stat(target)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if fi != nil {
		// fmt.Fprintf(os.Stderr, "# %s is already a file\n", target)
		// os.Exit(1)
		return false, fmt.Errorf("# %s already exists, could not link %s\n", target, source)
	}
	err = os.Symlink(source, target)
	if err != nil {
		return false, err
	}
	return true, nil
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

func isSymlink(fi os.FileInfo) bool {
	return fi != nil && fi.Mode()&os.ModeSymlink != 0
}
