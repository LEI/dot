package install

import (
	"github.com/spf13/cobra"

	"github.com/LEI/dot/cli/command"
)

// var (
// 	installAction
// )

// NewInstallCommand returns a cobra command for `install` subcommands
func NewInstallCommand(dotCli *command.DotCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Aliases: []string{"i"},
		Short: "Install",
		Args:  cobra.NoArgs,
		RunE:  dotCli.ShowHelp,
		// PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 	installAction = cmd.Name()
		// 	return nil
		// },
	}
	cmd.AddCommand(
		NewDirCommand(dotCli),
		// NewCopyCommand(dotCli),
		NewLinkCommand(dotCli),
		// NewTemplateCommand(dotCli),
	)
	return cmd
}
