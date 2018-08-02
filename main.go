package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/imdario/mergo"
	// "github.com/jessevdk/go-flags"

	"github.com/LEI/dot/dotfile"
	"github.com/LEI/dot/utils"
)

var (
	config         = &Config{} // {Name: ".dot"}
	configFileUsed string
)

func main() {
	// Options.Source = "" // .
	// Options.Target = "$HOME" // os.Getenv("HOME")
	Options.Config = func(name string) error {
		// fmt.Println("# Config name:", name)
		cfgPath, err := FindConfig(name)
		if err != nil {
			return err
		}
		if cfgPath == "" {
			fmt.Println("Unable to find roles config file")
			return nil
		}
		_, s := filepath.Split(cfgPath)

		// TODO allow alternative role config names
		ConfigName = s

		if err := config.Read(cfgPath); err != nil {
			return err
		}
		// fmt.Printf("CONFIG: %+v\n", config)
		configFileUsed = cfgPath
		// configFile != "" &&
		// if configFile != ConfigName {
		// 	configFileUsed = configFile
		// }
		return nil
	}

	// Parse arguments
	remaining, err := Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %s\n", err)
		os.Exit(1)
	}
	if len(remaining) > 0 {
		// Hint -h, --help?
		fmt.Fprintf(os.Stderr, "Remaining arguments: %+v\n", remaining)
		os.Exit(1)
	}
	if Options.Version {
		fmt.Println("v0")
		os.Exit(0)
	}
	// if Verbose > 0 {
	// 	fmt.Printf("Verbosity: %v\n", Verbose)
	// }
	// fmt.Printf("Config: %+v\n", config)
	// fmt.Printf("Config roles: %+v\n", config.Roles)
	// fmt.Printf("Options: %+v\n", Options)
	// fmt.Printf("Options role: %+v\n", GetParser().Find("install").FindOptionByLongName("roles"))
	if configFileUsed != "" && Verbose > 0 {
		fmt.Printf("# Using configuration file: %s\n", configFileUsed)
	}
	if err := execute(&Options); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	if err := dotfile.InitEnv(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func execute(options *DotCmd) error {
	// listOnly := dotfile.Contains(RunOnly, "list")
	// fmt.Println(len(config.Roles), "ROLES")
	// Initialize role config
	// enabledCount := 0
	for _, r := range config.Roles {
		if len(options.RoleFilter) > 0 && !dotfile.Contains(options.RoleFilter, r.Name) {
			// fmt.Fprintf(os.Stderr, "# [%s] Skipping (filtered)\n", r.Name)
			// config.Roles = removeRole(config.Roles, r)
			// r.Disable()
			continue
		}
		r.Enabled = true // Quickfix
		if r.OS != nil {
			if len(r.OS) > 0 && !dotfile.HasOSType(r.OS...) { // Skip role
				fmt.Fprintf(os.Stderr, "# [%s] Skipping (OS: %s)\n", r.Name, strings.Join(r.OS, ", "))
				if err := config.RemoveRole(r); err != nil {
					return err
				}
				continue
			}
		}
		// enabledCount++
	}
	// Compute length after removing roles
	length := len(config.Roles)
	if length == 0 {
		return fmt.Errorf("# No roles, at least for this OS")
	}
	// Sync roles
	// TODO: refactor to sync initally disabled roles
	// but enabled by deps after initial sync below
	errs := make(chan error, length)
	for _, r := range config.Roles {
		go func(r *Role) {
			if !r.IsEnabled() {
				errs <- nil
				return
			}
			if err := r.Init(); err != nil {
				errs <- fmt.Errorf("# [%s] init error: %s", r.Name, err)
				return
			}
			if Verbose > 0 {
				fmt.Printf("# [%s] Syncing %s %s\n", r.Name, r.Path, r.URL)
			}
			if err := r.Sync(); err != nil {
				if err == ErrDirtyRepo {
					errs <- fmt.Errorf("# [%s] Uncommitted changes in %s, use --force to continue", r.Name, r.Path)
					return
				}
				errs <- fmt.Errorf("# [%s] sync error: %s", r.Name, err)
				return
			}
			configFile, err := r.ReadConfig(ConfigName)
			if err != nil {
				errs <- err
				return
			}
			if configFile != "" && Verbose > 1 {
				fmt.Printf("# [%s] Using role configuration file: %s\n", r.Name, configFile)
			}
			if err := r.Prepare(); err != nil {
				errs <- err
				return
			}
			errs <- nil
		}(r)
	}
	for i := 0; i < length; i++ {
		if err := <-errs; err != nil {
			fmt.Printf("Role initialization failed: %d/%d\n", i+1, length)
			return err
		}
	}
	// if enabledCount == 0 {
	// 	return fmt.Errorf("# No roles to execute in: %s", config.Roles)
	// }
	// TODO: sync missing roles enabled by dependencies
	switch Action {
	case "install":
		// Resolve dependencies
		if err := config.Require(); err != nil {
			return err
			// switch err {
			// case ErrSkipDeps:
			// default:
			// 	return err
			// }
		}
	}
	// // Install, remove...
	// if err := config.Do(Action, RunOnly); err != nil {
	// 	return err
	// }
	for _, r := range config.Roles {
		// Skip disabled roles
		if !r.IsEnabled() {
			continue
		}
		// if err := r.Prepare(); err != nil {
		// 	return err
		// }
		if err := do(r, Action); err != nil {
			return err
		}
	}
	// switch Action {
	// case "remove":
	// 	if err := config.ClearCache(); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

func do(r *Role, a string) error {
	act := strings.Title(a)
	// if shouldRunTask("list") {
	// 	// Just print the role fields
	// 	fmt.Printf("# LIST ROLE %+v\n", r.Print(Verbose > 0))
	// 	// return nil
	// }
	fmt.Printf("# Role: %s\n", r.Name) // (%v)\n", r.Name, r.IsEnabled())
	// originalEnv := dotfile.GetEnv() // Saved in dotfile.OriginalEnv
	if r.Env != nil {
		for k, v := range r.Env {
			k = strings.ToTitle(k)
			// Set role environment
			if err := dotfile.SetEnv(k, v); err != nil {
				// fmt.Fprintf(os.Stderr, err)
				return err
			}
		}
	}
	switch Action {
	case "install", "remove":
		// Pre-install/remove hook
		v := r.GetField(act)
		if !v.IsValid() {
			return fmt.Errorf("could not get field %s: %s", act, v)
		}
		before := v.Interface().([]string)
		if len(before) > 0 && shouldRunTask("exec") {
			for _, c := range before {
				task := &dotfile.ExecTask{
					Cmd: c,
				}
				str, err := task.Do(act)
				if err != nil {
					return err
				}
				if str != "" {
					fmt.Println(str)
				}
			}
		}
	}
	// System packages
	if r.Pkg != nil && Options.Packages && shouldRunTask("package") {
		for _, v := range r.Pkg {
			if v.OS != nil && len(v.OS) > 0 && !dotfile.HasOSType(v.OS.Value()...) {
				continue
			}
			if v.Action != "" && strings.ToLower(v.Action) != strings.ToLower(act) {
				continue
			}
			// args := []string{}
			// for _, a := range v.Args {
			// 	args = append(args, a)
			// }
			task := &dotfile.PkgTask{
				Name: v.Name,
				Opts: v.Args,
				Type: v.Type,
				Sudo: Options.Sudo,
			}
			str, err := task.Do(act)
			if err != nil {
				return err
			}
			if str != "" {
				fmt.Println(str)
			}
		}
	}
	// NOOP: Copies
	if r.Copies != nil && shouldRunTask("copy") {
		for s, t := range r.Copies {
			task := &dotfile.CopyTask{
				Source: s,
				Target: t,
			}
			str, err := task.Do(act)
			if err != nil {
				return err
			}
			if str != "" {
				fmt.Println(str)
			}
		}
	}
	// Line in file
	if r.Lines != nil && shouldRunTask("line") {
		for s, t := range r.Lines {
			task := &dotfile.LineTask{
				File: s,
				Line: t,
			}
			str, err := task.Do(act)
			if err != nil {
				return err
			}
			if str != "" {
				fmt.Println(str)
			}
		}
	}
	// Symlink files
	if r.Links != nil && shouldRunTask("link") {
		for s, t := range r.Links {
			task := &dotfile.LinkTask{
				Source: s,
				Target: t,
			}
			str, err := task.Do(act)
			if err != nil {
				return err
			}
			if str != "" {
				fmt.Println(str)
			}
		}
	}
	// Templates
	if r.Templates != nil && shouldRunTask("template") {
		for _, t := range r.Templates {
			env := r.Env // map[string]string{}
			if err := mergo.Merge(&env, t.Env); err != nil {
				return err
			}
			vars := r.Vars // map[string]interface{}{}
			if err := mergo.Merge(&vars, t.Vars); err != nil {
				return err
			}
			task := &dotfile.TemplateTask{
				Source: t.Source,
				Target: t.Target,
				Env:    env,
				Vars:   vars,
			}
			str, err := task.Do(act)
			if err != nil {
				return err
			}
			if str != "" {
				fmt.Println(str)
			}
		}
	}
	// Restore original environment
	if r.Env != nil {
		if err := dotfile.RestoreEnv(dotfile.OriginalEnv); err != nil {
			return nil
		}
	}
	switch Action {
	// case "list":
	// 	fmt.Printf("========= %+v\n", r)
	case "install", "remove":
		// Post-install/remove hook
		v := r.GetField("Post" + act)
		if !v.IsValid() {
			return fmt.Errorf("could not get field %s: %s", act, v)
		}
		after := r.GetField("Post" + act).Interface().([]string)
		if len(after) > 0 && shouldRunTask("exec") {
			for _, c := range after {
				task := &dotfile.ExecTask{
					Cmd: c,
				}
				str, err := task.Do(act)
				if err != nil {
					return err
				}
				if str != "" {
					fmt.Println(str)
				}
			}
		}
		defaultsFile := filepath.Join(r.Path, "defaults.yml")
		if utils.Exist(defaultsFile) {
			defaults := &dotfile.Defaults{}
			if err := defaults.Read(defaultsFile); err != nil {
				return err
			}
			if err := defaults.Parse(); err != nil {
				return err
			}
			if err := defaults.Exec(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Check if a given task name should be run
func shouldRunTask(s string) bool {
	return dotfile.Contains(RunOnly, s)
}
