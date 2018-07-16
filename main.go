package main

import (
	"fmt"
	"os"
	// "path/filepath"
	"strings"

	// "github.com/jessevdk/go-flags"

	"github.com/LEI/dot/cmd"
	"github.com/LEI/dot/dotfile"
)

var (
	// Shell ...
	Shell = "bash"

	config = cmd.NewConfig()
	// *cmd.Config = &cmd.Config{} // {Name: ".dot"}
	configFileUsed string
)

func main() {
	// cmd.Options.Source = "" // .
	// cmd.Options.Target = "$HOME" // os.Getenv("HOME")
	cmd.Options.Config = func(s string) error {
		cmd.ConfigName = s
		configFile, err := config.Read(s)
		if err != nil {
			return err
		}
		configFileUsed = configFile
		return nil
	}

	// Parse arguments
	remaining, err := cmd.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %s\n", err)
		os.Exit(1)
	}
	if len(remaining) > 0 {
		// Hint -h, --help?
		fmt.Fprintf(os.Stderr, "Remaining arguments: %+v\n", remaining)
		os.Exit(1)
	}
	if cmd.Options.Version {
		fmt.Println("v0")
		os.Exit(0)
	}
	// if cmd.Verbose > 0 {
	// 	fmt.Printf("Verbosity: %v\n", cmd.Verbose)
	// }
	// fmt.Printf("Config: %+v\n", config)
	// fmt.Printf("Config roles: %+v\n", config.Roles)
	// fmt.Printf("Options: %+v\n", cmd.Options)
	// fmt.Printf("Options role: %+v\n", cmd.GetParser().Find("install").FindOptionByLongName("roles"))
	if configFileUsed != "" && cmd.Verbose > 0 {
		fmt.Printf("# Using configuration file: %s\n", configFileUsed)
	}
	if err := execute(&cmd.Options); err != nil {
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

func execute(options *cmd.DotCmd) error {
	// Do not attempt to sync for listing
	listOnly := dotfile.HasOne([]string{"list"}, cmd.RunOnly)
	// fmt.Println(len(config.Roles), "ROLES")
	// Initialize role config
	enabledCount := 0
	for _, r := range config.Roles {
		if len(options.RoleFilter) > 0 && !dotfile.HasOne([]string{r.Name}, options.RoleFilter) {
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
		enabledCount++
	}
	// Compute length after removing roles
	length := len(config.Roles)
	if length == 0 {
		return fmt.Errorf("# No roles, at least for this OS")
	}
	errs := make(chan error, length)
	for _, r := range config.Roles {
		go func(r *cmd.Role) {
			if !r.IsEnabled() {
				errs <- nil
				return
			}
			if err := r.Init(); err != nil {
				errs <- fmt.Errorf("# [%s] init error: %s", r.Name, err)
				return
			}
			if !listOnly {
				if cmd.Verbose > 0 {
					fmt.Printf("# [%s] Syncing %s %s\n", r.Name, r.Path, r.URL)
				}
				if err := r.Sync(); err != nil {
					if err == cmd.ErrDirtyRepo {
						errs <- fmt.Errorf("# [%s] Uncommitted changes in %s, use --force to continue", r.Name, r.Path)
						return
					}
					errs <- fmt.Errorf("# [%s] sync error: %s", r.Name, err)
					return
				}
			}
			configFile, err := r.ReadConfig(cmd.ConfigName)
			if err != nil {
				errs <- err
				return
			}
			if configFile != "" && cmd.Verbose > 1 {
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
	if enabledCount == 0 {
		return fmt.Errorf("# No roles to execute in: %s", config.Roles)
	}
	switch cmd.Action {
	case "install":
		// Resolve dependencies
		if err := config.Require(); err != nil {
			return err
			// switch err {
			// case cmd.ErrSkipDeps:
			// 	break
			// default:
			// 	return err
			// }
		}
		break
	}
	return config.Do(cmd.Action, cmd.RunOnly)
}
