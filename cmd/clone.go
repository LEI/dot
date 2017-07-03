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
	// "fmt"
	"os"

	"github.com/spf13/cobra"
	// "gopkg.in/src-d/go-git.v4"
)

var (
	URL string
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone a git repository",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fi, err := os.Stat(Directory)
		if err != nil && os.IsExist(err) {
			// fmt.Fprintln(os.Stderr, err)
			// os.Exit(1)
			return err
		}
		if fi != nil {
			return pullRepository(Directory, "origin", "master")
		}
		return cloneRepository(URL, Directory)
	},
}

func init() {
	RootCmd.AddCommand(cloneCmd)

	// cloneCmd.PersistentFlags().String("foo", "", "A help for foo")

	cloneCmd.Flags().StringVarP(&URL, "url", "u", "", "Repository URL")
}

func cloneRepository(url string, dir string) error {
	args := []string{"clone", url, dir, "--recursive", "--quiet"}
	err := executeCmd("git", args...)
	if err != nil {
		return err
	}
	return nil
}

func pullRepository(dir string, remote string, branch string) error {
	args := []string{"-C", dir, "pull", remote, branch, "--quiet"}
	err := executeCmd("git", args...)
	if err != nil {
		return err
	}
	return nil
}

/*func cloneRepository(url string, dir string) (*git.Repository, error) {
	var isBare = false
	fmt.Printf("git clone %s %s --recursive\n", url, dir)
	r, err := git.PlainClone(dir, isBare, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		return r, err
	}
	// // Branch being pointed by HEAD
	// ref, err := r.Head()
	// if err != nil {
	// 	return err
	// }
	// // Commit object
	// commit, err := r.CommitObject(ref.Hash())
	// if err != nil {
	// 	return err
	// }
	return r, nil
}*/
