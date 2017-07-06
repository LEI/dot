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
	"path/filepath"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			fmt.Println("Args:", args)
		}
		return installCommand(args)
	},
}

func init() {
	DotCmd.AddCommand(installCmd)

	// installCmd.Flags().StringVarP(&Directory, "dir", "d", Directory, "Repository path")
}

func installCommand(args []string) error {
	roles, err := filter(Config.Roles, args)
	if err != nil {
		return err
	}
	Config.Roles = roles
	if len(Config.Roles) == 0 {
		return fmt.Errorf("404 role not found\n")
	}
	for index, role := range Config.Roles {
		r, err := initRole(role)
		if err != nil {
			return err
		}
		Config.Roles[index] = r
	}
	return nil
}

func filter(roles []role, patterns []string) ([]role, error) {
	if len(patterns) == 0 {
		return roles, nil
	}
	out := roles[:0]
	for _, r := range roles {
		matched, err := match(r.Name, patterns...)
		if err != nil {
			return out, err
		}
		if matched {
			out = append(out, r)
		}
	}
	return out, nil
}

func match(str string, patterns ...string) (bool, error) {
	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, str)
		if err != nil || matched {
			return matched, err
		}
	}
	return false, nil
}
