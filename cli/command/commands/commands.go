package commands

import (
	"github.com/LEI/dot/cli/command"
	"github.com/LEI/dot/cli/command/list"
	"github.com/LEI/dot/cli/command/install"
	// "github.com/LEI/dot/cli/command/remove"
	"github.com/spf13/cobra"
)

// AddCommands adds all the commands from cli/command to the root command
func AddCommands(cmd *cobra.Command, dockerCli *command.DotCli) {
	cmd.AddCommand(
		list.NewListCommand(dockerCli),
		install.NewInstallCommand(dockerCli),
		// remove.NewRemoveCommand(dockerCli),
	)
}
