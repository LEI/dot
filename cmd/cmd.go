package cmd

import (
	// "fmt"
	// "os"
	// "reflect"
	// "strings"

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
	return dot.NewRole(a.Name)
}

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
