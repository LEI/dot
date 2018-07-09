package main

import (
	"fmt"
	"os"

	// "github.com/jessevdk/go-flags"

	"github.com/LEI/dot/cfg"
	"github.com/LEI/dot/cmd"
	"github.com/LEI/dot/dot"
)

var (
	config *dot.Config = &dot.Config{}
	configFile string
)

func main() {
	cfg.Debug = true

	cmd.GlobalConfig = config

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

	// fmt.Printf("Options: %+v\n", options)

	// fmt.Printf("Config: %+v\n", config)
	// fmt.Printf("Config roles: %+v\n", config.Roles)
	// fmt.Printf("Options: %+v\n", options)
	for i, r := range config.Roles {
		fmt.Println("Role", i+1, r.Name)
		fmt.Println("Copies", r.Copy)
		fmt.Println("Links", r.Link)
		fmt.Println("Templates", r.Template)
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
