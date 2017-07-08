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
	"github.com/spf13/viper"
)

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Symlink",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := getRole(source, URL)
		if err != nil {
			return err
		}
		r.Link = args
		if len(args) == 0 {
			if err := viper.UnmarshalKey("link", &r.Link); err != nil {
				return err
			}
		}
		return InstallLink(r.Link, source)
	},
}

func init() {
	installCmd.AddCommand(linkCmd)
	removeCmd.AddCommand(linkCmd)

	// linkCmd.Flags().StringVarP(&URL, "url", "u", URL, "Repository URL")
}

func InstallLink(in []string, dir string) error {
	return linkCommand(in, dir, linkGlob)
}

func RemoveLink(in []string, dir string) error {
	return linkCommand(in, dir, nil)
}

func linkCommand(in []string, dir string, action func(src, dst string) error) error {
	if action == nil {
		return nil // Skip
	}
	for _, arg := range in {
		err := parseArg(arg, dir, action)
		if err != nil {
			return err
		}
	}
	return nil
}

func linkGlob(src, dst string) error {
	paths, err := filepath.Glob(src)
	if err != nil {
		return err // false
	}
	for _, s := range paths {
		_, f := path.Split(s)
		t := path.Join(dst, f)
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

func linkFile(src, dst string) (bool, error) {
	real, err := readLink(dst)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if real == src { // Symlink already exists
		return false, nil
	}
	if real != "" {
		// fmt.Fprintf(os.Stderr, "# %s is a link to %s, not %s", dst, real, src)
		// os.Exit(1)
		return false, fmt.Errorf("# %s is a link to %s, not to %s", dst, real, src)
	}
	fi, err := os.Stat(dst)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if fi != nil {
		// fmt.Fprintf(os.Stderr, "# %s is already a file", dst)
		// os.Exit(1)
		return false, fmt.Errorf("# %s already exists, could not link %s", dst, src)
	}
	err = os.Symlink(src, dst)
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
