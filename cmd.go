package main

import (
	// "fmt"
	// "os"
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

	// List     ListCmd     `command:"list" alias:"ls" description:"List files"`
	Copy     CopyCmd     `command:"copy" alias:"cp" description:"Copy files"`
	Line     LineCmd     `command:"line-in-file" alias:"line" description:"Line in file"`
	Link     LinkCmd     `command:"link" alias:"ln" description:"Symlink files"`
	Template TemplateCmd `command:"template" alias:"tpl" description:"Template files"`
	Package  PackageCmd  `command:"package" alias:"p" description:"Manage packages"`
	// Exec     ExecCmd      `command:"exec" alias:"e" description:"Execute a command"`
}

// BaseRoleCmd ...
type BaseRoleCmd struct {
	BaseCmd

	Role struct {
		Name string // `default:"default"` // `required:"true"`
		Args []flags.Filename
	} `positional-args:"true" positional-arg-name:"paths"` // required:"1"

	Executable
}
