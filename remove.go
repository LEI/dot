package main

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
	if Verbose > 2 {
		fmt.Println("# RemoveCmd", args)
	}
	return nil // executeCommands(cmd, args)
}
