package command

import (
	"github.com/spf13/cobra"
)

type installOpts struct {
	// sync bool
}

// NewInstallCommand returns a cobra command for `install` subcommands
func NewInstallCommand(dotCli *DotCli) *cobra.Command {
	a := "install" // action
	// opts := installOpts{}
	cmd := &cobra.Command{
		Use:     "install [OPTIONS] [ACTION]",
		Aliases: []string{"i"},
		Short:   "Install",
		Args:    cobra.NoArgs,
		// PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 	// cmd.Name()
		// 	return nil
		// },
		RunE: func(cmd *cobra.Command, args []string) error {
			// if opts.sync {
			// 	if err := runSync(dotCli, syncOptions{}); err != nil {
			// 		return err
			// 	}
			// }
			if err := runDir(dotCli, dirOptions{action: a}); err != nil {
				return err
			}
			if err := runCopy(dotCli, copyOptions{action: a}); err != nil {
				return err
			}
			if err := runLink(dotCli, linkOptions{action: a}); err != nil {
				return err
			}
			return nil
		},
	}
	dotCli.AddCommands(cmd)
	Options.InstallActionFlags(cmd.Flags())
	Options.InstallActionPersistentFlags(cmd.PersistentFlags())

	// flags := cmd.Flags() // var flags *pflag.FlagSet
	// flags.BoolVarP(&opts.sync, "sync", "S", false, "Clone or pull git repositories")

	return cmd
}
