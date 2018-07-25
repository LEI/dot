package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	// "github.com/jessevdk/go-flags"

	"github.com/LEI/dot/dotfile"
)

var (
	config         *Config = &Config{} // {Name: ".dot"}
	configFileUsed string
)

func main() {
	// Options.Source = "" // .
	// Options.Target = "$HOME" // os.Getenv("HOME")
	Options.Config = func(name string) error {
		//fmt.Println("Config flag:", name)
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
	// Install, remove...
	if err := config.Do(Action, RunOnly); err != nil {
		return err
	}
	// switch Action {
	// case "remove":
	// 	if err := config.ClearCache(); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}
