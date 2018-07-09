package main

import (
	"fmt"
	"os"

	// "github.com/jessevdk/go-flags"

	"github.com/LEI/dot/cmd"
	"github.com/LEI/dot/cfg"
)

// Config ...
type Config struct {
	Roles []Role
}

// Role ...
type Role struct {
	Name string
	URL string
	OS []string
}

var (
	config Config
	configFile string
	options = cmd.Options
)

func main() {
	remaining, err := cmd.Parse()
	if err != nil {
		fmt.Println("Command error:", remaining)
		cmd.Help(1)
	}
	if len(remaining) > 0 {
		fmt.Println("Remaining arguments:", remaining)
		cmd.Help(1)
	}

	// TODO: before Execute()
	if options.Version {
		fmt.Println("v0")
		os.Exit(0)
	}

	verbosity := len(options.Verbose)

	if verbosity > 0 {
		fmt.Printf("Verbosity: %v\n", verbosity)
	}

	cfg.Debug = options.Debug

	configFile, err := cfg.Load(&config, options.Config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if configFile != "" {
		fmt.Println("Final configuration file:", configFile)
	}

	// FIXME: &config not working?
	//  dot i -c .dot.yml -d
	fmt.Printf("Config: %+v\n", config.Roles)
	fmt.Printf("Options: %+v\n", options.Roles)
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
