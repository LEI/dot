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
	fmt.Println("execute template command",
		cmd.Role.Name,
		cmd.Role.Paths,
		args)
	if cmd.Role.Name == "" {
		return nil
	}
	// fmt.Println("role:", cmd.Role)
	role := NewRole(cmd.Role.Name)
	if err := role.Register(GlobalConfig); err != nil {
		return err
	}
	for _, p := range cmd.Role.Paths {
		err := role.RegisterTemplate(string(p))
		if err != nil {
			return err
		}
	}
	return nil
}
