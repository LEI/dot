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
	"strings"

	// "github.com/LEI/dot/dotlib"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
)

// pkgCmd represents the package command
var pkgCmd = &cobra.Command{
	Use:   "pkg",
	Short: "Create or remove packages",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	installCmd.AddCommand(pkgCmd)
	removeCmd.AddCommand(pkgCmd)
}

// InstallPackages ...
func InstallPackages(args []interface{}) error {
	return pkgCommand(INSTALL, args)
}

// RemovePackages ...
func RemovePackages(args []interface{}) error {
	return pkgCommand(REMOVE, args)
}

func pkgCommand(method string, args []interface{}) error {
	if len(args) == 0 {
		return nil
	}
	pacaptArgs := []string{}
	if method == INSTALL {
		pacaptArgs = append(pacaptArgs, "-S")
	} else if method == REMOVE {
		pacaptArgs = append(pacaptArgs, "-R")
	}
	osList := listOS()
	if ok := hasOne([]string{"darwin"}, osList); !ok {
		pacaptArgs = append(pacaptArgs, "--noconfirm")
	}
	fmt.Println(osList)
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			pacaptArgs = append(pacaptArgs, v)
		// case map[interface{}]string:
		// 	fmt.Println("---------", v)
		case map[interface{}]interface{}:
			if v["os"] != nil {
				vOS := v["os"].([]interface{})
				osPkg := make([]string, len(vOS))
				for i := range vOS {
					osPkg[i] = vOS[i].(string)
				}
				if ok := hasOne(osPkg, osList); !ok {
					break
				}
				// TODO Check "!"+OS
			}
			pacaptArgs = append(pacaptArgs, v["name"].(string))
		}
	}
	fmt.Println("pacapt", strings.Join(pacaptArgs, " "))
	err := execute(pacaptBin, pacaptArgs...)
	// out, err := execPac(pacaptArgs...)
	if err != nil {
		return err
	}
	// fmt.Println(out)
	return nil
}

// func execPac(args ...string) (string, error) {
// 	if DryRun {
// 		return "# (SKIPPED: dry-run)", nil
// 	}
// 	output, err := exec.Command(pacaptBin, args...).CombinedOutput()
// 	str := strings.TrimRight(string(output), "\n")
// 	if err != nil {
// 		return str, err
// 	}
// 	return str, err
// }
