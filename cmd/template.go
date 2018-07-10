package cmd

import (
	"fmt"
)

// TemplateCmd ...
type TemplateCmd struct {
	BaseRoleCmd
}

// Template ...
var Template TemplateCmd

// Execute ...
func (cmd *TemplateCmd) Execute(args []string) error {
	if cmd.Role.Name == "" {
		return nil
	}
	fmt.Println("execute template command", args)
	// fmt.Println("role:", cmd.Role)
	role := cmd.Role.New() // Init dot.Role
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
