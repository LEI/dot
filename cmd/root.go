// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
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
	"os/exec"
	"path"
	"runtime"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	OS = runtime.GOOS
)

var (
	Directory string
	cfgFile   string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "dot",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cloneOrPull(Directory); err != nil {
			return err
		}
		linkList := viper.GetStringSlice("link")
		if err := linkCmd.RunE(linkCmd, linkList); err != nil {
			return err
		}
		templateList := viper.GetStringSlice("template")
		if err := templateCmd.RunE(templateCmd, templateList); err != nil {
			return err
		}
		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func init() {
	// err := os.Setenv("OS", OS)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, err)
	// }

	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "$HOME/.dot.yaml", "config file")
	RootCmd.PersistentFlags().StringVarP(&Directory, "dir", "d", "", "Directory path")

	// Local flags will onlt run when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".dot")  // name of config file (without extension)
	viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.AddConfigPath("$HOME/.dot")
	viper.AddConfigPath("/etc/dot")
	viper.AddConfigPath(Directory)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func parseArgs(key string, args []string, cb func(string, string) error) error {
	if len(args) == 1 && args[0] == "-" {
		format := "json"
		viper.SetConfigType(format)
		stdin, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("Error occured while reading from stdin: %s.", err)
		}
		viper.ReadConfig(bytes.NewBuffer(stdin))
		args = viper.GetStringSlice(key)
	} else if viper.ConfigFileUsed() != "" {
		args = viper.GetStringSlice(key)
	}
	for _, arg := range args {
		parts := strings.Split(arg, ":")
		if len(parts) == 1 {
			parts = append(parts, "$HOME")
		} else if len(parts) != 2 {
			fmt.Println("Invalid arg", arg)
			os.Exit(1)
		}
		source := os.ExpandEnv(parts[0])
		if !path.IsAbs(source) {
			source = path.Join(Directory, source)
		}
		source = path.Clean(source)
		target := os.ExpandEnv(parts[1])
		_, err := createDir(target)
		if err != nil {
			return err
		}
		err = cb(source, target)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetEnv() (map[string]string, error) {
	env := Env() // make(map[string]string, 0)
	for key, val := range Env() {
		env[key] = val
	}
	for k, v := range viper.GetStringMapString("env") {
		k = strings.ToTitle(k)
		if v == "" { // Lookup environment if the variable is empty
			val, ok := os.LookupEnv(k)
			if !ok {
				fmt.Fprintf(os.Stderr, "# LookupEnv failed for '%s'", k)
				continue
			}
			v = val
		}
		if v != "" { // Parse string as a template
			templ, err := template.New(k).Option("missingkey=zero").Parse(v)
			if err != nil {
				return env, err
			}
			buf := &bytes.Buffer{}
			err = templ.Execute(buf, Env())
			if err != nil {
				return env, err
			}
			v = buf.String()
		}
		// fmt.Printf("%s=\"%s\"\n", k, v)
		if v != "" { // Set the environment variable
			err := os.Setenv(k, v)
			if err != nil {
				return env, err
			}
		} else {
			fmt.Fprintf(os.Stderr, "# Empty environment variable '%s'", k)
		}
		env[k] = v
	}
	return env, nil
}

func Env() map[string]string {
	env := make(map[string]string, 0)
	for _, i := range os.Environ() {
		sep := strings.Index(i, "=")
		env[i[0:sep]] = i[sep+1:]
	}
	if _, ok := env["OS"]; !ok {
		env["OS"] = OS
	}
	return env
}

func createDir(dir string) (bool, error) {
	fi, err := os.Stat(dir)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if err == nil && fi.IsDir() {
		return false, nil
	}
	return true, os.MkdirAll(dir, defaultDirMode)
}

func executeCmd(name string, args ...string) error {
	fmt.Printf("%s %s\n", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
