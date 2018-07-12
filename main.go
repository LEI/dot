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

	verbosity = 0
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

	verbosity = len(cmd.Options.Verbose)
	// if verbosity > 0 {
	// 	fmt.Printf("Verbosity: %v\n", verbosity)
	// }

	// fmt.Printf("Config: %+v\n", config)
	// fmt.Printf("Config roles: %+v\n", config.Roles)
	// fmt.Printf("Options: %+v\n", cmd.Options)
	// fmt.Printf("Options role: %+v\n", cmd.GetParser().Find("install").FindOptionByLongName("roles"))

	if configFileUsed != "" && verbosity > 0 {
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
	// fmt.Println(len(config.Roles), "ROLES")
	// Initialize role config
	for _, r := range config.Roles {
		if len(options.RoleFilter) > 0 && !dotfile.HasOne([]string{r.Name}, options.RoleFilter) {
			// fmt.Fprintf(os.Stderr, "# [%s] Skipping (filtered)\n", r.Name)
			config.Roles = removeRole(config.Roles, r)
			continue
		}
		if r.OS != nil {
			if len(r.OS) > 0 && !dotfile.HasOSType(r.OS...) { // Skip role
				fmt.Fprintf(os.Stderr, "# [%s] Skipping (OS: %s)\n", r.Name, strings.Join(r.OS, ", "))
				config.Roles = removeRole(config.Roles, r)
				continue
			}
		}
	}
	// Compute length after filtering roles
	length := len(config.Roles)
	errs := make(chan error, length)
	for _, r := range config.Roles {
		go func(r *cmd.Role) {
			if err := r.Init(); err != nil {
				errs <- fmt.Errorf("# %s init error: %s", r.Name, err)
				return
			}
			configFile, err := r.ReadConfig(cmd.ConfigName)
			if err != nil {
				errs <- err
				return
			}
			if err := r.Prepare(); err != nil {
				errs <- err
				return
			}
			if configFile != "" && verbosity > 0 {
				fmt.Printf("# [%s] Using role configuration file: %s\n", r.Name, configFile)
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
	if len(config.Roles) == 0 {
		fmt.Fprintln(os.Stderr, "No roles to execute")
		return nil
	}
	for _, r := range config.Roles {
		// if err := r.Prepare(); err != nil {
		// 	return err
		// }
		if err := r.Do(cmd.Action, cmd.Only); err != nil {
			return err
		}
	}
	return nil
}

func removeRole(roles []*cmd.Role, rm *cmd.Role) (ret []*cmd.Role) {
	for _, r := range roles {
		if r == rm {
			continue
		}
		ret = append(ret, r)
	}
	return ret
}
