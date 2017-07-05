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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Fill go template",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := GetEnv(viper.GetStringMapString("env"))
		if err != nil {
			return err
		}
		return templateCommand(args, Directory, env)
	},
}

func init() {
	RootCmd.AddCommand(templateCmd)

	templateCmd.Flags().StringVarP(&Directory, "dir", "d", Directory, "Repository path")
}

func templateCommand(in []string, dir string, env map[string]string) error {
	for _, arg := range in {
		err := parseArg(arg, dir, func(source, target string) error {
			_, f := path.Split(source)
			target = path.Join(target, strings.TrimSuffix(f, ".tpl"))
			changed, err := templateGlob(source, target, env)
			if err != nil {
				return err
			}
			prefix := "# "
			if changed {
				prefix = ""
			}
			fmt.Printf("%senvsubst < %s | tee %s\n", prefix, source, target)
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func templateGlob(source, target string, env map[string]string) (bool, error) {
	tmpl, err := template.ParseGlob(source)
	tmpl = tmpl.Option("missingkey=zero")
	if err != nil {
		return false, err
	}
	buf := &bytes.Buffer{}
	// env, err := GetEnv()
	// if err != nil {
	// 	return false, err
	// }
	err = tmpl.Execute(buf, env)
	if err != nil {
		return false, err
	}
	str := buf.String()
	b, err := ioutil.ReadFile(target)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if str != string(b) {
		err := ioutil.WriteFile(target, []byte(str), FileMode)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func WriteString(path string, str string) (bool, error) {
	fi, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, FileMode)
	defer fi.Close()
	if err != nil {
		return false, err
	}
	n, err := fi.WriteString(str)
	if err != nil {
		return false, err
	}
	if n == 0 {
		return false, nil
	}
	return true, nil
}
