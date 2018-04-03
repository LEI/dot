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
	return lineCommand("Install", in)
}

func RemoveLine(in map[string]string) error {
	return lineCommand("Remove", in)
}

func lineCommand(method string, in map[string]string) error {
	// if action == nil {
	// 	return nil // Skip
	// }
	for file, line := range in {
		filePath := parsePath(os.ExpandEnv(file), Target)
		// if method == INSTALL {
		// 	if _, err := createDir(dir); err != nil {
		// 		return err
		// 	}
		// }
		lif := &dotlib.LineTask{
			File: filePath,
			Line: line,
		}
		err := callMethod(lif, method)
		if err != nil {
			return err
		}
	}
	return nil
}
