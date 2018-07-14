package cmd

import (
	"fmt"
)

// PkgCmd ...
type PkgCmd struct {
	BaseRoleCmd
}

// Execute ...
func (cmd *PkgCmd) Execute(args []string) error {
	if Verbose > 2 {
		fmt.Println("# PkgCmd", cmd.Role.Name, args)
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
		fmt.Println("=>", p)
		// err := role.RegisterPkg(string(p))
		// if err != nil {
		// 	return err
		// }
	}
	return nil
}
