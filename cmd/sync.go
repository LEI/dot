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
	// "github.com/spf13/viper"
	// "gopkg.in/src-d/go-git.v4"
)

var (
	Remote string = "origin"
	Branch string = "master"
	Pull   bool   = true
	synced []string
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Clone or pull a git repository",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return syncCommand(Directory, URL) // args...
	},
}

func init() {
	RootCmd.AddCommand(syncCmd)

	RootCmd.PersistentFlags().BoolVarP(&Pull, "pull", "n", Pull, "Update if already cloned")

	syncCmd.Flags().StringVarP(&Directory, "dir", "d", Directory, "Repository path")
	syncCmd.Flags().StringVarP(&URL, "url", "u", URL, "Repository URL")
	syncCmd.Flags().StringVarP(&Remote, "remote", "r", Remote, "Remote name")
	syncCmd.Flags().StringVarP(&Branch, "branch", "b", Branch, "Target ref")

	// viper.BindPFlag("url", RootCmd.PersistentFlags().Lookup("url"))
}

// syncCommand clone or pull
func syncCommand(dir, url string) error {
	for _, c := range synced {
		if c == dir {
			// Already updated
			return nil
		}
	}
	fi, err := os.Stat(dir)
	if err != nil && os.IsExist(err) {
		// fmt.Fprintln(os.Stderr, err)
		// os.Exit(1)
		return err
	}
	if fi != nil {
		return pullRepo(dir, Remote, Branch)
	}
	if err = cloneRepo(dir, url); err != nil {
		return err
	}
	// if dir == Directory {}
	// Read config file
	// initConfig()
	return nil
}

func cloneRepo(dir, url string) error {
	args := []string{"clone", url, dir, "--recursive", "--quiet"}
	err := executeCmd("git", args...)
	if err != nil {
		return err
	}
	synced = append(synced, dir)
	return nil
}

func pullRepo(dir, remote, branch string) error {
	if !Pull {
		return nil
	}
	args := []string{"-C", dir, "pull", remote, branch, "--quiet"}
	err := executeCmd("git", args...)
	if err != nil {
		return err
	}
	synced = append(synced, dir)
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
