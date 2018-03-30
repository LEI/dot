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

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	// Hidden: true,
	Use:        "install [flags]",
	Aliases:    []string{"i"},
	SuggestFor: []string{"add"},
	Short:      "Install dotfiles",
	Long:       ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return initCmd("install", args...)
	},
}

func init() {
	DotCmd.AddCommand(installCmd)

	installCmd.PersistentFlags().StringVarP(&source, "source", "s", source, "Source directory")
	installCmd.PersistentFlags().StringVarP(&destination, "target", "t", destination, "Destination directory")
	installCmd.PersistentFlags().StringVarP(&URL, "url", "u", URL, "Remote URL")

	// installCmd.Flags().StringVarP(&Directory, "dir", "d", Directory, "Repository path")
}
