package cmd

import (
	"fmt"
	// "os"
)

// RemoveCmd ...
type RemoveCmd struct {
	BaseTaskCmd
}

// var removeCmd RemoveCmd

// Execute ...
func (cmd *RemoveCmd) Execute(args []string) error {
	fmt.Println("exec remove cmd", args)
	return nil // executeCommands(cmd, args)
}
