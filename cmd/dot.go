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
	HomeDir  = os.Getenv("HOME")
	destination   = HomeDir
	source   string
	URL      string
	Config   config
	cfgType  string
	cfgFile  string
	cfgDir               = []string{"$HOME", "/etc/dot"}
	dotDir               = ".dot" // Default clone directory under $HOME
	dotCfg               = ".dot" // Default config file name without extension
	envKeys              = []string{"OS"}
	DirMode  os.FileMode = 0755
	FileMode os.FileMode = 0644
)

type config struct {
	Roles []role
}

type role struct {
	Name      string
	Dir       string `mapstructure:"directory"`
	URL       string
	OS        strSlice
	taskRoles `mapstructure:",squash"`
	taskExec  `mapstructure:",squash"`
}

type taskRoles struct {
	Env      map[string]string // Environment variables map
	Line     map[string]string // Lines map
	Link     strSlice          // Paths list `<source>[:<target>]`
	Template strSlice          // Paths list `<source>[:<target>]`
}

type taskExec struct {
	Install     strSlice // Exec before install
	PostInstall strSlice `mapstructure:"post_install"` // Exec after install
	Remove      strSlice // Exec before remove
	PostRemove  strSlice `mapstructure:"post_remove"` // Exec after remove
}

type strSlice []string

func (s *strSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *strSlice) Set(value string) error {
	// fmt.Printf("%s (%t)\n", value, value)
	*s = append(*s, value)
	return nil
}

// DotCmd represents the base command when called without any subcommands
var DotCmd = &cobra.Command{
	Use:   "dot",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return installCommand(args)
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := DotCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	DotCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "$HOME/.dot.yaml", "Config file")
	DotCmd.PersistentFlags().StringVarP(&cfgType, "format", "f", cfgType, "Config type: json, toml or yaml")

	// Local flags will only run when this action is called directly.
	// DotCmd.Flags().StringVarP(&Directory, "dir", "d", Directory, "Repository path")

	// viper.BindPFlag("directory", DotCmd.Flags().Lookup("directory"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if source != "" {
		cfgDir = append([]string{source}, cfgDir...)
	}
	readConfig(viper.GetViper(), cfgDir...)
	if err := viper.Unmarshal(&Config); err != nil {
		fmt.Fprintf(os.Stderr, "# Unable to decode into struct, %v", err)
		os.Exit(1)
	}
}

func readConfig(v *viper.Viper, dirs ...string) *viper.Viper {
	if cfgFile != "" { // Enable ability to specify config file via flag
		v.SetConfigFile(cfgFile)
	}
	if cfgType != "" { // Enable ability to specify config file format
		v.SetConfigType(cfgType)
	}
	v.SetConfigName(dotCfg) // Name of config file (without extension)

	for _, dir := range dirs { // Add directories to look for the config file in
		v.AddConfigPath(dir)
	}

	v.AutomaticEnv() // Read in environment variables that match
	// v.WatchConfig()  // Read config file while running
	// If a config file is found, read it in.
	if err := v.ReadInConfig(); err == nil {
		fmt.Println("# Using config file:", v.ConfigFileUsed())
	}
	return v
}

func readRoleConfig(r *role) error {
	v := viper.New()
	readConfig(v, r.Dir)
	if err := v.UnmarshalKey("role", &r); err != nil {
		return err
	}
	// fmt.Printf("%v\n", r)
	return nil
}

func getRole(dir, url string) (*role, error) {
	r := &role{Dir: dir, URL: url}
	if err := syncCommand(r.Dir, r.URL); err != nil {
		return r, err
	}
	if err := readRoleConfig(r); err != nil {
		return r, fmt.Errorf("# Unable to decode into struct, %v", err)
	}
	return r, nil
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
		parts = append(parts, destination)
	} else if len(parts) != 2 {
		fmt.Println("Invalid arg", arg)
		os.Exit(1)
	}
	src := os.ExpandEnv(parts[0])
	if !path.IsAbs(src) {
		src = path.Join(baseDir, src)
	}
	src = path.Clean(src)
	dst := os.ExpandEnv(parts[1])
	if !path.IsAbs(dst) {
		dst = path.Join(destination, dst)
	}
	changed, err := createDir(dst)
	if err != nil {
		return err
	}
	if changed {
		fmt.Printf("mkdir -p %s\n", dst)
	}
	return cb(src, dst)
}

func initCmd(action string, args []string) error {
	roles, err := filter(Config.Roles, args)
	if err != nil {
		return err
	}
	Config.Roles = roles
	if len(Config.Roles) == 0 {
		return fmt.Errorf("404 role not found\n")
	}
	for index, role := range Config.Roles {
		r, err := initRole(role)
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "# Install %s\n", role.Name)

		if err := syncCommand(role.Dir, role.URL); err != nil {
			return err
		}
		if err := readRoleConfig(&role); err != nil {
			fmt.Fprintf(os.Stderr, "# Unable to decode into struct, %v", err)
			// os.Exit(1)
			return nil
		}
		roleEnv, err := initEnv(role.Env)
		if err != nil {
			return err
		}
		switch action {
		case "install":
			if err := execCommand(role.Install); err != nil {
				return err
			}
			if err := linkCommand(role.Link, role.Dir); err != nil {
				return err
			}
			if err := templateCommand(role.Template, role.Dir, roleEnv); err != nil {
				return err
			}
			if err := lineCommand(role.Line); err != nil {
				return err
			}
			if err := execCommand(role.PostInstall); err != nil {
				return err
			}
		case "remove":
			if err := execCommand(role.Remove); err != nil {
				return err
			}
			// if err := linkCommand(role.Link, role.Dir); err != nil {
			// 	return err
			// }
			// if err := templateCommand(role.Template, role.Dir, roleEnv); err != nil {
			// 	return err
			// }
			// if err := lineCommand(role.Line); err != nil {
			// 	return err
			// }
			if err := execCommand(role.PostRemove); err != nil {
				return err
			}
		default:
			fmt.Printf("Unknown action '%s'", action)
		}
		// Config.Roles[index] = r
	}
	return nil
}

func filter(roles []role, patterns []string) ([]role, error) {
	if len(patterns) == 0 {
		return roles, nil
	}
	out := roles[:0]
	for _, r := range roles {
		matched, err := match(r.Name, patterns...)
		if err != nil {
			return out, err
		}
		if matched {
			out = append(out, r)
		}
	}
	return out, nil
}

func match(str string, patterns ...string) (bool, error) {
	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, str)
		if err != nil || matched {
			return matched, err
		}
	}
	return false, nil
}

func initRole(role role) (role, error) {
	if role.Name == "" {
		fmt.Fprintf(os.Stderr, "Missing role name in %v\n", role)
		os.Exit(1)
	}
	if role.URL == "" {
		fmt.Fprintf(os.Stderr, "Missing role url in %v\n", role)
		os.Exit(1)
	}
	if role.OS != nil {
		if ok := hasOne(role.OS, getOS()); !ok { // Skip role
			fmt.Fprintf(os.Stderr, "# Skip %s (%s)\n", role.Name, strings.Join(role.OS, ", "))
			return role, nil
		}
	}
	if role.Dir == "" {
		role.Dir = path.Join(destination, dotDir, role.Name)
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

func initEnv(in map[string]string) (map[string]string, error) {
	// env := Env()
	env := make(map[string]string, 0)
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
			// fmt.Printf("%s=\"%s\"\n", k, v)
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
	for _, k := range envKeys {
		if _, ok := env[k]; !ok {
			env[k] = OS
		}
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
