package cmd

import (
	"fmt"
)

// LinkCmd ...
type LinkCmd struct {
	BaseRoleCmd
}

// Link ...
var Link LinkCmd

// Execute ...
func (cmd *LinkCmd) Execute(args []string) error {
	if cmd.Role.Name == "" {
		return nil
	}
	fmt.Println("execute link command", args) // , cmd)
	// fmt.Println("role:", cmd.Role)
	role := cmd.Role.New() // Init dot.Role
	if err := role.Register(GlobalConfig); err != nil {
		return err
	}
	for _, p := range cmd.Role.Paths {
		err := role.RegisterLink(string(p))
		if err != nil {
			return err
		}
	}
	return nil
}
