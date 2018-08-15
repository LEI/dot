package command

import (
	"fmt"

	"github.com/LEI/dot/cli/config/tasks"
	"github.com/spf13/cobra"
)

type lineOptions struct {
	action string
}

// NewLineCommand creates a new `dot line` command
func NewLineCommand(dotCli *DotCli) *cobra.Command {
	opts := lineOptions{} // filter: opts.NewFilterOpt()
	cmd := &cobra.Command{
		Use:     "line [OPTIONS]",
		Aliases: []string{},
		Short:   "Line in file",
		Args:    cobra.NoArgs, // RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.action = cmd.Parent().Name()
			return runLine(dotCli, opts)
		},
	}
	Options.InstallTaskFlags(cmd.Flags())
	Options.InstallTaskPersistentFlags(cmd.PersistentFlags())
	return cmd
}

func runLine(dotCli *DotCli, opts lineOptions) error {
	roles := dotCli.Roles()
	if len(roles) == 0 {
		return fmt.Errorf("no roles to line")
	}
	for _, r := range roles {
		if err := tasks.Check(r.Lines); err != nil {
			return err
		}
	}
	for _, r := range roles {
		switch opts.action {
		case "install":
			if err := tasks.Install(r.Lines); err != nil {
				return err
			}
		case "remove":
			if err := tasks.Remove(r.Lines); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%s: not implemented", opts.action)
		}
	}
	return nil
}
