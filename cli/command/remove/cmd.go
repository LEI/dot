package remove

import (
	"github.com/spf13/cobra"

	"github.com/LEI/dot/cli"
	"github.com/LEI/dot/cli/command"
)

// NewRemoveCommand returns a cobra command for `remove` subcommands
func NewRemoveCommand(dotCli *command.DotCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Aliases: []string{"rm"},
		Short: "Remove",
		Args:  cli.NoArgs,
		RunE:  dotCli.ShowHelp,
		// PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 	return nil
		// },
	}
	cmd.AddCommand(
		NewRmDirCommand(dotCli),
		NewRmLinkCommand(dotCli),
	)
	return cmd
}
