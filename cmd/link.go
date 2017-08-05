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
	// "os"
	"path"
	"path/filepath"
	// "strings"

	"github.com/LEI/dot/dot"
	"github.com/LEI/dot/helpers"

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

func linkCommand(in []string, dir string) error {
	var linkList []dot.Link
	for _, arg := range in {
		err := parseArg(arg, dir, func(source, target string) error {
			paths, err := filepath.Glob(source)
			if err != nil {
				return err // false
			}
			for _, s := range paths {
				_, f := path.Split(s)
				t := path.Join(target, f)
				linkList = append(linkList, dot.Link{s, t})
				// changed, err := linkFile(s, t)
				// if err != nil {
				// 	return err
				// }
				// prefix := "# "
				// if changed {
				// 	prefix = ""
				// }
				// fmt.Printf("%sln -s %s %s\n", prefix, s, t)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	fmt.Println("Links:", linkList)
	return nil
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
		// fmt.Fprintf(os.Stderr, "# %s is a link to %s, not %s", target, real, source)
		// os.Exit(1)
		return false, fmt.Errorf("# %s is a link to %s, not to %s", target, real, source)
	}
	fi, err := os.Stat(target)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if fi != nil {
		// fmt.Fprintf(os.Stderr, "# %s is already a file", target)
		// os.Exit(1)
		return false, fmt.Errorf("# %s already exists, could not link %s", target, source)
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
