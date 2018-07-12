package cmd

import ()

// LineCmd ...
type LineCmd struct {
	BaseRoleCmd
}

// Execute ...
func (cmd *LineCmd) Execute(args []string) error {
	// fmt.Println("execute line command",
	// 	cmd.Role.Name,
	// 	cmd.Role.Paths,
	// 	args)
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
