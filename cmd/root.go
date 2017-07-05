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
	// "strconv"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	OS = runtime.GOOS
)

var (
	HomeDir   = os.Getenv("HOME")
	Target    = HomeDir
	Directory      string
	Cfg       Config
	cfgFormat string
	cfgFile   string
	cfgDir    = []string{"$HOME", "/etc/dot"}
	dotDir    = ".dot"
	dotCfg    = ".dot"
	DirMode os.FileMode = 0755
	FileMode os.FileMode = 0644
)

type Config struct {
	Roles []Role
}

type Role struct {
	Name     string
	Dir      string `mapstructure:"directory"`
	URL      string
	OS       stringSlice
	Task     `mapstructure:",squash"`
}

type Task struct {
	Exec     stringSlice
	Link     stringSlice
	Template stringSlice
	Line     map[string]string
	Done     stringSlice
	Env      map[string]string
}

// type Value interface {
// 	String() string
// 	Set(string) error
// }

type stringSlice []string

func (s *stringSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *stringSlice) Set(value string) error {
	// fmt.Printf("%s (%t)\n", value, value)
	*s = append(*s, value)
	return nil
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
		for index, role := range Cfg.Roles {
			r, err := initRole(role)
			if err != nil {
				return err
			}
			Cfg.Roles[index] = r
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
	RootCmd.PersistentFlags().StringVarP(&Target, "target", "t", Target, "Target directory")
	RootCmd.PersistentFlags().StringVarP(&cfgFormat, "format", "f", cfgFormat, "Data format (json|toml|yaml)")

	// Local flags will only run when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// viper.BindPFlag("directory", RootCmd.PersistentFlags().Lookup("directory"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if Directory != "" {
		cfgDir = append([]string{Directory}, cfgDir...)
	}
	readConfig(viper.GetViper(), cfgDir...)
	if err := viper.Unmarshal(&Cfg); err != nil {
		fmt.Fprintf(os.Stderr, "# Unable to decode into struct, %v", err)
		os.Exit(1)
	}
}

func readConfig(v *viper.Viper, dirs ...string) *viper.Viper {
	if cfgFile != "" { // Enable ability to specify config file via flag
		v.SetConfigFile(cfgFile)
	}
	if cfgFormat != "" { // Enable ability to specify config file format
		v.SetConfigType(cfgFormat)
	}
	v.SetConfigName(dotCfg)    // Name of config file (without extension)
	for _, dir := range dirs { // Add directories to look for the config file in
		v.AddConfigPath(dir)
	}
	v.AutomaticEnv() // Read in environment variables that match
	v.WatchConfig()  // Read config file while running
	// If a config file is found, read it in.
	if err := v.ReadInConfig(); err == nil {
		fmt.Println("# Using config file:", v.ConfigFileUsed())
	}
	return v
}

func readRoleConfig(r *Role) error {
	v := viper.New()
	readConfig(v, r.Dir)
	if err := v.UnmarshalKey("role", &r); err != nil {
		return err
	}
	// fmt.Printf("%v\n", r)
	return nil
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

func parseArg(arg, baseDir string, cb func(string, string) error) error {
	parts := strings.Split(arg, ":")
	if len(parts) == 1 {
		parts = append(parts, Target)
	} else if len(parts) != 2 {
		fmt.Println("Invalid arg", arg)
		os.Exit(1)
	}
	source := os.ExpandEnv(parts[0])
	if !path.IsAbs(source) {
		source = path.Join(baseDir, source)
	}
	source = path.Clean(source)
	target := os.ExpandEnv(parts[1])
	if !path.IsAbs(target) {
		target = path.Join(Target, target)
	}
	changed, err := createDir(target)
	if err != nil {
		return err
	}
	if changed {
		fmt.Printf("mkdir -p %s\n", target)
	}
	return cb(source, target)
}

func initRole(role Role) (Role, error) {
	if role.Name == "" {
		fmt.Fprintf(os.Stderr, "Missing role name in %v\n", role)
		os.Exit(1)
	}
	if role.URL == "" {
		fmt.Fprintf(os.Stderr, "Missing role url in %v\n", role)
		os.Exit(1)
	}
	if role.OS != nil {
		if ok := hasOne(role.OS, getOS()); !ok {
			fmt.Fprintf(os.Stderr, "## Skipping %s (%s)\n", role.Name, strings.Join(role.OS, ", "))
			return role, nil
		}
	}
	if role.Dir == "" {
		role.Dir = path.Join(Target, dotDir, role.Name)
	}
	if err := syncCommand(role.Dir, role.URL); err != nil {
		return role, err
	}
	if err := readRoleConfig(&role); err != nil {
		fmt.Fprintf(os.Stderr, "# Unable to decode into struct, %v", err)
		// os.Exit(1)
		return role, nil
	}
	if err := execCommand(role.Exec); err != nil {
		return role, err
	}
	if err := linkCommand(role.Link, role.Dir); err != nil {
		return role, err
	}
	if err := templateCommand(role.Template, role.Dir, role.Env); err != nil {
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

func GetEnv(in map[string]string) (map[string]string, error) {
	env := Env()
	// env := make(map[string]string, 0)
	// for key, val := range Env() {
	// 	env[key] = val
	// }
	for k, v := range in {
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
	return true, os.MkdirAll(dir, DirMode)
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

/*func getSlice(arg interface{}) []string {
	var slice []string
	switch v := arg.(type) {
	case []interface{}:
		// fmt.Println("[]interface")
		for _, r := range v {
			slice = append(slice, r.(string))
		}
	case interface{}:
		// fmt.Println("interface")
		slice = append(slice, v.(string))
	case string:
		// fmt.Println("string")
		slice = append(slice, v)
		// :
		// 	fmt.Println("other:", v)
	}
	return slice
}*/
