package cmd

import (
	"fmt"
)

// ListCmd ...
type ListCmd struct {
	BaseRoleCmd
}

// Execute ...
func (cmd *ListCmd) Execute(args []string) error {
	if Verbose > 2 {
		fmt.Println("# ListCmd",
			cmd.Role.Name,
			cmd.Role.Args,
			args)
	}
	if cmd.Role.Name == "" {
		return nil
	}
	// role := NewRole(cmd.Role.Name)
	// if err := role.Register(config); err != nil {
	// 	return err
	// }
	// for _, p := range cmd.Role.Args {
	// 	err := role.RegisterList(string(p))
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}
