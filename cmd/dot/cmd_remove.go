package main

import (
	"github.com/spf13/cobra"
)

var cmdRemove = &cobra.Command{
	Use:     "remove [flags]",
	Aliases: []string{"rm"},
	Short:   "Remove managed files",
	Long: `
The "remove" command removes roles by executing their tasks.
`,
	DisableAutoGenTag: true,
	Args:              cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRemove(removeOptions, globalOptions, args)
	},
}

// RemoveOptions collects all options for the remove command.
type RemoveOptions struct {
	// Empty bool
}

var removeOptions RemoveOptions

func init() {
	cmdRoot.AddCommand(cmdRemove)

	// flags := cmdRemove.Flags()
	// flags.BoolVarP(&removeOptions.Empty, "remove-empty", "", false, "remove empty directories and empty files")
}

func runRemove(opts RemoveOptions, gopts GlobalOptions, args []string) error {
	for _, r := range globalConfig.Roles {
		for _, d := range r.Dirs {
			if err := undoDir(d, DirOptions{}, gopts); err != nil {
				return err
			}
		}
	}
	return nil
}
