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
	"github.com/LEI/dot/dot"
)

var (
	// OS ...
	OS = runtime.GOOS

	// Shell ...
	Shell = "bash"

	config *dot.Config = &dot.Config{} // {Name: ".dot"}
	configFileUsed string
)

func main() {
	cmd.Options.Source = ""
	cmd.Options.Target = "$HOME" // os.Getenv("HOME")
	cmd.Options.ConfigName = ".dot.yml"
	cmd.Options.Config = func(s string) error {
		configFile, err := config.Load(s)
		if err != nil {
			return err
		}
		configFileUsed = configFile
		return nil
	}

	cmd.GlobalConfig = config

	// Parse arguments
	remaining, err := cmd.Parse()
	if err != nil {
		fmt.Println("Command error:", err)
		cmd.Help(1)
	}
	if len(remaining) > 0 {
		fmt.Println("Remaining arguments:", remaining)
		cmd.Help(1)
	}

	// FIXME: &config not working?
	//  dot i -c .dot.yml -d

	// TODO: before Execute()
	if cmd.Options.Version {
		fmt.Println("v0")
		os.Exit(0)
	}

	verbosity := len(cmd.Options.Verbose)
	// if verbosity > 0 {
	// 	fmt.Printf("Verbosity: %v\n", verbosity)
	// }

	dot.NoCheck = cmd.Options.NoCheck
	dot.NoSync = cmd.Options.NoSync
	dot.Source = string(cmd.Options.Source)
	dot.Target = string(cmd.Options.Target)
	dot.Verbose = verbosity > 0

	// fmt.Printf("Config: %+v\n", config)
	// fmt.Printf("Config roles: %+v\n", config.Roles)
	// fmt.Printf("Options: %+v\n", cmd.Options)

	if configFileUsed != "" && dot.Verbose {
		fmt.Println("# Using configuration file:", configFileUsed)
	}

	if err := execute(&cmd.Options); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// for k, v := range cmd.Options.Roles {
	// 	fmt.Println(k, v)
	// 	config.Roles = append(config.Roles, dot.Role{
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
		go func (r *dot.Role) {
			if err := r.Init(); err != nil {
				errs <- fmt.Errorf("# %s init error: %s", r.Name, err)
				return
			}
			configFile, err := r.LoadConfig(options.ConfigName)
			if err != nil {
				errs <- err
				return
			}
			if configFile != "" && dot.Verbose {
				fmt.Printf("# [%s] Using role configuration file: %s\n", r.Name, configFile)
			}
			errs <- nil
		}(r)
	}

	for i := 0; i < length; i++ {
		if err := <- errs; err != nil {
			fmt.Printf("Role initialization failed: %d/%d\n", i+1, length)
			return err
		}
	}
	if len(config.Roles) == 0 {
		fmt.Fprintln(os.Stderr, "No roles to execute")
		return nil
	}
	return config.Do(options.ActionFilter)
}

func removeRole(roles []*dot.Role, rm *dot.Role) (ret []*dot.Role) {
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
