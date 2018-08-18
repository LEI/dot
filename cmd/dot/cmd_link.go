package main

import (
	"github.com/spf13/cobra"
)

// Options for the link command.
type linkOptions struct {
}

var linkOpts linkOptions

var cmdLink = &cobra.Command{
	Use:     "link [flags]",
	Aliases: []string{"ln"},
	Short:   "Link managed files",
	Long: `
The "link" task installs or removes symlinks.
`,
	DisableAutoGenTag: true,
	Args:              cobra.NoArgs,
	RunE:              runLink,
}

func init() {
	cmdInstall.AddCommand(cmdLink)
	cmdRemove.AddCommand(cmdLink)

	// flags := cmdLink.Flags()
}

func runLink(cmd *cobra.Command, args []string) error {
	action := cmd.Parent().Name()
	for _, r := range globalConfig.Roles {
		for _, l := range r.Links {
			if err := runTask(action, l); err != nil {
				return err
			}
		}
	}
	return nil
}
