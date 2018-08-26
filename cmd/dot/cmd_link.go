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
	Aliases: []string{"ln"}, // ~stow
	Short:   "Link managed files",
	Long:    `The "link" task installs or removes symlinks.`,
	Args:    cobra.NoArgs,
	PreRunE: preRunTask,
	RunE:    runLink,
	// DisableAutoGenTag: true,
}

var cmdRmLink = &cobra.Command{}

func init() {
	cmdInstall.AddCommand(cmdLink)
	*cmdRmLink = *cmdLink
	cmdRemove.AddCommand(cmdRmLink)

	// flags := cmdLink.Flags()
}

func runLink(cmd *cobra.Command, args []string) error {
	action := cmd.Parent().Name()
	for _, r := range dotConfig.Roles {
		for _, l := range r.Links {
			if err := runTask(action, l); err != nil {
				return err
			}
		}
	}
	return nil
}
