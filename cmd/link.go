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
	"path"
	"path/filepath"
	// "strings"

	"github.com/LEI/dot/dotlib"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Create or remove symlinks",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		r := &role{
			Dir: Source,
			URL: URL,
		}
		if r.URL == "" {
			if err := viper.UnmarshalKey("url", &r.URL); err != nil {
				return err
			}
		}
		if err := r.Init(); err != nil {
			return err
		}
		r.Link = args
		if len(args) == 0 {
			if err := viper.UnmarshalKey("link", &r.Link); err != nil {
				return err
			}
		}
		return InstallLink(r.Link, Source)
	},
}

func init() {
	installCmd.AddCommand(linkCmd)
	removeCmd.AddCommand(linkCmd)

	// linkCmd.Flags().StringVarP(&URL, "url", "u", URL, "Repository URL")
}

// InstallLink ...
func InstallLink(in []string, dir string) error {
	return linkCommand(in, dir, linkGlob)
}

// RemoveLink ...
func RemoveLink(in []string, dir string) error {
	return linkCommand(in, dir, nil)
}

func linkCommand(in []string, dir string, action func(src, dst string, env map[string]string) error) error {
	// func linkCommand(in []string, dir string) error {
	if action == nil {
		return nil // Skip
	}
	env := map[string]string{} // Placeholder
	for _, arg := range in {
		err := parseArg(arg, dir, env, action)
		if err != nil {
			return err
		}
	}
	return nil
}

func linkGlob(src, dst string, env map[string]string) error {
	// var linkList []*dotlib.LinkTask
	paths, err := filepath.Glob(src)
	if err != nil {
		return err // false
	}
	for _, s := range paths {
		_, f := path.Split(s)
		t := path.Join(dst, f)

		link := &dotlib.LinkTask{
			Source: s,
			Target: t,
		}
		// TODO: switch Install/Remove Link
		if err = link.Install(); err != nil {
			return err
		}
		// linkList = append(linkList, &dotlib.LinkTask{s, t, })
	}
	// fmt.Println("Links:", linkList)
	return nil
}
