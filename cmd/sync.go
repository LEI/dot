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
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
)

var (
	directory string
	remote    = "origin"
	branch    = "master"
	gitCheck  = true
	gitClone  = true
	gitPull   = true
	synced    []string
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Clone or pull a git repository",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return CloneOrPull(directory, URL) // args...
	},
}

func init() {
	DotCmd.AddCommand(syncCmd)

	// DotCmd.PersistentFlags().BoolVarP(&pull, "pull", "n", pull, "Update if already cloned")

	syncCmd.Flags().StringVarP(&directory, "dir", "d", directory, "Repository directory")
	syncCmd.Flags().StringVarP(&URL, "url", "u", URL, "Repository URL")
	syncCmd.Flags().StringVarP(&remote, "remote", "r", remote, "Remote name")
	syncCmd.Flags().StringVarP(&branch, "branch", "b", branch, "Target ref")

	// viper.BindPFlag("url", DotCmd.PersistentFlags().Lookup("url"))
}

// CloneOrPull clone or pull
func CloneOrPull(dir, url string) error {
	if dir == "" {
		return fmt.Errorf("Missing repository directory")
	}
	if url == "" {
		return fmt.Errorf("Missing repository url")
	}
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
	// if err = checkURL(url, remote); err != nil {
	// 	return err
	// }
	if fi != nil {
		err := checkRepo(dir, remote, branch)
		if err != nil {
			return err
		}
		return pullRepo(dir, remote, branch)
	}
	if err = cloneRepo(dir, url); err != nil {
		return err
	}
	// if dir == source {}
	// Read config file
	// initConfig()
	return nil
}

func checkURL(url, remote string) error {
	fmt.Println(url, remote)
	args := []string{"config", "--local", "--get", "remote.origin.url"}
	// fmt.Printf("git %s\n", strings.Join(args, " "))

	// TODO: catpure output and check url against `user/repo`
	err := executeCmd("git", args...)
	if err != nil {
		return err
	}
	return nil
}

func cloneRepo(dir, url string) error {
	if !gitClone {
		return nil
	}
	args := []string{"clone", url, dir, "--recursive", "--quiet"}
	fmt.Printf("git %s\n", strings.Join(args, " "))
	err := executeCmd("git", args...)
	if err != nil {
		return err
	}
	synced = append(synced, dir)
	return nil
}

func checkRepo(dir, remote, branch string) error {
	if !gitCheck {
		return nil
	}
	args := []string{"-C", dir, "diff-index", "--quiet", "HEAD"}
	c := exec.Command("git", args...)
	err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Uncommited changes in '%s', aborting", dir)
		os.Exit(1)
		// return err
	}
	return nil
}

func pullRepo(dir, remote, branch string) error {
	if !gitPull {
		return nil
	}
	args := []string{"-C", dir, "pull", remote, branch, "--quiet"}
	fmt.Printf("git %s\n", strings.Join(args, " "))
	err := executeCmd("git", args...)
	if err != nil {
		return err
	}
	synced = append(synced, dir)
	return nil
}
