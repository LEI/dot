package main

import (
	"github.com/spf13/cobra"
)

// Options for the pkg command.
type packageOptions struct {
}

var pkgOpts packageOptions

var cmdPkg = &cobra.Command{
	Use:     "package [flags]",
	Aliases: []string{"pkg"},
	Short:   "Package tasks",
	Long:    `The "pkg" task installs or removes packages.`,
	Args:    cobra.NoArgs,
	PreRunE: preRunTask,
	RunE:    runPkg,
	// DisableAutoGenTag: true,
}

var cmdRmPkg = &cobra.Command{}

func init() {
	cmdInstall.AddCommand(cmdPkg)
	*cmdRmPkg = *cmdPkg
	cmdRemove.AddCommand(cmdRmPkg)

	// flags := cmdPkg.Flags()
}

func runPkg(cmd *cobra.Command, args []string) error {
	action := cmd.Parent().Name()
	for _, r := range dotConfig.Roles {
		for _, l := range r.Pkgs {
			if err := runTask(action, l); err != nil {
				return err
			}
		}
	}
	return nil
}
