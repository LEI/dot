package cmd

import (
	"fmt"
	// "os"

	"github.com/jessevdk/go-flags"
)

// CopyCmd ...
type CopyCmd struct {
	*BaseCmd

	Source string `short:"s" long:"source" description:""`
	Target string `short:"t" long:"target" description:""`
	// *CopyArg `required:"true" positional-args:"true"`
}

// CopyArg ...
type CopyArg struct {
	Source flags.Filename `name:"SOURCE" description:"Path to source file"`
	Target flags.Filename `name:"TARGET" description:"Path to target file"`
}

// Copy ...
var copyCmd CopyCmd

// Execute ...
func (cmd *CopyCmd) Execute(args []string) error {
	fmt.Println("execute copy command", args)
	if cmd.Source == "" && cmd.Target == "" {
		return nil
	}

	fmt.Println("cp", cmd.Source, cmd.Target)

	return nil
}

func init() {
}
