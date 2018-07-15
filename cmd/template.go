package cmd

import (
	"fmt"
)

// TemplateCmd ...
type TemplateCmd struct {
	BaseRoleCmd
}

// Execute ...
func (cmd *TemplateCmd) Execute(args []string) error {
	if Verbose > 2 {
		fmt.Println("# TemplateCmd",
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
		err := role.RegisterTemplate(string(p))
		if err != nil {
			return err
		}
	}
	return nil
}
