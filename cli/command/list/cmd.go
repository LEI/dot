package list
// https://github.com/moby/moby/tree/17.05.x/cli/command/image

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/LEI/dot/cli"
	"github.com/LEI/dot/cli/command"
	// "github.com/LEI/dot/cli/config"
)

type listOptions struct {}

// NewListCommand returns a cobra command for `list` subcommands
func NewListCommand(dotCli *command.DotCli) *cobra.Command {
	opts := listOptions{}
	cmd := &cobra.Command{
		Use:   "list",
		Aliases: []string{"ls"},
		Short: "List",
		Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// if len(args) > 0 {
			// 	opts.matchName = args[0]
			// }
			return runList(dotCli, opts)
		},
	}
	return cmd
}

func runList(dotCli *command.DotCli, opts listOptions) error {
	fmt.Fprintf(dotCli.Out(), "RUN LIST %+v\n", opts)
	for _, role := range dotCli.Roles() {
		if err := role.Link.Check(); err != nil {
			return err
		}
	}
	for i, role := range dotCli.Roles() {
		fmt.Println("ROLE", i, role)
	}
	return nil
}
