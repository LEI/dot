package main

import (
	"fmt"
	"os"
	"path/filepath"

	// "github.com/jessevdk/go-flags"

	"github.com/LEI/dot/cfg"
	"github.com/LEI/dot/cmd"
	"github.com/LEI/dot/dot"
)

var (
	// OS ...
	// OS = runtime.GOOS

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
	for i, r := range config.Roles {
		if r.Path == "" {
			r.Path = filepath.Join(
				string(options.Target),
				string(options.RoleDir),
				r.Name)
		}
		if err := r.Init(); err != nil {
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
