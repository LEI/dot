package main

import (
	"fmt"
)

// PackageCmd ...
type PackageCmd struct {
	BaseRoleCmd
}

// Execute ...
func (cmd *PackageCmd) Execute(args []string) error {
	if Verbose > 2 {
		fmt.Println("# PackageCmd", cmd.Role.Name, args)
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
		fmt.Println("=>", p)
		// err := role.RegisterPackage(string(p))
		// if err != nil {
		// 	return err
		// }
	}
	return nil
}
