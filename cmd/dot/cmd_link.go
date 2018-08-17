package main

import (
	"github.com/spf13/cobra"
)

// LinkOptions collects all options for the link command.
type LinkOptions struct {
}

var linkOptions LinkOptions

var cmdLink = &cobra.Command{
	Use:     "link [flags]",
	Aliases: []string{"ln"},
	Short:   "Link managed files",
	Long: `
The "link" task installs or removes symlinks.
`,
	DisableAutoGenTag: true,
	Args:              cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runLink(linkOptions, globalOptions, args)
	},
}

func init() {
	cmdInstall.AddCommand(cmdLink)
	cmdRemove.AddCommand(cmdLink)

	// flags := cmdLink.Flags()
}

func runLink(opts LinkOptions, gopts GlobalOptions, args []string) error {
	return nil
}
