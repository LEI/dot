package command

import (
	"github.com/spf13/cobra"
)

type removeOpts struct {
	removeEmpty bool
}

// NewRemoveCommand returns a cobra command for `remove` subcommands
func NewRemoveCommand(dotCli *DotCli) *cobra.Command {
	a := "remove"
	// opts := removeOpts{}
	cmd := &cobra.Command{
		Use:   "remove",
		Aliases: []string{"rm"},
		Short: "Remove",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runLink(dotCli, linkOptions{action: a}); err != nil {
				return err
			}
			if err := runDir(dotCli, dirOptions{action: a}); err != nil {
				return err
			}
			return nil
		},
	}
	// flags := cmd.Flags() // var flags *pflag.FlagSet
	// flags.BoolVarP(&opts.removeEmpty, "remove-empty", "R", true, "Remove empty files and directories")
	dotCli.AddCommands(cmd)
	return cmd
}
