package cmd

import (
	"fmt"

	"github.com/jessevdk/go-flags"

	"github.com/LEI/dot/dot"
)

// LinkCmd ...
type LinkCmd struct {
	*BaseCmd

	*LinkArg // `required:"true" positional-args:"true"`
}

// LinkArg ...
type LinkArg struct {
	// Paths map[flags.Filename]flags.Filename `name:"PATHS" description:"Paths"`
	// Paths map[string]string `name:"PATHS" description:"Paths"`
	// Source flags.Filename `name:"SOURCE" description:"Path to source file"`
	// Target flags.Filename `name:"TARGET" description:"Path to target link"`
	Source flags.Filename `short:"s" long:"source" description:"Path to source file"`
	Target flags.Filename `short:"t" long:"target" description:"Path to target link"`
}

// Link ...
var Link LinkCmd

// Execute ...
func (cmd *LinkCmd) Execute(args []string) error {
	fmt.Println("execute link command", args)
	if cmd.Source == "" && cmd.Target == "" {
		return nil
	}

	s := string(cmd.Source)
	t := string(cmd.Target)

	// fmt.Printf("%v\n", cmd)

	return dot.RegisterLink(s, t)
}
