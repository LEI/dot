package main

import (
	"github.com/spf13/cobra"
)

// Options for the pkg command.
type pkgOptions struct {
}

var pkgOpts pkgOptions

var cmdPkg = &cobra.Command{
	Use:     "pkg [flags]",
	Aliases: []string{"li"},
	Short:   "Package manager",
	Long: `
The "pkg" task installs or removes packages.
`,
	DisableAutoGenTag: true,
	Args:              cobra.NoArgs,
	RunE:              runPkg,
}

func init() {
	cmdInstall.AddCommand(cmdPkg)

	cmdRmPkg := &cobra.Command{}
	*cmdRmPkg = *cmdPkg
	cmdRemove.AddCommand(cmdRmPkg)

	// flags := cmdPkg.Flags()
}

func runPkg(cmd *cobra.Command, args []string) error {
	action := cmd.Parent().Name()
	for _, r := range globalConfig.Roles {
		for _, l := range r.Pkgs {
			if err := runTask(action, l); err != nil {
				return err
			}
		}
	}
	return nil
}
