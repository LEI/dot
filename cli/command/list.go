package command
// https://github.com/moby/moby/tree/17.05.x/cli/command/image

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/LEI/dot/cli/config/tasks"
)

type listOptions struct {}

// NewListCommand returns a cobra command for `list` subcommands
func NewListCommand(dotCli *DotCli) *cobra.Command {
	opts := listOptions{}
	cmd := &cobra.Command{
		Use:   "list",
		Aliases: []string{"ls"},
		Short: "List",
		Args:  cobra.NoArgs,
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
	for i, r := range dotCli.Roles() {
		if tasks.Verbose {
			fmt.Printf("Role #%d: %+v\n", i+1, r)
		} else {
			fmt.Printf("Role #%d: %s\n", i+1, r.Name)
		}
	}
	return nil
}
