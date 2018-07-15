package cmd

import (
	"fmt"
)

// CopyCmd ...
type CopyCmd struct {
	BaseRoleCmd
}

// Execute ...
func (cmd *CopyCmd) Execute(args []string) error {
	if Verbose > 2 {
		fmt.Println("# CopyCmd",
			cmd.Role.Name,
			cmd.Role.Args,
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
	for _, p := range cmd.Role.Args {
		err := role.RegisterCopy(string(p))
		if err != nil {
			return err
		}
	}
	return nil
}
