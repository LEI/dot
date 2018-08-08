package command

import (
	"github.com/spf13/cobra"
)

// NewRemoveCommand returns a cobra command for `remove` subcommands
func NewRemoveCommand(dotCli *DotCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Aliases: []string{"rm"},
		Short: "Remove",
		Args:  cobra.NoArgs,
		// PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 	return nil
		// },
		RunE: func(cmd *cobra.Command, args []string) error {
			// for _, c := range cmd.Commands() {
			// 	c.SetArgs([]string{})
			// 	if err := c.Execute(); err != nil {
			// 		return err
			// 	}
			// }
			// if err := NewLinkCommand(dotCli).Execute(); err != nil {
			// 	return err
			// }
			if err := runLink(dotCli, linkOptions{
				action: "remove",
			}); err != nil {
				return err
			}
			if err := runDir(dotCli, dirOptions{
				action: "remove",
			}); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.AddCommand(
		NewDirCommand(dotCli),
		NewLinkCommand(dotCli),
	)
	return cmd
}
