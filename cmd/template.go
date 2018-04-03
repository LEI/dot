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
	// "bytes"
	// "fmt"
	// "io/ioutil"
	// "os"
	// "path"
	// "strings"
	// "text/template"

	"github.com/LEI/dot/dotlib"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Generate a file",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		r := &role{
			Dir:  Source,
			URL:  URL,
			Env: Env(),
		}
		if err := r.Init(); err != nil {
			return err
		}
		r.Template = args
		if len(args) == 0 {
			if err := viper.UnmarshalKey("template", &r.Template); err != nil {
				return err
			}
		}
		env := viper.GetStringMapString("env")
		env, err := initEnv(env)
		if err != nil {
			return err
		}
		/*
			roleEnv, err := initEnv(role.Env)
			if err != nil {
				return role, err
			}
		*/
		// for k, e := range role.task.Env {
		// 	env[k] = r
		// }
		return InstallTemplate(r.Template, Source, env)
	},
}

func init() {
	installCmd.AddCommand(templateCmd)

	// templateCmd.Flags().StringVarP(&URL, "url", "u", URL, "Repository URL")
	// templateCmd.Flags().StringVarP(&Extra, "env", "e", Extra, "Extra env var")
}

// InstallTemplate ...
func InstallTemplate(in []string, dir string, env map[string]string) error {
	return templateCommand(INSTALL, in, dir, env)
}

// RemoveTemplate ...
func RemoveTemplate(in []string, dir string, env map[string]string) error {
	return templateCommand(REMOVE, in, dir, env)
}

func templateCommand(method string, in []string, dir string, env map[string]string) error {
	// switch method.(type) {
	// case string:
	// 	break;
	// }
	// if action == nil {
	// 	return nil // Skip
	// }
	for _, arg := range in {
		src, dst := parseArg(arg, dir)
		if method == INSTALL {
			if _, err := createDir(dst); err != nil {
				return err
			}
		}
		tpl := &dotlib.TemplateTask{
			Source: src,
			Target: dst,
			Env:    env,
		}
		if err := callMethod(tpl, method); err != nil {
			return err
		}
	}
	return nil
}
