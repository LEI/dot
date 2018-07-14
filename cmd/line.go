package cmd

import (
	"fmt"
)

// LineCmd ...
type LineCmd struct {
	BaseRoleCmd
}

// Execute ...
func (cmd *LineCmd) Execute(args []string) error {
	if Verbose > 2 {
		fmt.Println("# LineCmd",
			cmd.Role.Name,
			cmd.Role.Paths,
			args)
	}
	if cmd.Role.Name == "" {
		return nil
	}
	// fmt.Println("role:", cmd.Role)
	role := NewRole(cmd.Role.Name)
	if err := role.Register(config); err != nil {
		return err
	}
	for _, p := range cmd.Role.Paths {
		err := role.RegisterTask("Line", string(p))
		if err != nil {
			return err
		}
		// err := role.RegisterLine(string(p))
		// if err != nil {
		// 	return err
		// }
	}
	return nil
}
