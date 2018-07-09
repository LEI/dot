package cmd

import (
	"fmt"
	"os"
	// "reflect"
	"strings"

	"github.com/LEI/dot/dot"

	"github.com/jessevdk/go-flags"
)

// Executable ...
type Executable interface {
	Execute(args []string) error
}

// BaseCmd ...
type BaseCmd struct {
	// FIXME *flags.Command
	flags.Commander

	// Role map[string]string `short:"r" long:"role" description:""`
}

// BaseRoleCmd ...
type BaseRoleCmd struct {
	BaseCmd

	// Paths map[flags.Filename]flags.Filename `name:"PATHS" description:"Paths"`
	// Paths map[string]string `name:"PATHS" description:"Paths"`
	// Source flags.Filename `name:"SOURCE" description:"Path to source file"`
	// Target flags.Filename `name:"TARGET" description:"Path to target link"`

	// Role dot.Role `required:"false" positional-args:"true" positional-arg-name:"role"`
	// Name string `short:"n" long:"name" description:""`
	// URL string `short:"u" long:"url" description:""`

	// *PathArg // `required:"true" positional-args:"true"`
	// Paths dot.Role `required:"" positional-args:"true" positional-arg-name:"paths"`
	Role RoleArg `positional-args:"true" positional-arg-name:"paths"` // required:"1"

	Executable
}

// RoleArg ...
type RoleArg struct {
	Name string // `default:"default"` // `required:"true"`
	Paths []flags.Filename
}

// New dot.Role
func (a *RoleArg) New() *dot.Role {
	var name = a.Name
	switch name {
	case "":
		name = "default"
		break
	case "all":
		name = "*"
		break
	}
	url := ""
	if strings.Contains(name, ":") {
		parts := strings.Split(name, ":")
		if len(parts) == 2 {
			name = parts[0]
			url = parts[1]
		} else {
			fmt.Println("Unhandled", name)
			os.Exit(1)
		}
	}
	if strings.Contains(name, "*") {
		// find glob
	}
	return dot.NewRole(name, url)
}

// PathArg ...
// type PathArg struct {
// 	Source flags.Filename `short:"s" long:"source" description:"Path to source file"`
// 	Target flags.Filename `short:"t" long:"target" description:"Path to target link"`
// }

/*
var commands = []string{
	"Copy",
	"Link",
	"Template",
}

func executeCommands(cmd interface{}, args []string) error {
	fmt.Println(commands)
	for _, c := range commands {
		r := reflect.ValueOf(cmd)
		f := reflect.Indirect(r).FieldByName(c)
		fmt.Println(f)
		// fmt.Println(f.Interface().(Executable))
		// f.Execute([]string{})
		// if err := cmd[a].Execute([]string{}); err != nil {
		// 	return err
		// }
	}
	return nil
}
*/
