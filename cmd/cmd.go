package cmd

import (
	"fmt"
	"os"
	// "reflect"
	// "strings"

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

// BaseTaskCmd ...
type BaseTaskCmd struct {
	BaseCmd

	Copy     CopyCmd     `command:"copy" alias:"cp" description:"Copy"`
	Link     LinkCmd     `command:"link" alias:"ln" description:"Symlink"`
	Template TemplateCmd `command:"template" alias:"tpl" description:"Template"`
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
	Name  string // `default:"default"` // `required:"true"`
	Paths []flags.Filename
}

var parser = flags.NewParser(&Options, flags.HelpFlag | flags.PassDoubleDash)

// GetParser ...
func GetParser() *flags.Parser {
	return parser
}

// Parse ...
func Parse() ([]string, error) {
	// TODO: control (mute) output?
	remaining, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			switch flagsErr.Type {
			case flags.ErrHelp:
				parser.WriteHelp(os.Stdout)
				os.Exit(0)
				// break
			// case flags.ErrCommandRequired:
			// 	// FIXME: DotCmd.Execute() never called
			// 	// when first-level sub commands are optional
			// 	err = Options.Install.Execute(remaining)
			// 	remaining = []string{}
			// 	break
			default:
				fmt.Fprintf(os.Stderr, "Error parsing args: %s\n", err)
				parser.WriteHelp(os.Stdout)
				os.Exit(1)
			}
		}
	}
	// WriteIniConfig(parser)
	return remaining, err
}
