package main

import (
	"github.com/spf13/cobra"
)

var cmdInstall = &cobra.Command{
	Use:     "install [flags]",
	Aliases: []string{"i"},
	Short:   "Install managed files",
	Long: `
The "install" command installs roles by executing their tasks.
`,
	DisableAutoGenTag: true,
	Args:              cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInstall(installOptions, globalOptions, args)
	},
}

// InstallOptions collects all options for the install command.
type InstallOptions struct {
	Sync bool
}

var installOptions InstallOptions

func init() {
	cmdRoot.AddCommand(cmdInstall)

	flags := cmdInstall.Flags()
	flags.BoolVarP(&installOptions.Sync, "sync", "S", false, "synchronize repositories")
	// flags.BoolVarP(&installOptions.Pull, "pull", "p", false, "pull repositories")
	// flags.BoolVarP(&installOptions.Force, "force", "f", false, "ignore uncommitted changes")
}

func runInstall(opts InstallOptions, gopts GlobalOptions, args []string) error {
	// if opts.Sync {
	// 	for _, r := range globalConfig.Roles {
	// 		if err := r.RunSync(); err != nil {
	// 			return err
	// 		}
	// 	}
	// }
	// for _, r := range globalConfig.Roles {
	// 	if err := r.RunInstall(); err != nil {
	// 		return err
	// 	}
	// }
	for _, r := range globalConfig.Roles {
		for _, d := range r.Dirs {
			if err := doDir(d, DirOptions{}, gopts); err != nil {
				return err
			}
		}
	}
	return nil
}
