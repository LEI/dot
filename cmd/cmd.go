package cmd

import (
	"fmt"
	"os"
	// "reflect"
	// "strings"

	"github.com/jessevdk/go-flags"

	"github.com/LEI/dot/dotfile"
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
	Line     LineCmd     `command:"line-in-file" alias:"line" description:"Line in file"`
	Link     LinkCmd     `command:"link" alias:"ln" description:"Symlink"`
	Template TemplateCmd `command:"template" alias:"tpl" description:"Template"`
	Package  PackageCmd  `command:"package" alias:"p" description:"Manage packages"`
	// Exec     ExecCmd      `command:"exec" alias:"e" description:"Execute a command"`
}

// BaseRoleCmd ...
type BaseRoleCmd struct {
	BaseCmd

	Role RoleArg `positional-args:"true" positional-arg-name:"paths"` // required:"1"

	Executable
}

// RoleArg ...
type RoleArg struct {
	Name  string // `default:"default"` // `required:"true"`
	Paths []flags.Filename
}

var parser = flags.NewParser(&Options, flags.HelpFlag|flags.PassDoubleDash)

var target string

// GetParser ...
func GetParser() *flags.Parser {
	return parser
}

// Parse ...
func Parse() ([]string, error) {
	remaining, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			switch flagsErr.Type {
			case flags.ErrHelp:
				parser.WriteHelp(os.Stdout)
				os.Exit(0)
				// break
			case flags.ErrCommandRequired:
				// err = Options.Install.Execute(remaining)
				// remaining = []string{}
				err = nil
				break
			default:
				fmt.Fprintf(os.Stderr, "Error parsing args: %s\n", err)
				parser.WriteHelp(os.Stdout)
				os.Exit(1)
			}
		}
	}
	// Update variables
	target = string(Options.Target)
	target = dotfile.ExpandEnv(target)
	dotfile.DryRun = Options.DryRun
	Verbose = len(Options.Verbose)
	dotfile.Verbose = Verbose
	// WriteIniConfig(parser)
	return remaining, err
}
