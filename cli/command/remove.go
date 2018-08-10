package command

import (
	// cliflags "github.com/LEI/dot/cli/flags"
	"github.com/spf13/cobra"
)

// TODO
type removeOpts struct {
	removeEmpty bool
}

// NewRemoveCommand returns a cobra command for `remove` subcommands
func NewRemoveCommand(dotCli *DotCli) *cobra.Command {
	a := "remove"
	// opts := removeOpts{}
	cmd := &cobra.Command{
		Use:     "remove [OPTIONS] [ACTION]",
		Aliases: []string{"rm"},
		Short:   "Remove",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runLink(dotCli, linkOptions{action: a}); err != nil {
				return err
			}
			if err := runCopy(dotCli, copyOptions{action: a}); err != nil {
				return err
			}
			if err := runDir(dotCli, dirOptions{action: a}); err != nil {
				return err
			}
			return nil
		},
		// PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 	dotCli.InitializeAction(..., cmd.PersistentFlags())
		// 	return nil
		// },
	}
	dotCli.AddCommands(cmd)
	Options.InstallActionFlags(cmd.Flags())
	Options.InstallActionPersistentFlags(cmd.PersistentFlags())

	// flags := cmd.Flags() // var flags *pflag.FlagSet
	// flags.BoolVarP(&opts.removeEmpty, "remove-empty", "R", true, "Remove empty files and directories")

	return cmd
}
