package cmd

import (
	"fmt"
	// "os"
)

// RemoveCmd ...
type RemoveCmd struct {
	*BaseCmd

	Copy     CopyCmd     `command:"copy" alias:"cp" description:""`
	Link     LinkCmd     `command:"link" alias:"ln" description:""`
	Template TemplateCmd `command:"template" alias:"tpl" description:""`
}

var removeCmd RemoveCmd

// Execute ...
func (cmd *RemoveCmd) Execute(args ...[]string) error {
	fmt.Println("exec remove cmd", args)
	return nil
}

func init() {
	// Remove, err := parser.AddCommand("remove",
	// 	"Remove",
	// 	"",
	// 	&removeCmd)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// Remove.Aliases = []string{"rm"}
	// Remove.SubcommandsOptional = true

	// RemoveCopy, err := Remove.AddCommand("copy",
	// 	"Remove copy",
	// 	"",
	// 	&copyCmd)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// RemoveCopy.Aliases = []string{"rm"}
}
