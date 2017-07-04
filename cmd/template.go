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
)

var (
	defaultFileMode os.FileMode = 0644
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Fill go template",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return templateCommand(args)
	},
}

func init() {
	RootCmd.AddCommand(templateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// templateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// templateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func templateCommand(in []string) error {
	for _, arg := range in {
		err := parseArg(arg, func(source, target string) error {
			err := templatePattern(source, target, Directory)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func templatePattern(source, target, dir string) error {
	_, f := path.Split(source)
	target = path.Join(target, strings.TrimSuffix(f, ".tpl"))
	tmpl, err := template.ParseGlob(source)
	tmpl = tmpl.Option("missingkey=zero")
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	env, err := GetEnv()
	if err != nil {
		return err
	}
	err = tmpl.Execute(buf, env)
	if err != nil {
		return err
	}
	str := buf.String()
	b, err := ioutil.ReadFile(target)
	if err != nil && os.IsExist(err) {
		return err
	}
	if str != string(b) {
		err := ioutil.WriteFile(target, []byte(str), defaultFileMode)
		if err != nil {
			return err
		}
	}
	fmt.Printf("envsubst < %s | tee %s\n", source, target)
	return nil
}

func WriteString(path string, str string) (bool, error) {
	fi, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, defaultFileMode)
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
