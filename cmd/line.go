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
	"path"
	// "path/filepath"
	// "strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
)

// lineCmd represents the line command
var lineCmd = &cobra.Command{
	Use:   "line",
	Short: "Add or remove a line in file",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return lineCommand(viper.GetStringMapString("line"))
	},
}

func init() {
	RootCmd.AddCommand(lineCmd)

	// lineCmd.Flags().BoolVarP(&Bool, "bool", "b", Bool, "Example boolean")
}

func lineCommand(in map[string]string) error {
	for dst, str := range in {
		dst = os.ExpandEnv(dst)
		if !path.IsAbs(dst) {
			dst = path.Join(Target, dst)
		}
		dst = path.Clean(dst)
		// _, err := createDir(dst)
		// if err != nil {
		// 	return err
		// }
		err := lineInFile(dst, str)
		if err != nil {
			return err
		}
	}
	return nil
}

func lineInFile(file string, line string) error {
	str := `if ! test -f "` + file + `" || ! grep -Fxq '` + line + `' "` + file + `"; then echo '` + line + `' >> "` + file + `"; fi;`
	err := executeCmd(Shell, []string{"-c", str}...)
	if err != nil {
		return err
	}

	return nil
}
