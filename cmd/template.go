package cmd

import (
	"fmt"

	"github.com/jessevdk/go-flags"
)

// TemplateCmd ...
type TemplateCmd struct {
	*BaseCmd

	Source string `short:"s" long:"source" description:""`
	Target string `short:"t" long:"target" description:""`
	// *TemplateArg `required:"true" positional-args:"true"`
}

// TemplateArg ...
type TemplateArg struct {
	Source flags.Filename `name:"SOURCE" description:"Path to template file"`
	Target flags.Filename `name:"TARGET" description:"Path to target file"`
}

// Template ...
var Template TemplateCmd

// Execute ...
func (cmd *TemplateCmd) Execute(args []string) error {
	fmt.Println("execute template command", args)
	if cmd.Source == "" && cmd.Target == "" {
		return nil
	}

	fmt.Println("tpl", cmd.Source, cmd.Target)

	return nil
}
