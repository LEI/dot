package list
// https://github.com/moby/moby/tree/17.05.x/cli/command/image

import (
	"github.com/spf13/cobra"

	"github.com/LEI/dot/cli"
	"github.com/LEI/dot/cli/command"
	// "github.com/LEI/dot/cli/command/install"
)

// NewListCommand returns a cobra command for `list` subcommands
func NewListCommand(dotCli *command.DotCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List",
		Args:  cli.NoArgs,
		RunE:  dotCli.ShowHelp,
	}
	cmd.AddCommand(
		// install.NewLinkCommand(dotCli),
		// NewHistoryCommand(dotCli),
		// NewImportCommand(dotCli),
		// NewLoadCommand(dotCli),
		// NewPullCommand(dotCli),
		// NewPushCommand(dotCli),
		// NewSaveCommand(dotCli),
		// NewTagCommand(dotCli),
		// newListCommand(dotCli),
		// newRemoveCommand(dotCli),
		// newInspectCommand(dotCli),
		// NewPruneCommand(dotCli),
	)
	return cmd
}
