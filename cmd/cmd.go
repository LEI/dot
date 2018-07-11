package cmd

import (
	// "fmt"
	"os"
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

// Parse ...
func Parse() ([]string, error) {
	// TODO: control (mute) output?
	remaining, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			switch flagsErr.Type {
			case flags.ErrHelp:
				os.Exit(1)
			case flags.ErrCommandRequired:
				// FIXME: DotCmd.Execute() never called
				// when first-level sub commands are optional
				err = Options.Install.Execute(remaining)
				remaining = []string{}
				// os.Exit(1)
			// default:
			// 	fmt.Println("Error parsing args:", err)
			// 	os.Exit(1)
			}
		}
	}
	// WriteIniConfig(parser)
	return remaining, err
}

// WriteHelp ...
// func WriteHelp(o io.Writer) {
// 	parser.WriteHelp(o)
// }

// Help ...
func Help(rc int) {
	parser.WriteHelp(os.Stdout)
	os.Exit(rc)
}

// GetParser ...
func GetParser() *flags.Parser {
	return parser
}
