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
	// "io/ioutil"
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
	HomeDir     = os.Getenv("HOME")
	Target      = HomeDir
	Directory   string
	Config      config
	cfgFormat   string
	cfgFile     string
	dotDir      = ".dot"
)

type config struct {
	Roles []role
}

type role struct {
	Name string
	URL string `mapstructure:"url"`
	OS []string
	Directory string
	Exec []string
	Link []string
	Template []string
	Line map[string]string
	Done []string
	Env map[string]string
}

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
		if err := viper.Unmarshal(&Config); err != nil {
			fmt.Fprintf(os.Stderr, "# Unable to decode into struct, %v", err)
			os.Exit(1)
		}
		for index, role := range Config.Roles {
			r, err := initRole(role)
			if err != nil {
				return err
			}
			Config.Roles[index] = r
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
	RootCmd.PersistentFlags().StringVarP(&Directory, "dir", "d", Directory, "Repository path")
	RootCmd.PersistentFlags().StringVarP(&Target, "target", "t", Target, "Target directory")
	RootCmd.PersistentFlags().StringVarP(&cfgFormat, "format", "f", cfgFormat, "Data format (json|toml|yaml)")

	// Local flags will only run when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// viper.BindPFlag("directory", RootCmd.PersistentFlags().Lookup("directory"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}
	if cfgFormat != "" {
		viper.SetConfigType(cfgFormat)
	}
	viper.SetConfigName(".dot")  // name of config file (without extension)
	viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.AddConfigPath("$HOME/.dot")
	viper.AddConfigPath("/etc/dot")
	viper.AddConfigPath(Directory)
	viper.AutomaticEnv() // read in environment variables that match

	viper.WatchConfig() // Read config file while running

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("# Using config file:", viper.ConfigFileUsed())
	}
}

/*func readStdin(args []string, cb func(string, string) error) error {
	// if len(args) >= 1 && args[0] == "-" { // Read config from stdin
	// 	in, err := ioutil.ReadAll(os.Stdin)
	// 	if err != nil {
	// 		return fmt.Errorf("Error occured while reading from stdin: %s.", err)
	// 	}
	// 	viper.ReadConfig(bytes.NewBuffer(in))
	// 	args = viper.GetStringSlice(key)
	// } else if viper.ConfigFileUsed() != "" {
	// 	args = viper.GetStringSlice(key)
	// }
}*/

func parseArg(arg string, cb func(string, string) error) error {
	parts := strings.Split(arg, ":")
	if len(parts) == 1 {
		parts = append(parts, Target)
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
	if !path.IsAbs(target) {
		target = path.Join(Target, target)
	}
	_, err := createDir(target)
	if err != nil {
		return err
	}
	return cb(source, target)
}

func initRole(role role) (role, error) {
	if len(role.OS) > 0 {
		if ok := hasOne(role.OS, getOS()); !ok {
			fmt.Printf("## Skipping %s (OS: %s)\n", role.Name, strings.Join(role.OS, ", "))
			return role, nil
		}
	}

	if role.Name == "" {
		fmt.Fprintln(os.Stderr, "Missing role name")
		os.Exit(1)
	}
	if role.URL == "" {
		fmt.Fprintln(os.Stderr, "Missing role url")
		os.Exit(1)
	}

	URL = role.URL

	if role.Directory == "" {
		role.Directory = path.Join(Target, dotDir, role.Name)
	}

	Directory = role.Directory

	if err := syncCommand(role.Directory, role.URL); err != nil {
		return role, err
	}
	if err := execCommand(role.Exec); err != nil {
		return role, err
	}
	if err := linkCommand(role.Link); err != nil {
		return role, err
	}
	if err := templateCommand(role.Template); err != nil {
		return role, err
	}
	if err := lineCommand(role.Line); err != nil {
		return role, err
	}
	if err := execCommand(role.Done); err != nil {
		return role, err
	}
	return role, nil
}

func getOS() []string {
	types := []string{OS}
	OSTYPE, ok := os.LookupEnv("OSTYPE")
	if ok && OSTYPE != "" {
		types = append(types, OSTYPE)
	} else { // !ok || OSTYPE == ""
		// fmt.Printf("OSTYPE='%s' (%v)\n", OSTYPE, ok)
		out, err := exec.Command(Shell, "-c", "printf '%s' \"$OSTYPE\"").Output()
		if err != nil {
			fmt.Println(err)
		}
		if len(out) > 0 {
			OSTYPE = string(out)
			o := strings.Split(OSTYPE, ".")
			if len(o) > 0 {
				types = append(types, o[0])
			}
			types = append(types, OSTYPE)
		}
	}
	if OSTYPE == "" {
		fmt.Println("OSTYPE is not set or empty")
	}
	return types
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
		if v != "" { // Set the environment variable
			fmt.Printf("%s=\"%s\"\n", k, v)
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
	c := exec.Command(name, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func hasOne(in []string, list []string) bool {
	for _, a := range in {
		for _, b := range list {
			if b == a {
				return true
			}
		}
	}
	return false
}
