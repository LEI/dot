package command

import (
	"github.com/spf13/cobra"
)

type installOpts struct {
	sync bool
}

// NewInstallCommand returns a cobra command for `install` subcommands
func NewInstallCommand(dotCli *DotCli) *cobra.Command {
	a := "install" // action
	opts := installOpts{}
	cmd := &cobra.Command{
		Use:   "install [OPTIONS]",
		Aliases: []string{"i"},
		Short: "Install",
		Args:  cobra.NoArgs,
		// PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 	// cmd.Name()
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
			if err := runDir(dotCli, dirOptions{action: a}); err != nil {
				return err
			}
			if err := runLink(dotCli, linkOptions{action: a}); err != nil {
				return err
			}
			return nil
		},
	}
	flags := cmd.Flags() // var flags *pflag.FlagSet
	flags.BoolVarP(&opts.sync, "sync", "S", false, "Clone or pull git repositories")
	dotCli.AddCommands(cmd)
	return cmd
}
