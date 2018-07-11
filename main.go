package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	// "path/filepath"
	"runtime"
	"strings"

	// "github.com/jessevdk/go-flags"

	"github.com/LEI/dot/cmd"
)

var (
	// OS ...
	OS = runtime.GOOS

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
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err)
		os.Exit(1)
	}
	if len(remaining) > 0 {
		// Hint -h, --help?
		fmt.Fprintf(os.Stderr, "Remaining arguments: %+v\n", remaining)
		os.Exit(1)
	}

	// FIXME: &config not working?
	//  dot i -c .dot.yml -d

	// TODO: before Execute()
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

	if configFileUsed != "" && verbosity > 0 {
		fmt.Println("# Using configuration file:", configFileUsed)
	}

	if err := execute(&cmd.Options); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// for k, v := range cmd.Options.Roles {
	// 	fmt.Println(k, v)
	// 	config.Roles = append(config.Roles, Role{
	// 		Name: k,
	// 		URL: v,
	// 	})
	// }
	// fmt.Printf("=> %+v roles\n", len(config.Roles))
	// fmt.Println("CLI role:", cmd.GetParser().Find("install").FindOptionByLongName("roles"))
}

func init() {
	if err := os.Setenv("OS", OS); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// i, err := parser.AddCommand("install",
	// 	"Install",
	// 	"",
	// 	&cmd.Install)

	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// // cmd.Install = i

	// r, err := parser.AddCommand("remove",
	// 	"Remove",
	// 	"",
	// 	&cmd.Remove)

	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// // cmd.Remove = r
	// fmt.Println(i, r)
	// fmt.Println(cmd.Install, cmd.Remove)
}

func execute(options *cmd.DotCmd) error {
	// fmt.Println(len(config.Roles), "ROLES")
	// Initialize role config
	for _, r := range config.Roles {
		if len(options.RoleFilter) > 0 && !hasOne([]string{r.Name}, options.RoleFilter) {
			// fmt.Fprintf(os.Stderr, "# [%s] Skipping (filtered)\n", r.Name)
			config.Roles = removeRole(config.Roles, r)
			continue
		}
		if r.OS != nil {
			if !hasOne(r.OS, getOsTypes()) { // Skip role
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

// List of OS name and family/type
func getOsTypes() []string {
	types := []string{OS}

	// Add OS family
	c := exec.Command(Shell, "-c", "cat /etc/*-release")
	stdout, _ := c.StdoutPipe()
	// stderr, _ := c.StderrPipe()
	c.Start()
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
		v := strings.TrimLeft(m, "ID=")
		if m != v {
			types = append(types, v)
			break
		}
	}
	c.Wait()

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
