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

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Fill go template",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := getRole(source, URL)
		if err != nil {
			return err
		}
		r.Template = args
		if len(args) == 0 {
			if err := viper.UnmarshalKey("template", &r.Template); err != nil {
				return err
			}
		}
		env, err := initEnv(viper.GetStringMapString("env"))
		if err != nil {
			return err
		}
		// roleEnv, err := initEnv(role.Env)
		// if err != nil {
		// 	return role, err
		// }
		return InstallTemplate(r.Template, source, env)
	},
}

func init() {
	installCmd.AddCommand(templateCmd)

	// templateCmd.Flags().StringVarP(&URL, "url", "u", URL, "Repository URL")
	// templateCmd.Flags().StringVarP(&Extra, "env", "e", Extra, "Extra env var")
}

func InstallTemplate(in []string, dir string, env map[string]string) error {
	return templateCommand(in, dir, env, templateGlob)
}

func RemoveTemplate(in []string, dir string, env map[string]string) error {
	return templateCommand(in, dir, env, nil)
}

func templateCommand(in []string, dir string, env map[string]string, action func(src, dst string, env map[string]string) (bool, error)) error {
	if action == nil {
		return nil // Skip
	}
	for _, arg := range in {
		err := parseArg(arg, dir, func(src, dst string) error {
			_, f := path.Split(src)
			dst = path.Join(dst, strings.TrimSuffix(f, ".tpl"))
			changed, err := action(src, dst, env)
			if err != nil {
				return err
			}
			prefix := "# "
			if changed {
				prefix = ""
			}
			for k, v := range env {
				fmt.Printf("%s=\"%s\"\n", k, v)
			}
			// fmt.Printf("%senvsubst < %s | tee %s\n", prefix, src, dst)
			fmt.Printf("%stemplate %s -> %s\n", prefix, src, dst)
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func templateGlob(src, dst string, env map[string]string) (bool, error) {
	tmpl, err := template.ParseGlob(src)
	if err != nil {
		return false, err
	}
	tmpl = tmpl.Option("missingkey=zero")
	buf := &bytes.Buffer{}
	// env, err := GetEnv()
	// if err != nil {
	// 	return false, err
	// }
	if err = tmpl.Execute(buf, env); err != nil {
		return false, err
	}
	str := buf.String()
	b, err := ioutil.ReadFile(dst)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if str == string(b) {
		return false, nil
	}
	if err := ioutil.WriteFile(dst, []byte(str), FileMode); err != nil {
		return false, err
	}
	return true, nil
}
