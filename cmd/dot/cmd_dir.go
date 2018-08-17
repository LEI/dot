package main

import (
	"github.com/spf13/cobra"
)

// DirOptions collects all options for the dir command.
type DirOptions struct {
	// Empty bool
}

var dirOptions DirOptions

var cmdDir = &cobra.Command{
	Use:     "dir [flags]",
	Aliases: []string{"d"},
	Short:   "Dir managed files",
	Long: `
The "dir" task installs or removes directories.
`,
	DisableAutoGenTag: true,
	Args:              cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDir(cmd.Parent().Name(), args)
	},
}

func init() {
	cmdInstall.AddCommand(cmdDir)
	cmdRemove.AddCommand(cmdDir)

	// flags := cmdDir.Flags()
}

// opts DirOptions, gopts GlobalOptions
func runDir(action string, args []string) error {
	for _, r := range globalConfig.Roles {
		for _, d := range r.Dirs {
			if err := runTask(action, d); err != nil {
				return err
			}
		}
	}
	return nil
}
