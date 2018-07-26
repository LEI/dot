package main

import (
	"fmt"
)

// ListCmd ...
type ListCmd struct {
	BaseTaskCmd
}

// Execute ...
func (cmd *ListCmd) Execute(args []string) error {
	if Verbose > 2 {
		fmt.Println("# ListCmd", args)
	}
	return nil
}
