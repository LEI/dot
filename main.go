package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"runtime"

	// "github.com/jessevdk/go-flags"

	"github.com/LEI/dot/cfg"
	"github.com/LEI/dot/cmd"
	"github.com/LEI/dot/dot"
)

var (
	// OS ...
	OS = runtime.GOOS

	// Shell ...
	Shell = "bash"

	config *dot.Config = &dot.Config{}
	configFile string
)

// var defaults = map[string]interface{} {
// 	"Source": "a",
// 	"Config": ...,
// }

func main() {
	// cfg.Debug = true

	cmd.Options.Source = ""
	cmd.Options.Target = "$HOME" // os.Getenv("HOME")
	cmd.Options.RoleDir = ".dot"

	cmd.GlobalConfig = config

	cmd.Options.Config = func(s string) error {
		configFile, err := cfg.Load(config, s)
		if err != nil {
			return err
		}

		if configFile != "" {
			fmt.Println("Using configuration file:", configFile)
		}

		return nil
	}

	remaining, err := cmd.Parse()
	if err != nil {
		fmt.Println("Command error:", remaining)
		cmd.Help(1)
	}
	if len(remaining) > 0 {
		fmt.Println("Remaining arguments:", remaining)
		cmd.Help(1)
	}

	// Parsed arguments
	options := cmd.Options

	// FIXME: &config not working?
	//  dot i -c .dot.yml -d

	// TODO: before Execute()
	if options.Version {
		fmt.Println("v0")
		os.Exit(0)
	}

	// verbosity := len(options.Verbose)
	// if verbosity > 0 {
	// 	fmt.Printf("Verbosity: %v\n", verbosity)
	// }
	// cfg.Debug = verbosity > 0

	// fmt.Printf("Config: %+v\n", config)
	// fmt.Printf("Config roles: %+v\n", config.Roles)
	// fmt.Printf("Options: %+v\n", options)
	target := string(options.Target)
	for i, r := range config.Roles {
		if r.OS != nil {
			if ok := hasOne(r.OS, getOsTypes()); !ok { // Skip role
				fmt.Fprintf(os.Stderr, "# Skip %s (%s)\n", r.Name, strings.Join(r.OS, ", "))
				continue
			}
		}
		if err := r.Init(
			target,
			string(options.RoleDir)); err != nil {
			fmt.Println("Role", i, "error:", remaining)
			os.Exit(1)
		}
	}

	// fmt.Println("CFG")
	// cmd.WriteIniConfig(cmd.GetParser())
	// fmt.Println("ENDCFG")

	if err := config.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// for k, v := range options.Roles {
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
