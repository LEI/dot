package command

import (
	"fmt"

	"github.com/LEI/dot/cli/config/tasks"
	"github.com/spf13/cobra"
)

type dirOptions struct {
	action string
	// quiet bool
}

// NewDirCommand creates a new `dot dir` command
func NewDirCommand(dotCli *DotCli) *cobra.Command {
	opts := dirOptions{} // filter: opts.NewFilterOpt()
	cmd := &cobra.Command{
		Use:     "dir [OPTIONS]",
		Aliases: []string{"d", "mkdir"},
		Short:   "Create directory",
		Args:    cobra.NoArgs, // RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.action = cmd.Parent().Name()
			return runDir(dotCli, opts)
		},
	}
	Options.InstallTaskFlags(cmd.Flags())
	Options.InstallTaskPersistentFlags(cmd.PersistentFlags())
	// flags := cmd.Flags()
	// flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only show numeric IDs")
	return cmd
}

// func newDirCommand(dotCli *DotCli) *cobra.Command {
// 	cmd := *NewDirCommand(dotCli)
// 	cmd.Aliases = []string{"d"}
// 	cmd.Use = "dir [OPTIONS]"
// 	return &cmd
// }

func runDir(dotCli *DotCli, opts dirOptions) error {
	roles := dotCli.Roles()
	if len(roles) == 0 {
		return fmt.Errorf("no roles to dir")
	}
	for _, r := range roles {
		if err := tasks.Check(r.Dirs); err != nil {
			return err
		}
	}
	for _, r := range roles {
		switch opts.action {
		case "install":
			if err := tasks.Install(r.Dirs); err != nil {
				return err
			}
		case "remove":
			if err := tasks.Remove(r.Dirs); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%s: not implemented", opts.action)
		}
	}
	return nil
}
