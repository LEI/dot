package cmd

import (
	"fmt"
)

// LinkCmd ...
type LinkCmd struct {
	BaseRoleCmd
}

// Execute ...
func (cmd *LinkCmd) Execute(args []string) error {
	fmt.Println("execute link command",
		cmd.Role.Name,
		cmd.Role.Paths,
		args)
	if cmd.Role.Name == "" {
		return nil
	}
	// fmt.Println("role:", cmd.Role)
	role := NewRole(cmd.Role.Name)
	if err := role.Register(config); err != nil {
		return err
	}
	for _, p := range cmd.Role.Paths {
		// err := role.RegisterTask("Link", string(p))
		// if err != nil {
		// 	return err
		// }
		err := role.RegisterLink(string(p))
		if err != nil {
			return err
		}
	}
	return nil
}
