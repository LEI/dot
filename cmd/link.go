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
	"runtime"
	// "strings"

	"github.com/spf13/cobra"
)

const (
	OS = runtime.GOOS
)

var (
	Target         string
	OnlyDirs       bool
	defaultDirMode os.FileMode = 0755
)

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Symlink",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return parseArgs(args, func(source, target string) error {
			_, err := linkPattern(source, target, Directory)
			if err != nil {
				return err
			}
			return nil
		})
	},
}

func init() {
	err := os.Setenv("OS", OS)
	if err != nil {
		fmt.Println(err)
	}

	RootCmd.AddCommand(linkCmd)

	// linkCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory")
	linkCmd.Flags().StringVarP(&Target, "target", "t", "", "Target `path`, directory or file")
	// linkCmd.Flags().BoolVarP(&OnlyDirs, "only-dirs", "", OnlyDirs, "Ignore files")
}

func linkPattern(source, target, dir string) (bool, error) {
	paths, err := filepath.Glob(source)
	if err != nil {
		return false, err
	}
	for _, p := range paths {
		_, f := path.Split(p)
		t := path.Join(target, f)
		fmt.Printf("ln -s %s %s\n", p, t)
		err := os.Symlink(p, t)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}
