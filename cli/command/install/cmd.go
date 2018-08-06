package install

import (
	"github.com/spf13/cobra"

	"github.com/LEI/dot/cli"
	"github.com/LEI/dot/cli/command"
)

// NewInstallCommand returns a cobra command for `install` subcommands
func NewInstallCommand(dotCli *command.DotCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install",
		Args:  cli.NoArgs,
		RunE:  dotCli.ShowHelp,
	}
	cmd.AddCommand(
		NewLinkCommand(dotCli),
	)
	return cmd
}
