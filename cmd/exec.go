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
	// "os"
	// "os/exec"
	// "strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Shell = "bash"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cloneOrPull(Directory)
		if err != nil {
			return err
		}
		if len(args) == 0 {
			args = append(args, viper.GetStringSlice("exec")...)
		}
		if len(args) == 0 {
			return nil
		}
		args = append([]string{"-c"}, args...)
		err = executeCmd(Shell, args...)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(execCmd)

	execCmd.Flags().StringVarP(&Shell, "shell", "s", Shell, "Shell (bash|sh)")
}
