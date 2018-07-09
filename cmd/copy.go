package cmd

import (
)

// CopyCmd ...
type CopyCmd struct {
	BaseRoleCmd
}

// Copy ...
var Copy CopyCmd

// Execute ...
func (cmd *CopyCmd) Execute(args []string) error {
	role := cmd.Role.New() // Init dot.Role
	if err := role.Register(GlobalConfig); err != nil {
		return err
	}
	for _, p := range cmd.Role.Paths {
		err := role.RegisterCopy(string(p))
		if err != nil {
			return err
		}
	}
	return nil
}
