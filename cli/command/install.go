package command

import (
	"github.com/spf13/cobra"
)

// NewInstallCommand returns a cobra command for `install` subcommands
func NewInstallCommand(dotCli *DotCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Aliases: []string{"i"},
		Short: "Install",
		// Args:  cobra.NoArgs,
		// PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 	// cmd.Name()
		// 	return nil
		// },
		// RunE: dotCli.ShowHelp,
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
			if err := runDir(dotCli, dirOptions{
				action: "install",
			}); err != nil {
				return err
			}
			if err := runLink(dotCli, linkOptions{
				action: "install",
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
