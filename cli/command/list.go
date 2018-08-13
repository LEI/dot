package command

// https://github.com/moby/moby/tree/17.05.x/cli/command/image

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/LEI/dot/cli/config/tasks"
)

type listOptions struct{}

// NewListCommand returns a cobra command for `list` subcommands
func NewListCommand(dotCli *DotCli) *cobra.Command {
	opts := listOptions{}
	cmd := &cobra.Command{
		Use:     "list [OPTIONS]",
		Aliases: []string{"ls"},
		Short:   "List",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// if len(args) > 0 {
			// 	opts.matchName = args[0]
			// }
			return runList(dotCli, opts)
		},
	}
	return cmd
}

func runList(dotCli *DotCli, opts listOptions) error {
	roles := dotCli.Roles()
	if len(roles) == 0 {
		return fmt.Errorf("no roles to list")
	}
	for i, r := range roles {
		fmt.Fprintf(dotCli.Out(), "Role #%d: %s\n", i+1, r.Name)
		if tasks.Verbose > 0 {
			// fmt.Fprintf(dotCli.Out(), "Role #%d: %+v\n", i+1, r)
			// fmt.Fprintf(dotCli.Out(), " OS: %+v\n", r.OS)
			fmt.Fprintf(dotCli.Out(), " Env: %d\n", len(r.Env))
			for k, v := range r.Env {
				fmt.Fprintf(dotCli.Out(), "   %s=%s\n", k, v)
			}
			fmt.Fprintf(dotCli.Out(), " Dirs: %d\n", len(r.Dirs))
			for _, d := range r.Dirs {
				fmt.Fprintf(dotCli.Out(), "   %+v\n", d)
			}
			fmt.Fprintf(dotCli.Out(), " Files: %d\n", len(r.Files))
			for _, d := range r.Files {
				fmt.Fprintf(dotCli.Out(), "   %+v\n", d)
			}
			// fmt.Fprintf(dotCli.Out(), " Links: %d\n", len(r.Links))
			// for _, l := range r.Links {
			// 	fmt.Fprintf(dotCli.Out(), "   %+v\n", l)
			// }
		}
	}
	return nil
}
