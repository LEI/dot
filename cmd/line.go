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

	"github.com/LEI/dot/dotlib"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// lineCmd represents the line command
var lineCmd = &cobra.Command{
	Use:   "line",
	Short: "Add or remove a line in file",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		r := &role{
			Dir: Source,
			URL: URL,
		}
		if err := r.Init(); err != nil {
			return err
		}
		r.Line = viper.GetStringMapString("line")
		// switch action {
		// case "install":
		// 	err := InstallLine(in)
		// case "remove":
		// 	err := RemoveLine(in)
		// }
		return InstallLine(r.Line)
	},
}

func init() {
	installCmd.AddCommand(lineCmd)
	removeCmd.AddCommand(lineCmd)

	// lineCmd.Flags().StringVarP(&URL, "url", "u", URL, "Repository URL")
}

func InstallLine(in map[string]string) error {
	return lineCommand(in, dotlib.LineInFile)
}

func RemoveLine(in map[string]string) error {
	return lineCommand(in, nil)
}

func lineCommand(in map[string]string, action func(file string, line string) (bool, error)) error {
	if action == nil {
		return nil // Skip
	}
	for file, line := range in {
		p := parsePath(os.ExpandEnv(file), Target)
		// _, err := createDir(p)
		// if err != nil {
		// 	return err
		// }
		changed, err := action(p, line)
		if err != nil {
			return err
		}
		prefix := "# "
		if changed {
			prefix = ""
		}
		fmt.Printf("%secho '%s' >> \"%s\"\n", prefix, line, file)
	}
	return nil
}
